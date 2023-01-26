package repository

import (
	"context"
	"fmt"
	. "github.com/OVantsevich/GolangInternship/FMicroservice/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type PRepository struct {
	Pool *pgxpool.Pool
}

func (r *PRepository) CreateUser(ctx context.Context, user *User) error {
	var id string
	user.Created = time.Now()
	user.Updated = time.Now()
	err := r.Pool.QueryRow(ctx,
		"insert into users (login, email, password, name, age) select $1, $2, $3, $4, $5 returning name",
		user.Login, user.Email, user.Password, user.Name, user.Age).Scan(&id)
	if err != nil {
		return fmt.Errorf("repository - PRepository - CreateUser: %v", err)
	}

	return nil
}

func (r *PRepository) GetUserByLogin(ctx context.Context, login string) (*User, error) {
	user := User{}
	err := r.Pool.QueryRow(ctx, "select * from users where login=$1 and not deleted", login).Scan(
		&user.ID, &user.Login, &user.Email, &user.Password, &user.Name, &user.Age, &user.Token, &user.Deleted, &user.Created, &user.Updated)
	if err != nil {
		return nil, fmt.Errorf("repository - PRepository - GetUserByName: %v", err)
	}

	return &user, nil
}
func (r *PRepository) UpdateUser(ctx context.Context, login string, user *User) error {
	var id int
	err := r.Pool.QueryRow(ctx, "update users set email=$1, name=$2, age=$3, updated=$4 where login=$5 and not deleted returning id",
		user.Email, user.Name, user.Age, user.Updated, login).Scan(&id)
	if err != nil {
		return fmt.Errorf("repository - PRepository - CreateUser: %v", err)
	}

	return nil
}

func (r *PRepository) RefreshUser(ctx context.Context, login, token string) error {
	var id int
	err := r.Pool.QueryRow(ctx, "update users set token=$1, updated=$2 where login=$3 and not deleted returning id",
		token, time.Now(), login).Scan(&id)
	if err != nil {
		return fmt.Errorf("repository - PRepository - RefreshUser: %v", err)
	}

	return nil
}

func (r *PRepository) DeleteUser(ctx context.Context, login string) error {
	var id int
	err := r.Pool.QueryRow(ctx, "update users set deleted=true, updated=$1 where name=$2 and deleted=false returning id",
		time.Now(), login).Scan(&id)
	if err != nil {
		return fmt.Errorf("repository - PRepository - CreateUser: %v", err)
	}

	return nil
}
