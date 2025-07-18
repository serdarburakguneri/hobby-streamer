package auth

import (
	"context"
	"crypto/rsa"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/constants"
)

type KeycloakValidator struct {
	keycloakURL string
	realm       string
	clientID    string
	httpClient  *http.Client
	keys        map[string]*rsa.PublicKey
	keysMutex   sync.RWMutex
	lastFetch   time.Time
}

type JWKSResponse struct {
	Keys []JWK `json:"keys"`
}

type JWK struct {
	Kid string `json:"kid"`
	Kty string `json:"kty"`
	Alg string `json:"alg"`
	Use string `json:"use"`
	N   string `json:"n"`
	E   string `json:"e"`
}

func NewKeycloakValidator(keycloakURL, realm, clientID string) *KeycloakValidator {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // nolint:gosec // Development only - disable TLS verification for local Keycloak
	}
	httpClient := &http.Client{
		Transport: tr,
		Timeout:   30 * time.Second,
	}

	return &KeycloakValidator{
		keycloakURL: keycloakURL,
		realm:       realm,
		clientID:    clientID,
		httpClient:  httpClient,
		keys:        make(map[string]*rsa.PublicKey),
	}
}

func (k *KeycloakValidator) ValidateToken(ctx context.Context, token string) (*User, error) {
	token = strings.TrimPrefix(token, constants.BearerPrefix)

	parsedToken, _, err := jwt.NewParser().ParseUnverified(token, jwt.MapClaims{})
	if err != nil {
		return nil, fmt.Errorf("invalid token format: %w", err)
	}

	header, ok := parsedToken.Header["kid"].(string)
	if !ok {
		return nil, fmt.Errorf("missing key ID in token header")
	}

	publicKey, err := k.getPublicKey(ctx, header)
	if err != nil {
		return nil, fmt.Errorf("failed to get public key: %w", err)
	}

	parsedToken, err = jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return publicKey, nil
	})
	if err != nil {
		return nil, fmt.Errorf("token signature validation failed: %w", err)
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	if err := k.validateClaims(claims); err != nil {
		return nil, err
	}

	user := &User{
		ID:       getStringClaim(claims, "sub"),
		Username: getStringClaim(claims, "preferred_username"),
		Email:    getStringClaim(claims, "email"),
	}

	if realmAccess, ok := claims["realm_access"].(map[string]interface{}); ok {
		if roles, ok := realmAccess["roles"].([]interface{}); ok {
			for _, role := range roles {
				if roleStr, ok := role.(string); ok {
					user.Roles = append(user.Roles, roleStr)
				}
			}
		}
	}

	return user, nil
}

func (k *KeycloakValidator) getPublicKey(ctx context.Context, kid string) (*rsa.PublicKey, error) {
	k.keysMutex.RLock()
	if key, exists := k.keys[kid]; exists && time.Since(k.lastFetch) < 5*time.Minute {
		k.keysMutex.RUnlock()
		return key, nil
	}
	k.keysMutex.RUnlock()

	k.keysMutex.Lock()
	defer k.keysMutex.Unlock()

	jwksURL := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/certs", k.keycloakURL, k.realm)
	resp, err := k.httpClient.Get(jwksURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch JWKS: %w", err)
	}
	defer resp.Body.Close()

	var jwks JWKSResponse
	if err := json.NewDecoder(resp.Body).Decode(&jwks); err != nil {
		return nil, fmt.Errorf("failed to decode JWKS: %w", err)
	}

	for _, jwk := range jwks.Keys {
		if jwk.Kid == kid {
			publicKey, err := k.jwkToPublicKey(jwk)
			if err != nil {
				return nil, fmt.Errorf("failed to convert JWK to public key: %w", err)
			}
			k.keys[kid] = publicKey
			k.lastFetch = time.Now()
			return publicKey, nil
		}
	}

	return nil, fmt.Errorf("key with kid %s not found", kid)
}

func (k *KeycloakValidator) jwkToPublicKey(jwk JWK) (*rsa.PublicKey, error) {
	nBytes, err := base64.RawURLEncoding.DecodeString(jwk.N)
	if err != nil {
		return nil, fmt.Errorf("failed to decode modulus: %w", err)
	}

	eBytes, err := base64.RawURLEncoding.DecodeString(jwk.E)
	if err != nil {
		return nil, fmt.Errorf("failed to decode exponent: %w", err)
	}

	n := new(big.Int).SetBytes(nBytes)
	e := new(big.Int).SetBytes(eBytes)

	return &rsa.PublicKey{
		N: n,
		E: int(e.Int64()),
	}, nil
}

func (k *KeycloakValidator) validateClaims(claims jwt.MapClaims) error {
	if exp, ok := claims["exp"].(float64); ok {
		if time.Unix(int64(exp), 0).Before(time.Now()) {
			return fmt.Errorf(constants.ErrTokenExpired)
		}
	}

	if iss, ok := claims["iss"].(string); ok {
		expectedIss := fmt.Sprintf("%s/realms/%s", k.keycloakURL, k.realm)
		if iss != expectedIss {
			return fmt.Errorf("invalid issuer: expected %s, got %s", expectedIss, iss)
		}
	}

	if aud, ok := claims["aud"].(string); ok {
		if aud != k.clientID {
			return fmt.Errorf("invalid audience: expected %s, got %s", k.clientID, aud)
		}
	}

	return nil
}

func (k *KeycloakValidator) HasRole(user *User, role string) bool {
	for _, userRole := range user.Roles {
		if userRole == role {
			return true
		}
	}
	return false
}

func (k *KeycloakValidator) HasAnyRole(user *User, roles []string) bool {
	for _, requiredRole := range roles {
		if k.HasRole(user, requiredRole) {
			return true
		}
	}
	return false
}

func (k *KeycloakValidator) HasAllRoles(user *User, roles []string) bool {
	for _, requiredRole := range roles {
		if !k.HasRole(user, requiredRole) {
			return false
		}
	}
	return true
}

func getStringClaim(claims jwt.MapClaims, key string) string {
	if val, ok := claims[key].(string); ok {
		return val
	}
	return ""
}
