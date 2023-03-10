// Package repository pUser
package repository

import (
	"context"
	"fmt"
	"time"

	"GolangInternship/FMicroservice/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

// PUser mongo entity
type PUser struct {
	Pool *pgxpool.Pool
}

// NewPostgresRepository creating new PUser
func NewPostgresRepository(pool *pgxpool.Pool) *PUser {
	return &PUser{Pool: pool}
}

// CreateUser create user
func (r *PUser) CreateUser(ctx context.Context, user *model.User) (*model.User, error) {
	user.Created = time.Now()
	user.Updated = time.Now()
	_, err := r.Pool.Exec(ctx,
		"insert into users (login, email, password, name, age) values ($1, $2, $3, $4, $5);",
		user.Login, user.Email, user.Password, user.Name, user.Age)
	if err != nil {
		return nil, fmt.Errorf("PUser - CreateUser - Exec: %w", err)
	}

	return user, nil
}

// GetUserByLogin get user by login
func (r *PUser) GetUserByLogin(ctx context.Context, login string) (*model.User, error) {
	user := model.User{}
	err := r.Pool.QueryRow(ctx, `select u.name, u.age, u.login, u.password, u.token,u.email, r.name
									from users u
											 join l_role_user l on u.id = l.user_id
											 join roles r on l.role_id = r.id
									where u.login = $1 and u.deleted=false`, login).Scan(
		&user.Name, &user.Age, &user.Login, &user.Password, &user.Token, &user.Email, &user.Role)
	if err != nil {
		return nil, fmt.Errorf("PUser - GetUserByName - QueryRow: %w", err)
	}

	return &user, nil
}

// UpdateUser update user
func (r *PUser) UpdateUser(ctx context.Context, login string, user *model.User) error {
	var id int
	err := r.Pool.QueryRow(ctx, "update users set email=$1, name=$2, age=$3, updated=$4 where login=$5 and Deleted=false returning id",
		user.Email, user.Name, user.Age, user.Updated, login).Scan(&id)
	if err != nil {
		return fmt.Errorf("PUser - UpdateUser - Exec: %w", err)
	}

	return nil
}

// RefreshUser refresh user
func (r *PUser) RefreshUser(ctx context.Context, login, token string) error {
	var id int
	err := r.Pool.QueryRow(ctx, "update users set token=$1, updated=$2 where login=$3 and Deleted=false returning id",
		token, time.Now(), login).Scan(&id)
	if err != nil {
		return fmt.Errorf("PUser - RefreshUser - Exec: %w", err)
	}

	return nil
}

// DeleteUser delete user
func (r *PUser) DeleteUser(ctx context.Context, login string) error {
	var id int
	err := r.Pool.QueryRow(ctx, "update users set Deleted=true, updated=$1 where login=$2 and Deleted=false returning id",
		time.Now(), login).Scan(&id)
	if err != nil {
		return fmt.Errorf("PUser - DeleteUser - Exec: %w", err)
	}

	return nil
}
