package repository

import (
	"context"
	"fmt"
	"github.com/OVantsevich/GolangInternship/FMicroservice/internal/model"
	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
)

type RedisCache struct {
	Client redis.Client
}

func (c *RedisCache) GetByLogin(ctx context.Context, login string) (user *model.User, err error) {
	mycache := cache.New(&cache.Options{
		Redis: c.Client,
	})

	user = &model.User{}
	err = mycache.Get(ctx, login, user)
	if err != nil {
		if err.Error() == "cache: key is missing" {
			return nil, nil
		}
		return nil, fmt.Errorf("RedisCache - GetByLogin - Get: %w", err)
	}

	return
}

func (c *RedisCache) CreateUser(ctx context.Context, user *model.User) error {
	mycache := cache.New(&cache.Options{
		Redis: c.Client,
	})

	err := mycache.Set(&cache.Item{
		Ctx:   ctx,
		Key:   user.Login,
		Value: user,
	})
	if err != nil {
		return fmt.Errorf("RedisCache - CreateUser - Set: %w", err)
	}
	return nil
}
