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

func NewPostgresRepository(pool *pgxpool.Pool) *PUser {
	return &PUser{Pool: pool}
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
	err := r.Pool.QueryRow(ctx, "select (select name from roles where id = (select role_id from l_role_user where user_id = u.id)) role_name, "+
		"* "+
		"	from users u "+
		"where login=$1 "+
		"and deleted=false", login).Scan(
		&user.Role, &user.ID, &user.Login, &user.Email, &user.Password, &user.Name, &user.Age, &user.Token, &user.Deleted, &user.Created, &user.Updated)
	if err != nil {
		return nil, fmt.Errorf("PUser - GetUserByName - QueryRow: %w", err)
	}

	return &user, nil
}
func (r *PUser) UpdateUser(ctx context.Context, login string, user *model.User) error {
	var id int
	err := r.Pool.QueryRow(ctx, "update users set email=$1, name=$2, age=$3, updated=$4 where login=$5 and deleted=false returning id",
		user.Email, user.Name, user.Age, user.Updated, login).Scan(&id)
	if err != nil {
		return fmt.Errorf("PUser - UpdateUser - Exec: %w", err)
	}

	return nil
}

func (r *PUser) RefreshUser(ctx context.Context, login, token string) error {
	var id int
	err := r.Pool.QueryRow(ctx, "update users set token=$1, updated=$2 where login=$3 and deleted=false returning id",
		token, time.Now(), login).Scan(&id)
	if err != nil {
		return fmt.Errorf("PUser - RefreshUser - Exec: %w", err)
	}

	return nil
}

func (r *PUser) DeleteUser(ctx context.Context, login string) error {
	var id int
	err := r.Pool.QueryRow(ctx, "update users set deleted=true, updated=$1 where login=$2 and deleted=false returning id",
		time.Now(), login).Scan(&id)
	if err != nil {
		return fmt.Errorf("PUser - DeleteUser - Exec: %w", err)
	}

	return nil
}
