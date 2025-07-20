package user

import (
	"context"
)

type Repository interface {
	Save(ctx context.Context, user *User) error
	FindByID(ctx context.Context, id UserID) (*User, error)
	FindByUsername(ctx context.Context, username Username) (*User, error)
	FindByEmail(ctx context.Context, email Email) (*User, error)
	FindAll(ctx context.Context) ([]*User, error)
	Delete(ctx context.Context, id UserID) error
	ExistsByUsername(ctx context.Context, username Username) (bool, error)
	ExistsByEmail(ctx context.Context, email Email) (bool, error)
}
