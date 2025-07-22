package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/serdarburakguneri/hobby-streamer/backend/auth-service/internal/domain/token"
	"github.com/serdarburakguneri/hobby-streamer/backend/auth-service/internal/domain/user"
	pkgerrors "github.com/serdarburakguneri/hobby-streamer/backend/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type TokenValidationService interface {
	ValidateToken(tokenString string) (*TokenValidationResult, error)
}

type TokenValidationResult struct {
	IsValid   bool
	User      *user.User
	Message   string
	ExpiresAt time.Time
}

type DomainTokenValidationService struct {
	keyFunc jwt.Keyfunc
}

func NewDomainTokenValidationService(keyFunc jwt.Keyfunc) *DomainTokenValidationService {
	return &DomainTokenValidationService{
		keyFunc: keyFunc,
	}
}

func (s *DomainTokenValidationService) ValidateToken(tokenString string) (*TokenValidationResult, error) {
	if tokenString == "" {
		return &TokenValidationResult{
			IsValid: false,
			Message: "Token is empty",
		}, nil
	}

	accessToken, err := token.NewAccessToken(tokenString)
	if err != nil {
		return &TokenValidationResult{
			IsValid: false,
			Message: "Invalid token format",
		}, nil
	}

	expiresAt, err := s.extractExpirationFromToken(accessToken.Value())
	if err != nil {
		return &TokenValidationResult{
			IsValid: false,
			Message: "Invalid token payload",
		}, nil
	}

	if time.Now().After(expiresAt) {
		return &TokenValidationResult{
			IsValid:   false,
			Message:   "Token expired",
			ExpiresAt: expiresAt,
		}, nil
	}

	user, err := s.extractUserFromToken(accessToken.Value())
	if err != nil {
		return &TokenValidationResult{
			IsValid:   false,
			Message:   "Invalid user data in token",
			ExpiresAt: expiresAt,
		}, nil
	}

	return &TokenValidationResult{
		IsValid:   true,
		User:      user,
		Message:   "",
		ExpiresAt: expiresAt,
	}, nil
}

func (s *DomainTokenValidationService) extractExpirationFromToken(tokenString string) (time.Time, error) {
	parsedToken, err := jwt.Parse(tokenString, s.keyFunc)
	if err != nil {
		return time.Time{}, pkgerrors.WithContext(err, map[string]interface{}{
			"operation": "extract_expiration",
			"error":     err.Error(),
		})
	}

	if !parsedToken.Valid {
		return time.Time{}, pkgerrors.NewValidationError("invalid token", nil)
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return time.Time{}, pkgerrors.NewValidationError("invalid token claims", nil)
	}

	expClaim, exists := claims["exp"]
	if !exists {
		return time.Time{}, pkgerrors.NewValidationError("token missing expiration claim", nil)
	}

	var expTime time.Time
	switch exp := expClaim.(type) {
	case float64:
		expTime = time.Unix(int64(exp), 0)
	case int64:
		expTime = time.Unix(exp, 0)
	case string:
		expInt, err := time.Parse(time.RFC3339, exp)
		if err != nil {
			return time.Time{}, pkgerrors.WithContext(err, map[string]interface{}{
				"operation": "parse_expiration_string",
				"exp_value": exp,
			})
		}
		expTime = expInt
	default:
		return time.Time{}, pkgerrors.NewValidationError("invalid expiration claim format", nil)
	}

	return expTime, nil
}

