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

func (c *Redis) GetByLogin(ctx context.Context, login string) (user *model.User, notCached bool, err error) {
	mycache := cache.New(&cache.Options{
		Redis: c.Client,
	})

	user = &model.User{}
	err = mycache.Get(ctx, login, user)
	if err != nil {
		if err.Error() == "cache: key is missing" {
			notCached = true
			return
		}
		err = fmt.Errorf("redis - GetByLogin - Get: %w", err)
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

func (c *Redis) ProduceUser(ctx context.Context, user *model.User) error {
	mu, _ := json.Marshal(user)
	_, err := c.Client.XAdd(ctx, &redis.XAddArgs{
		Stream: "example",
		Values: map[string]interface{}{
			"data": mu,
		},
	}).Result()
	if err != nil {
		return fmt.Errorf("redis - RedisStreamInit - XAdd: %w", err)
	}
	return nil
}

func (c *Redis) ConsumeUser(stream string) {
	go func() {
		for {
			var err error
			var data []redis.XMessage
			data, err = c.Client.XRangeN(context.Background(), stream, "-", "+", 10).Result()
			if err != nil {
				logrus.Error(err)
			}
			for _, element := range data {
				dataFromStream := []byte(element.Values["data"].(string))
				var user = &model.User{}
				err := json.Unmarshal(dataFromStream, user)
				if err != nil {
					logrus.Error(err)
					continue
				}
				logrus.Infof("user created:%v", user.Login)
				c.Client.XDel(context.Background(), stream, element.ID)
			}
		}
	}()
}
