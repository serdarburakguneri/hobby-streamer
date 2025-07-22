package user

import (
	"regexp"
	"strings"
	"time"

	pkgerrors "github.com/serdarburakguneri/hobby-streamer/backend/pkg/errors"
)

type UserID struct {
	value string
}

func NewUserID(value string) (*UserID, error) {
	if value == "" {
		return nil, ErrInvalidUserID
	}
	return &UserID{value: value}, nil
}

func (id UserID) Value() string {
	return id.value
}

func (id UserID) Equals(other UserID) bool {
	return id.value == other.value
}

type Username struct {
	value string
}

func NewUsername(value string) (*Username, error) {
	if value == "" || len(value) < 3 || len(value) > 50 {
		return nil, ErrInvalidUsername
	}

	usernameRegex := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	if !usernameRegex.MatchString(value) {
		return nil, ErrInvalidUsername
	}

	return &Username{value: strings.ToLower(value)}, nil
}

func (u Username) Value() string {
	return u.value
}

func (u Username) Equals(other Username) bool {
	return u.value == other.value
}

type Email struct {
	value string
}

func NewEmail(value string) (*Email, error) {
	if value == "" {
		return nil, ErrInvalidEmail
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(value) {
		return nil, ErrInvalidEmail
	}

	return &Email{value: strings.ToLower(value)}, nil
}

func (e Email) Value() string {
	return e.value
}

func (e Email) Equals(other Email) bool {
	return e.value == other.value
}

type Role struct {
	value string
}

func NewRole(value string) (*Role, error) {
	if value == "" {
		return nil, ErrInvalidRole
	}

	roleRegex := regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)
	if !roleRegex.MatchString(value) {
		return nil, ErrInvalidRole
	}

	return &Role{value: value}, nil
}

func (r Role) Value() string {
	return r.value
}

func (r Role) Equals(other Role) bool {
	return r.value == other.value
}

type Roles struct {
	values []Role
}

func NewRoles(roles []string) (*Roles, error) {
	roleObjects := make([]Role, 0, len(roles))
	for _, roleStr := range roles {
		role, err := NewRole(roleStr)
		if err != nil {
			return nil, err
		}
		roleObjects = append(roleObjects, *role)
	}
	return &Roles{values: roleObjects}, nil
}

func (r Roles) Values() []Role {
	return r.values
}

func (r Roles) Contains(role Role) bool {
	for _, existingRole := range r.values {
		if existingRole.Equals(role) {
			return true
		}
	}
	return false
}

func (r Roles) Add(role Role) *Roles {
	if !r.Contains(role) {
		newRoles := make([]Role, len(r.values)+1)
		copy(newRoles, r.values)
		newRoles[len(r.values)] = role
		return &Roles{values: newRoles}
	}
	return &Roles{values: r.values}
}

func (r Roles) Remove(role Role) *Roles {
	newRoles := make([]Role, 0, len(r.values))
	for _, existingRole := range r.values {
		if !existingRole.Equals(role) {
			newRoles = append(newRoles, existingRole)
		}
	}
	return &Roles{values: newRoles}
}

type CreatedAt struct {
	value time.Time
}

func NewCreatedAt(value time.Time) *CreatedAt {
	return &CreatedAt{value: value}
}

func (c CreatedAt) Value() time.Time {
	return c.value
}

type UpdatedAt struct {
	value time.Time
}

func NewUpdatedAt(value time.Time) *UpdatedAt {
	return &UpdatedAt{value: value}
}

func (u UpdatedAt) Value() time.Time {
	return u.value
}

var (
	ErrInvalidUserID   = pkgerrors.NewValidationError("invalid user ID", nil)
	ErrInvalidUsername = pkgerrors.NewValidationError("invalid username", nil)
	ErrInvalidEmail    = pkgerrors.NewValidationError("invalid email", nil)
	ErrInvalidRole     = pkgerrors.NewValidationError("invalid role", nil)
)