func (s *DomainTokenValidationService) extractUserFromToken(tokenString string) (*user.User, error) {
	parsedToken, err := jwt.Parse(tokenString, s.keyFunc)
	if err != nil {
		return nil, pkgerrors.WithContext(err, map[string]interface{}{
			"operation": "extract_user",
			"error":     err.Error(),
		})
	}

	if !parsedToken.Valid {
		return nil, pkgerrors.NewValidationError("invalid token", nil)
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, pkgerrors.NewValidationError("invalid token claims", nil)
	}

	userIDStr := s.getStringClaim(claims, "sub")
	if userIDStr == "" {
		return nil, pkgerrors.NewValidationError("token missing user ID", nil)
	}

	userID, err := user.NewUserID(userIDStr)
	if err != nil {
		return nil, pkgerrors.WithContext(err, map[string]interface{}{
			"operation": "create_user_id",
			"user_id":   userIDStr,
		})
	}

	usernameStr := s.getStringClaim(claims, "preferred_username")
	if usernameStr == "" {
		usernameStr = s.getStringClaim(claims, "username")
	}
	if usernameStr == "" {
		usernameStr = "unknown"
	}

	username, err := user.NewUsername(usernameStr)
	if err != nil {
		return nil, pkgerrors.WithContext(err, map[string]interface{}{
			"operation": "create_username",
			"username":  usernameStr,
		})
	}

	emailStr := s.getStringClaim(claims, "email")
	if emailStr == "" {
		emailStr = usernameStr + "@example.com"
	}

	email, err := user.NewEmail(emailStr)
	if err != nil {
		return nil, pkgerrors.WithContext(err, map[string]interface{}{
			"operation": "create_email",
			"email":     emailStr,
		})
	}

	var roleStrings []string
	if realmAccess, ok := claims["realm_access"].(map[string]interface{}); ok {
		if rolesInterface, ok := realmAccess["roles"].([]interface{}); ok {
			for _, role := range rolesInterface {
				if roleStr, ok := role.(string); ok {
					roleStrings = append(roleStrings, roleStr)
				}
			}
		}
	}

	if len(roleStrings) == 0 {
		roleStrings = []string{"user"}
	}

	roles, err := user.NewRoles(roleStrings)
	if err != nil {
		return nil, pkgerrors.WithContext(err, map[string]interface{}{
			"operation": "create_roles",
			"roles":     roleStrings,
		})
	}

	createdAt := user.NewCreatedAt(time.Now())
	updatedAt := user.NewUpdatedAt(time.Now())

	return user.NewUser(*userID, *username, *email, *roles, true, *createdAt, *updatedAt), nil
}

func (s *DomainTokenValidationService) getStringClaim(claims jwt.MapClaims, key string) string {
	if val, ok := claims[key].(string); ok {
		return val
	}
	return ""
}

type PasswordPolicyService interface {
	ValidatePassword(password string) error
	GeneratePasswordHash(password string) (string, error)
	VerifyPassword(password, hash string) bool
}

type DomainPasswordPolicyService struct{}

func NewDomainPasswordPolicyService() *DomainPasswordPolicyService {
	return &DomainPasswordPolicyService{}
}

func (s *DomainPasswordPolicyService) ValidatePassword(password string) error {
	if len(password) < 8 {
		return pkgerrors.NewValidationError("password too short", nil)
	}

	if len(password) > 128 {
		return pkgerrors.NewValidationError("password too long", nil)
	}

	hasUpper := false
	hasLower := false
	hasDigit := false
	hasSpecial := false

	for _, char := range password {
		switch {
		case char >= 'A' && char <= 'Z':
			hasUpper = true
		case char >= 'a' && char <= 'z':
			hasLower = true
		case char >= '0' && char <= '9':
			hasDigit = true
		case char >= 33 && char <= 47 || char >= 58 && char <= 64 || char >= 91 && char <= 96 || char >= 123 && char <= 126:
			hasSpecial = true
		}
	}

	if !hasUpper {
		return pkgerrors.NewValidationError("password missing uppercase letter", nil)
	}
	if !hasLower {
		return pkgerrors.NewValidationError("password missing lowercase letter", nil)
	}
	if !hasDigit {
		return pkgerrors.NewValidationError("password missing digit", nil)
	}
	if !hasSpecial {
		return pkgerrors.NewValidationError("password missing special character", nil)
	}

	return nil
}

func (s *DomainPasswordPolicyService) GeneratePasswordHash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", pkgerrors.WithContext(err, map[string]interface{}{
			"operation": "generate_password_hash",
		})
	}
	return string(hash), nil
}

func (s *DomainPasswordPolicyService) VerifyPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

var (
	ErrPasswordTooShort         = errors.New("password too short")
	ErrPasswordTooLong          = errors.New("password too long")
	ErrPasswordMissingUpperCase = errors.New("password missing uppercase letter")
	ErrPasswordMissingLowerCase = errors.New("password missing lowercase letter")
	ErrPasswordMissingDigit     = errors.New("password missing digit")
	ErrPasswordMissingSpecial   = errors.New("password missing special character")
)
