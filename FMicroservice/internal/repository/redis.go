package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/OVantsevich/GolangInternship/FMicroservice/internal/model"
	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
)

type Redis struct {
	Client redis.Client
}

func (c *Redis) GetByLogin(ctx context.Context, login string) (user *model.User, err error) {
	mycache := cache.New(&cache.Options{
		Redis: c.Client,
	})

	user = &model.User{}
	err = mycache.Get(ctx, login, user)
	if err != nil {
		if err.Error() == "cache: key is missing" {
			return nil, nil
		}
		return nil, fmt.Errorf("redis - GetByLogin - Get: %w", err)
	}

	return
}

func (c *Redis) CreateUser(ctx context.Context, user *model.User) error {
	mycache := cache.New(&cache.Options{
		Redis: c.Client,
	})

	err := mycache.Set(&cache.Item{
		Ctx:   ctx,
		Key:   user.Login,
		Value: user,
	})
	if err != nil {
		return fmt.Errorf("redis - CreateUser - Set: %w", err)
	}
	return nil
}

func (c *Redis) RedisStreamInit(ctx context.Context) error {
	_, err := c.Client.XAdd(ctx, &redis.XAddArgs{
		Stream: "example",
	}).Result()
	if err != nil {
		return fmt.Errorf("redis - RedisStreamInit - XAdd: %w", err)
	}

	_, err = c.Client.XGroupCreate(ctx, "example", "user", "").Result()
	if err != nil {
		return fmt.Errorf("redis - RedisStreamInit - XGroupCreate: %w", err)
	}

	return nil
}

func (c *Redis) CreatingUser(ctx context.Context, user *model.User) error {
	mu, _ := json.Marshal(user)
	_, err := c.Client.XAdd(ctx, &redis.XAddArgs{
		Stream: "example",
		Values: map[string]interface{}{
			"user": mu,
		},
	}).Result()
	if err != nil {
		return fmt.Errorf("redis - RedisStreamInit - XAdd: %w", err)
	}
	streams := c.Client.XReadStreams(ctx, "example").Val()
	for _, str := range streams {
		for _, msg := range str.Messages {
			for _, val := range msg.Values {
				logrus.Infof("user created: %v", val)
			}
		}
	}
	return nil
}

func (c *Redis) CreatConsumer() error {
	go c.UserCreatHandler("example")
	return nil
}

func (c *Redis) UserCreatHandler(stream string) {
	ctx := context.Background()
	for {
		streams := c.Client.XReadStreams(ctx, stream).Val()
		for _, str := range streams {
			for _, msg := range str.Messages {
				for _, val := range msg.Values {
					logrus.Infof("user created: %v", val)
				}
			}
		}
	}
}
