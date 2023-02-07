package handler

import (
	"GolangInternship/FMicroserviceGRPC/internal/repository"
	"GolangInternship/FMicroserviceGRPC/internal/service"
	pr "GolangInternship/FMicroserviceGRPC/proto"
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
	"time"
)

var (
	TestPgUser          = "postgres"
	TestPgPassword      = "postgres"
	TestPgDB            = "postgres"
	TestPgPort          = "11111"
	TestPgContainerName = "postgres"
)

var MongoURL = "mongodb://mongo:mongo@localhost:27017"
var JwtKey = "testJWTKey"

var handlerTest *UserClassic
var postgresPool *pgxpool.Pool

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("unix:///home/olegvantsevich/.docker/desktop/docker.sock")
	if err != nil {
		logrus.Fatalf("Could not construct pool: %s", err)
	}
	err = pool.Client.Ping()
	if err != nil {
		logrus.Fatalf("Could not connect to Docker: %s", err)
	}

	network, err := pool.Client.CreateNetwork(docker.CreateNetworkOptions{Name: "testN"})
	if err != nil {
		logrus.Fatalf("Could not create network: %s", err)
	}

	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Name:       TestPgContainerName,
		Repository: "postgres",
		Tag:        "latest",
		NetworkID:  network.ID,
		Env: []string{
			fmt.Sprintf("POSTGRES_USER=%s", TestPgUser),
			fmt.Sprintf("POSTGRES_PASSWORD=%s", TestPgPassword),
			fmt.Sprintf("POSTGRES_DB=%s", TestPgDB),
		},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"5432/tcp": {{HostIP: "localhost", HostPort: fmt.Sprintf("%s/tcp", TestPgPort)}},
		}})
	if err != nil {
		logrus.Fatalf("Could not start postgres: %s", err)
	}

	flyway, err := pool.RunWithOptions(&dockertest.RunOptions{
		Name:       "flyway",
		Repository: "flyway/flyway",
		Tag:        "latest",
		NetworkID:  network.ID,
		Cmd: []string{
			fmt.Sprintf("-user=%s", TestPgUser),
			fmt.Sprintf("-password=%s", TestPgPassword),
			fmt.Sprintf("-url=jdbc:postgresql://%s:%s/%s", TestPgContainerName, "5432", TestPgDB),
			"-locations=filesystem:/flyway/sql",
			"-connectRetries=60",
			"migrate",
		},
		Mounts: []string{"/home/olegvantsevich/GolandProjects/GolangInternship/FMicroservice/migrations:/flyway/sql"},
	})
	if err != nil {
		logrus.Fatalf("Could not start flyway: %s", err)
	}

	time.Sleep(time.Second * 10)

	postgresPool, err = pgxpool.New(context.Background(), fmt.Sprintf("postgres://%s:%s@localhost:%s/%s", TestPgUser, TestPgPassword, TestPgPort, TestPgDB))
	if err != nil {
		logrus.Fatalf("Could not connect to db %s", err)
	}
	defer postgresPool.Close()

	cache, err := pool.RunWithOptions(&dockertest.RunOptions{
		Name:       "redis-test",
		Repository: "redis",
		Tag:        "latest",
		NetworkID:  network.ID,
	})
	if err != nil {
		logrus.Fatalf("Could not start redis: %s", err)
	}

	client := redis.NewClient(&redis.Options{
		Addr: "localhost:" + cache.GetPort("6379/tcp"),
	})
	defer client.Close()

	rds := &repository.Redis{Client: *client}

	rds.RedisStreamInit(context.Background())
	rds.ConsumeUser("example")

	userService := service.NewUserServiceClassic(repository.NewPostgresRepository(postgresPool), rds, rds, JwtKey)
	handlerTest = NewUserHandlerClassic(userService, JwtKey)

	code := m.Run()

	if err := pool.Purge(resource); err != nil {
		logrus.Fatalf("Could not purge postgres: %s", err)
	}

	if err := pool.Purge(flyway); err != nil {
		logrus.Fatalf("Could not purge flyway: %s", err)
	}

	if err := pool.Purge(cache); err != nil {
		logrus.Fatalf("Could not purge cache: %s", err)
	}

	if err := pool.Client.RemoveNetwork(network.ID); err != nil {
		logrus.Fatalf("Could not remove network: %s", err)
	}

	os.Exit(code)
}

