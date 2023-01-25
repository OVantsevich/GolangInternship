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

func (r *PRepository) CreateEntity(ctx context.Context, e *Entity) error {
	var name string
	err := r.Pool.QueryRow(ctx,
		"INSERT INTO entity (name, age) SELECT $1, $2 WHERE NOT EXISTS(SELECT 1 FROM entity WHERE name=$3) RETURNING name",
		e.Name, e.Age, e.Name).Scan(&name)
	if err != nil {
		return fmt.Errorf("repository - PRepository - CreateEntity: %v", err)
	}

	return nil
}

func (r *PRepository) GetEntityByName(ctx context.Context, name string) (*Entity, error) {
	e := Entity{}
	err := r.Pool.QueryRow(ctx, "select * from entity where name=$1 and not deleted", name).Scan(
		&e.ID, &e.Name, &e.Age, &e.Deleted)
	if err != nil {
		return nil, fmt.Errorf("repository - PRepository - GetEntityByName: %v", err)
	}

	return &e, nil
}
func (r *PRepository) UpdateEntity(ctx context.Context, name string, e *Entity) error {
	var id int
	err := r.Pool.QueryRow(ctx, "update entity set age=$1 where name=$2 and deleted=false returning id",
		e.Age, name).Scan(&id)
	if err != nil {
		return fmt.Errorf("repository - PRepository - CreateEntity: %v", err)
	}

	return nil
}
func (r *PRepository) DeleteEntity(ctx context.Context, name string) error {
	var id int
	err := r.Pool.QueryRow(ctx, "update entity set deleted=true where name=$1 and deleted=false returning id",
		name).Scan(&id)
	if err != nil {
		return fmt.Errorf("repository - PRepository - CreateEntity: %v", err)
	}

	return nil
}
