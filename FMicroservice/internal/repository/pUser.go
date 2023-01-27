package repository

import (
	"context"
	"fmt"
	"github.com/OVantsevich/GolangInternship/FMicroservice/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type PUser struct {
	Pool *pgxpool.Pool
}

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

func (r *PUser) GetUserByLogin(ctx context.Context, login string) (*model.User, error) {
	user := model.User{}
	err := r.Pool.QueryRow(ctx, "select * from users where login=$1 and deleted=false", login).Scan(
		&user.ID, &user.Login, &user.Email, &user.Password, &user.Name, &user.Age, &user.Token, &user.Deleted, &user.Created, &user.Updated)
	if err != nil {
		return nil, fmt.Errorf("PUser - GetUserByName - QueryRow: %w", err)
	}

	return &user, nil
}
func (r *PUser) UpdateUser(ctx context.Context, login string, user *model.User) error {
	_, err := r.Pool.Exec(ctx, "update users set email=$1, name=$2, age=$3, updated=$4 where login=$5 and deleted=false",
		user.Email, user.Name, user.Age, user.Updated, login)
	if err != nil {
		return fmt.Errorf("PUser - UpdateUser - Exec: %w", err)
	}

	return nil
}

func (r *PUser) RefreshUser(ctx context.Context, login, token string) error {
	_, err := r.Pool.Exec(ctx, "update users set token=$1, updated=$2 where login=$3 and deleted=false",
		token, time.Now(), login)
	if err != nil {
		return fmt.Errorf("PUser - RefreshUser - Exec: %w", err)
	}

	return nil
}

func (r *PUser) DeleteUser(ctx context.Context, login string) error {
	_, err := r.Pool.Exec(ctx, "update users set deleted=true, updated=$1 where name=$2 and deleted==false",
		time.Now(), login)
	if err != nil {
		return fmt.Errorf("PUser - DeleteUser - Exec: %w", err)
	}

	return nil
}