var testSignUpValid = []pr.SignupRequest{
	{
		Name:     `NameTest`,
		Age:      99,
		Login:    `TestLogin`,
		Email:    `login@test.com`,
		Password: `strongTestPassword`,
	},
}

var testSignUpInvalid = []pr.SignupRequest{
	{
		Name:     `NameTest`,
		Age:      101,
		Login:    `TestLogin`,
		Email:    `login@test.com`,
		Password: `strongTestPassword`,
	},
	{
		Name:     `NameTest`,
		Age:      0,
		Login:    `TestLogin`,
		Email:    `login@test.com`,
		Password: `strongTestPassword`,
	},
	{
		Name:     `NameTest`,
		Age:      99,
		Login:    `TestLogin`,
		Email:    `login`,
		Password: `strongTestPassword`,
	},
	{
		Name:     `NameTest`,
		Age:      99,
		Login:    `TestLogin`,
		Email:    `0test.com`,
		Password: `strongTestPassword`,
	},
	{
		Name:     `NameTest`,
		Age:      99,
		Email:    `login@test.com`,
		Password: `strongTestPassword`,
	},
	{
		Age:      99,
		Login:    `TestLogin`,
		Email:    `login@test.com`,
		Password: `strongTestPassword`,
	},
	{
		Name:     `NameTest`,
		Age:      99,
		Login:    `TestLogin`,
		Password: `strongTestPassword`,
	},
	{
		Name:  `NameTest`,
		Age:   99,
		Login: `TestLogin`,
		Email: `login@test.com`,
	},
	{
		Name:     `NameTest`,
		Age:      99,
		Login:    `TestLogin`,
		Email:    `login@test.com`,
		Password: `weakpass`,
	},
	{
		Name:     `NameTest`,
		Age:      99,
		Login:    `TestLogin`,
		Email:    `login@test.com`,
		Password: `111`,
	},
	{
		Name:     `NameTest`,
		Age:      99,
		Login:    `TestLogin`,
		Email:    `login@test.com`,
		Password: `1234567890`,
	},
}

func TestUserClassic_Signup(t *testing.T) {
	var response *pr.SignupResponse
	var err error

	for _, user := range testSignUpValid {
		_, err = postgresPool.Exec(context.Background(), "delete from users where login=$1", user.Login)
		require.NoError(t, err)

		response, err = handlerTest.Signup(context.Background(), &user)
		require.NoError(t, err)

		_, err = handlerTest.Verify(response.AccessToken)
		require.NoError(t, err)
	}

	for _, user := range testSignUpInvalid {
		response, err = handlerTest.Signup(context.Background(), &user)
		require.Error(t, err)
	}

	for _, user := range testSignUpValid {
		_, err = postgresPool.Exec(context.Background(), "delete from users where login=$1", user.Login)
		require.NoError(t, err)

		response, err = handlerTest.Signup(context.Background(), &user)
		require.NoError(t, err)

		response, err = handlerTest.Signup(context.Background(), &user)
		require.Error(t, err)
	}
}

