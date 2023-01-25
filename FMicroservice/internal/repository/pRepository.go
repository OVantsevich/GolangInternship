package repository

import (
	"context"
	"fmt"
	. "github.com/OVantsevich/GolangInternship/FMicroservice/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PRepository struct {
	Pool *pgxpool.Pool
}

func (r *PRepository) CreateUser(ctx context.Context, user *User) error {
	var name string
	err := r.Pool.QueryRow(ctx,
		"insert into user (login, email, password, name, age) select $1, $2, $3, $4, $5 returning name",
		user.Login, user.Email, user.Password, user.Name, user.Age).Scan(&name)
	if err != nil {
		return fmt.Errorf("repository - PRepository - CreateUser: %v", err)
	}

	return nil
}

func (r *PRepository) GetUserByName(ctx context.Context, name string) (*User, error) {
	user := User{}
	err := r.Pool.QueryRow(ctx, "select * from user where name=$1 and not deleted", name).Scan(
		&user.ID, &user.Login, &user.Email, &user.Password, &user.Name, &user.Age, &user.Deleted)
	if err != nil {
		return nil, fmt.Errorf("repository - PRepository - GetUserByName: %v", err)
	}

	return &user, nil
}
func (r *PRepository) UpdateUser(ctx context.Context, name string, user *User) error {
	var id int
	err := r.Pool.QueryRow(ctx, "update user set email=$1, name=$2, age=$3 where login=$4 and not deleted returning id",
		user.Age, name).Scan(&id)
	if err != nil {
		return fmt.Errorf("repository - PRepository - CreateUser: %v", err)
	}

	return nil
}
func (r *PRepository) DeleteUser(ctx context.Context, name string) error {
	var id int
	err := r.Pool.QueryRow(ctx, "update user set deleted=true where name=$1 and deleted=false returning id",
		name).Scan(&id)
	if err != nil {
		return fmt.Errorf("repository - PRepository - CreateUser: %v", err)
	}

	return nil
}
