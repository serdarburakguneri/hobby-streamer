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
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/errors"
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
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
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
		return nil, errors.NewInternalError("invalid token format", err)
	}

	header, ok := parsedToken.Header["kid"].(string)
	if !ok {
		return nil, errors.NewInternalError("missing key ID in token header", nil)
	}

	publicKey, err := k.getPublicKey(ctx, header)
	if err != nil {
		return nil, errors.NewInternalError("failed to get public key", err)
	}

	parsedToken, err = jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.NewInternalError(fmt.Sprintf("unexpected signing method: %v", token.Header["alg"]), nil)
		}
		return publicKey, nil
	})
	if err != nil {
		return nil, errors.NewInternalError("token signature validation failed", err)
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.NewInternalError("invalid token claims", nil)
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
		return nil, errors.NewInternalError("failed to fetch JWKS", err)
	}
	defer resp.Body.Close()

	var jwks JWKSResponse
	if err := json.NewDecoder(resp.Body).Decode(&jwks); err != nil {
		return nil, errors.NewInternalError("failed to decode JWKS", err)
	}

	for _, jwk := range jwks.Keys {
		if jwk.Kid == kid {
			publicKey, err := k.jwkToPublicKey(jwk)
			if err != nil {
				return nil, errors.NewInternalError("failed to convert JWK to public key", err)
			}
			k.keys[kid] = publicKey
			k.lastFetch = time.Now()
			return publicKey, nil
		}
	}

	return nil, errors.NewInternalError(fmt.Sprintf("key with kid %s not found", kid), nil)
}

func (k *KeycloakValidator) jwkToPublicKey(jwk JWK) (*rsa.PublicKey, error) {
	nBytes, err := base64.RawURLEncoding.DecodeString(jwk.N)
	if err != nil {
		return nil, errors.NewInternalError("failed to decode modulus", err)
	}

	eBytes, err := base64.RawURLEncoding.DecodeString(jwk.E)
	if err != nil {
		return nil, errors.NewInternalError("failed to decode exponent", err)
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
			return errors.NewInternalError(constants.ErrTokenExpired, nil)
		}
	}

	if iss, ok := claims["iss"].(string); ok {
		expectedIss := fmt.Sprintf("%s/realms/%s", k.keycloakURL, k.realm)
		if iss != expectedIss {
			return errors.NewInternalError(fmt.Sprintf("invalid issuer: expected %s, got %s", expectedIss, iss), nil)
		}
	}

	if aud, ok := claims["aud"].(string); ok {
		if aud != k.clientID {
			return errors.NewInternalError(fmt.Sprintf("invalid audience: expected %s, got %s", k.clientID, aud), nil)
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