func TestUserClassic_Login(t *testing.T) {
	var response *pr.LoginResponse
	var err error

	for _, user := range testSignUpValid {
		_, err = postgresPool.Exec(context.Background(), "delete from users where login=$1", user.Login)
		require.NoError(t, err)

		_, err = handlerTest.Signup(context.Background(), &user)
		require.NoError(t, err)

		response, err = handlerTest.Login(context.Background(), &pr.LoginRequest{
			Login:    user.Login,
			Password: user.Password,
		})
		require.NoError(t, err)

		_, err = handlerTest.Verify(response.AccessToken)
		require.NoError(t, err)
	}

	for _, user := range testSignUpValid {
		_, err = postgresPool.Exec(context.Background(), "delete from users where login=$1", user.Login)
		require.NoError(t, err)

		response, err = handlerTest.Login(context.Background(), &pr.LoginRequest{
			Login:    user.Login,
			Password: user.Password,
		})
		require.Error(t, err)
	}

	for _, user := range testSignUpValid {
		_, err = postgresPool.Exec(context.Background(), "delete from users where login=$1", user.Login)
		require.NoError(t, err)

		_, err = handlerTest.Signup(context.Background(), &user)
		require.NoError(t, err)

		response, err = handlerTest.Login(context.Background(), &pr.LoginRequest{
			Login:    user.Login,
			Password: "wrong password",
		})
		require.Error(t, err)
	}
}

func TestUserClassic_Refresh(t *testing.T) {
	var response *pr.RefreshResponse
	var signupResponse *pr.SignupResponse
	var err error

	for _, user := range testSignUpValid {
		_, err = postgresPool.Exec(context.Background(), "delete from users where login=$1", user.Login)
		require.NoError(t, err)

		signupResponse, err = handlerTest.Signup(context.Background(), &user)
		require.NoError(t, err)

		response, err = handlerTest.Refresh(context.Background(), &pr.RefreshRequest{
			Login:        user.Login,
			RefreshToken: signupResponse.RefreshToken,
		})
		require.NoError(t, err)

		_, err = handlerTest.Verify(response.AccessToken)
		require.NoError(t, err)
	}

	for _, user := range testSignUpValid {
		_, err = postgresPool.Exec(context.Background(), "delete from users where login=$1", user.Login)
		require.NoError(t, err)

		response, err = handlerTest.Refresh(context.Background(), &pr.RefreshRequest{
			Login:        user.Login,
			RefreshToken: "",
		})
		require.Error(t, err)
	}

	for _, user := range testSignUpValid {
		_, err = postgresPool.Exec(context.Background(), "delete from users where login=$1", user.Login)
		require.NoError(t, err)

		_, err = handlerTest.Signup(context.Background(), &user)
		require.NoError(t, err)

		response, err = handlerTest.Refresh(context.Background(), &pr.RefreshRequest{
			Login:        user.Login,
			RefreshToken: "ad",
		})
		require.Error(t, err)
	}
}
func TestUserClassic_Update(t *testing.T) {
	var response *pr.UpdateResponse
	var signupResponse *pr.SignupResponse
	var err error

	for _, user := range testSignUpValid {
		_, err = postgresPool.Exec(context.Background(), "delete from users where login=$1", user.Login)
		require.NoError(t, err)

		signupResponse, err = handlerTest.Signup(context.Background(), &user)
		require.NoError(t, err)

		response, err = handlerTest.Update(context.Background(), &pr.UpdateRequest{
			User: &pr.User{
				Email: user.Email,
				Name:  user.Name,
				Age:   user.Age,
			},
			AccessToken: signupResponse.AccessToken,
		})
		require.NoError(t, err)
		require.Equal(t, response.Login, user.Login)
	}

	for i, user := range testSignUpInvalid {
		if i == 4 {
			break
		}
		_, err = postgresPool.Exec(context.Background(), "delete from users where login=$1", testSignUpValid[0].Login)
		require.NoError(t, err)

		signupResponse, err = handlerTest.Signup(context.Background(), &testSignUpValid[0])
		require.NoError(t, err)

		response, err = handlerTest.Update(context.Background(), &pr.UpdateRequest{
			User: &pr.User{
				Email: user.Email,
				Name:  user.Name,
				Age:   user.Age,
			},
			AccessToken: signupResponse.AccessToken,
		})
		require.Error(t, err)
	}

}
