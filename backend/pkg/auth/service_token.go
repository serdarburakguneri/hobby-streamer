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

type ServiceTokenValidator struct {
	keycloakURL string
	realm       string
	clientID    string
	httpClient  *http.Client
	keys        map[string]*rsa.PublicKey
	keysMutex   sync.RWMutex
	lastFetch   time.Time
}

type ServiceUser struct {
	ID       string   `json:"id"`
	Username string   `json:"username"`
	Email    string   `json:"email"`
	Roles    []string `json:"roles"`
	ClientID string   `json:"client_id"`
}

func NewServiceTokenValidator(keycloakURL, realm, clientID string) *ServiceTokenValidator {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	httpClient := &http.Client{
		Transport: tr,
		Timeout:   30 * time.Second,
	}

	return &ServiceTokenValidator{
		keycloakURL: keycloakURL,
		realm:       realm,
		clientID:    clientID,
		httpClient:  httpClient,
		keys:        make(map[string]*rsa.PublicKey),
	}
}

func (s *ServiceTokenValidator) ValidateServiceToken(ctx context.Context, token string) (*ServiceUser, error) {
	token = strings.TrimPrefix(token, constants.BearerPrefix)

	parsedToken, _, err := jwt.NewParser().ParseUnverified(token, jwt.MapClaims{})
	if err != nil {
		return nil, fmt.Errorf("invalid token format: %w", err)
	}

	header, ok := parsedToken.Header["kid"].(string)
	if !ok {
		return nil, fmt.Errorf("missing key ID in token header")
	}

	publicKey, err := s.getPublicKey(ctx, header)
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

	if err := s.validateServiceClaims(claims); err != nil {
		return nil, err
	}

	serviceUser := &ServiceUser{
		ID:       getStringClaim(claims, "sub"),
		Username: getStringClaim(claims, "preferred_username"),
		Email:    getStringClaim(claims, "email"),
		ClientID: getStringClaim(claims, "azp"),
	}

	if realmAccess, ok := claims["realm_access"].(map[string]interface{}); ok {
		if roles, ok := realmAccess["roles"].([]interface{}); ok {
			for _, role := range roles {
				if roleStr, ok := role.(string); ok {
					serviceUser.Roles = append(serviceUser.Roles, roleStr)
				}
			}
		}
	}

	return serviceUser, nil
}

func (s *ServiceTokenValidator) getPublicKey(ctx context.Context, kid string) (*rsa.PublicKey, error) {
	s.keysMutex.RLock()
	if key, exists := s.keys[kid]; exists && time.Since(s.lastFetch) < 5*time.Minute {
		s.keysMutex.RUnlock()
		return key, nil
	}
	s.keysMutex.RUnlock()

	s.keysMutex.Lock()
	defer s.keysMutex.Unlock()

	jwksURL := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/certs", s.keycloakURL, s.realm)
	resp, err := s.httpClient.Get(jwksURL)
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
			publicKey, err := s.jwkToPublicKey(jwk)
			if err != nil {
				return nil, fmt.Errorf("failed to convert JWK to public key: %w", err)
			}
			s.keys[kid] = publicKey
			s.lastFetch = time.Now()
			return publicKey, nil
		}
	}

	return nil, fmt.Errorf("key with kid %s not found", kid)
}

func (s *ServiceTokenValidator) jwkToPublicKey(jwk JWK) (*rsa.PublicKey, error) {
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

func (s *ServiceTokenValidator) validateServiceClaims(claims jwt.MapClaims) error {
	if exp, ok := claims["exp"].(float64); ok {
		if time.Unix(int64(exp), 0).Before(time.Now()) {
			return fmt.Errorf(constants.ErrTokenExpired)
		}
	}

	if iss, ok := claims["iss"].(string); ok {
		expectedIss := fmt.Sprintf("%s/realms/%s", s.keycloakURL, s.realm)
		if iss != expectedIss {
			return fmt.Errorf("invalid issuer: expected %s, got %s", expectedIss, iss)
		}
	}

	if aud, ok := claims["aud"].(string); ok {
		if aud != s.clientID {
			return fmt.Errorf("invalid audience: expected %s, got %s", s.clientID, aud)
		}
	}

	tokenType, ok := claims["typ"].(string)
	if !ok || tokenType != "Bearer" {
		return fmt.Errorf("invalid token type: expected Bearer")
	}

	return nil
}

func (s *ServiceTokenValidator) IsServiceToken(user *ServiceUser) bool {
	return user.ClientID == "streaming-api"
}

func (s *ServiceTokenValidator) HasServiceRole(user *ServiceUser, role string) bool {
	for _, userRole := range user.Roles {
		if userRole == role {
			return true
		}
	}
	return false
}
