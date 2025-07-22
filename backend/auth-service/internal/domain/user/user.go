package user

import (
	"time"

	pkgerrors "github.com/serdarburakguneri/hobby-streamer/backend/pkg/errors"
)

type User struct {
	id        UserID
	username  Username
	email     Email
	roles     Roles
	enabled   bool
	createdAt CreatedAt
	updatedAt UpdatedAt
	version   int
}

func NewUser(
	id UserID,
	username Username,
	email Email,
	roles Roles,
	enabled bool,
	createdAt CreatedAt,
	updatedAt UpdatedAt,
) *User {
	return &User{
		id:        id,
		username:  username,
		email:     email,
		roles:     roles,
		enabled:   enabled,
		createdAt: createdAt,
		updatedAt: updatedAt,
		version:   1,
	}
}

func (u *User) ID() UserID {
	return u.id
}

func (u *User) Username() Username {
	return u.username
}

func (u *User) Email() Email {
	return u.email
}

func (u *User) Roles() Roles {
	return u.roles
}

func (u *User) Enabled() bool {
	return u.enabled
}

func (u *User) CreatedAt() CreatedAt {
	return u.createdAt
}

func (u *User) UpdatedAt() UpdatedAt {
	return u.updatedAt
}

func (u *User) Version() int {
	return u.version
}

func (u *User) HasRole(role Role) bool {
	return u.roles.Contains(role)
}

func (u *User) HasAnyRole(roles []Role) bool {
	for _, role := range roles {
		if u.HasRole(role) {
			return true
		}
	}
	return false
}

func (u *User) HasAllRoles(roles []Role) bool {
	for _, role := range roles {
		if !u.HasRole(role) {
			return false
		}
	}
	return true
}

func (u *User) IsActive() bool {
	return u.enabled
}

func (u *User) AddRole(role Role) error {
	if u.HasRole(role) {
		return ErrRoleAlreadyExists
	}

	u.roles = *u.roles.Add(role)
	u.updatedAt = *NewUpdatedAt(time.Now())
	u.version++
	return nil
}

func (u *User) RemoveRole(role Role) error {
	if !u.HasRole(role) {
		return ErrRoleNotFound
	}

	u.roles = *u.roles.Remove(role)
	u.updatedAt = *NewUpdatedAt(time.Now())
	u.version++
	return nil
}

func (u *User) Enable() {
	if !u.enabled {
		u.enabled = true
		u.updatedAt = *NewUpdatedAt(time.Now())
		u.version++
	}
}

func (u *User) Disable() {
	if u.enabled {
		u.enabled = false
		u.updatedAt = *NewUpdatedAt(time.Now())
		u.version++
	}
}

func (u *User) UpdateEmail(email Email) error {
	if u.email.Equals(email) {
		return nil
	}

	u.email = email
	u.updatedAt = *NewUpdatedAt(time.Now())
	u.version++
	return nil
}

func (u *User) UpdateUsername(username Username) error {
	if u.username.Equals(username) {
		return nil
	}

	u.username = username
	u.updatedAt = *NewUpdatedAt(time.Now())
	u.version++
	return nil
}

var (
	ErrRoleAlreadyExists = pkgerrors.NewValidationError("role already exists", nil)
	ErrRoleNotFound      = pkgerrors.NewValidationError("role not found", nil)
)
