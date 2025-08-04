package auth

import (
	"context"
	"crypto/rsa"
	"crypto/tls"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/constants"
	"github.com/serdarburakguneri/hobby-streamer/backend/pkg/errors"
)

type ServiceTokenValidator struct {
	keycloakURL string
	realm       string
	clientID    string
	httpClient  *http.Client
	jwks        *jwksCache
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
		jwks:        newJWKSCache(fmt.Sprintf("%s/realms/%s/protocol/openid-connect/certs", keycloakURL, realm)),
	}
}

func (s *ServiceTokenValidator) ValidateServiceToken(ctx context.Context, token string) (*ServiceUser, error) {
	token = strings.TrimPrefix(token, constants.BearerPrefix)

	parsedToken, _, err := jwt.NewParser().ParseUnverified(token, jwt.MapClaims{})
	if err != nil {
		return nil, errors.NewInternalError("invalid token format", err)
	}

	header, ok := parsedToken.Header["kid"].(string)
	if !ok {
		return nil, errors.NewInternalError("missing key ID in token header", nil)
	}

	publicKey, err := s.getPublicKey(ctx, header)
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
	return s.jwks.getKey(ctx, kid)
}

func (s *ServiceTokenValidator) validateServiceClaims(claims jwt.MapClaims) error {
	if exp, ok := claims["exp"].(float64); ok {
		if time.Unix(int64(exp), 0).Before(time.Now()) {
			return errors.NewInternalError(constants.ErrTokenExpired, nil)
		}
	}

	if iss, ok := claims["iss"].(string); ok {
		expectedIss := fmt.Sprintf("%s/realms/%s", s.keycloakURL, s.realm)
		if iss != expectedIss {
			return errors.NewInternalError(fmt.Sprintf("invalid issuer: expected %s, got %s", expectedIss, iss), nil)
		}
	}

	if aud, ok := claims["aud"].(string); ok {
		if aud != s.clientID {
			return errors.NewInternalError(fmt.Sprintf("invalid audience: expected %s, got %s", s.clientID, aud), nil)
		}
	}

	tokenType, ok := claims["typ"].(string)
	if !ok || tokenType != "Bearer" {
		return errors.NewInternalError("invalid token type: expected Bearer", nil)
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
