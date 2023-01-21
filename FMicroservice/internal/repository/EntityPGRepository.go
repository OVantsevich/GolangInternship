package repository

import (
	"context"
	"fmt"
	"github.com/OVantsevich/GolangInternship/FMicroservice/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
	log "github.com/sirupsen/logrus"
)

type PRepository struct {
	pool *pgxpool.Pool
}

func (r *PRepository) OpenPool(ctx context.Context) error {
	if r.pool == nil {
		var err error
		r.pool, err = pgxpool.New(ctx, domain.Cfg.PostgresUrl)
		if err != nil {
			log.Fatalf("database connection error: %v", err)
			return err
		}
		if r.pool.Ping(ctx) != nil {
			log.Fatalf("database connection error: database not responding")
			return fmt.Errorf("database not responding")
		}
	}
	if r.pool.Ping(ctx) != nil {
		log.Errorf("database connection error: database not responding")
		return fmt.Errorf("database not responding")
	}
	return nil
}

func (r *PRepository) ClosePool() {
	if r.pool != nil {
		r.pool.Close()
	}
}

func (r *PRepository) CreateEntity(ctx context.Context, e *domain.Entity) error {
	if err := r.OpenPool(ctx); err != nil {
		return err
	}

	var name string
	if err := r.pool.QueryRow(ctx,
		"INSERT INTO entity (name, age) SELECT $1, $2 WHERE NOT EXISTS(SELECT 1 FROM entity WHERE name=$3) RETURNING name",
		e.Name, e.Age, e.Name).Scan(&name); err != nil {
		log.Errorf("database error while creating entity: %v", err)
		return fmt.Errorf("entity with this name already exist")
	}

	return nil
}
func (r *PRepository) GetEntityByName(ctx context.Context, name string) (*domain.Entity, error) {
	if err := r.OpenPool(ctx); err != nil {
		return nil, err
	}

	e := domain.Entity{}
	if err := r.pool.QueryRow(ctx, "select * from entity where name=$1 and not is_deleted", name).Scan(
		&e.ID, &e.Name, &e.Age, &e.IsDeleted); err != nil {
		log.Errorf("database error while getting entity: %v", err)
		return nil, fmt.Errorf("entity with this name doesn't exist")
	}

	return &e, nil
}
func (r *PRepository) UpdateEntity(ctx context.Context, name string, e *domain.Entity) error {
	if err := r.OpenPool(ctx); err != nil {
		return err
	}

	var id int
	if err := r.pool.QueryRow(ctx, "update entity set age=$1 where name=$2 and is_deleted=false returning id",
		e.Age, name).Scan(&id); err != nil {
		log.Errorf("database error while updating entity: %v", err)
		return fmt.Errorf("entity with this name doesn't exist")
	}

	return nil
}
func (r *PRepository) DeleteEntity(ctx context.Context, name string) error {
	if err := r.OpenPool(ctx); err != nil {
		return err
	}

	var id int
	if err := r.pool.QueryRow(ctx, "update entity set is_deleted=true where name=$1 and is_deleted=false returning id",
		name).Scan(&id); err != nil {
		log.Errorf("database error while deleting an entity: %v", err)
		return fmt.Errorf("entity with this name doesn't exist")
	}

	return nil
}
