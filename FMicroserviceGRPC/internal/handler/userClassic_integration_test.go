package handler

import (
	"GolangInternship/FMicroserviceGRPC/internal/repository"
	"GolangInternship/FMicroserviceGRPC/internal/service"
	pr "GolangInternship/FMicroserviceGRPC/proto"
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt/v4"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
	"math"
	"net"
	"os"
	"path/filepath"
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
	defer pool.Client.RemoveNetwork(network.ID)

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
	defer pool.Purge(resource)

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
	defer pool.Purge(flyway)

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
	defer pool.Purge(cache)

	client := redis.NewClient(&redis.Options{
		Addr: "localhost:" + cache.GetPort("6379/tcp"),
	})
	defer client.Close()

	rds := &repository.Redis{Client: *client}

	rds.ConsumeUser("example")

	userService := service.NewUserServiceClassic(repository.NewPostgresRepository(postgresPool), rds, rds, JwtKey)
	fileService := service.NewFile("../../fileStore")
	handlerTest = NewUserHandlerClassic(userService, fileService, JwtKey)

	code := m.Run()

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

func verify(token string) (claims *service.CustomClaims, err error) {
	claims = &service.CustomClaims{}

	_, err = jwt.ParseWithClaims(
		token,
		claims,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(JwtKey), nil
		},
	)
	if err != nil {
		err = fmt.Errorf("invalid token: %w", err)
	}

	return
}

func TestUserClassic_Signup(t *testing.T) {
	var response *pr.SignupResponse
	var err error

	for _, user := range testSignUpValid { //nolint:govet //all ok
		_, err = postgresPool.Exec(context.Background(), "delete from users where login=$1", user.Login)
		require.NoError(t, err)

		response, err = handlerTest.Signup(context.Background(), &user)
		require.NoError(t, err)

		_, err = verify(response.AccessToken)
		require.NoError(t, err)
	}

	for _, user := range testSignUpInvalid { //nolint:govet //all ok
		_, err = handlerTest.Signup(context.Background(), &user)
		require.Error(t, err)
	}

	_, err = postgresPool.Exec(context.Background(), "delete from users where login=$1", testSignUpValid[0].Login)
	require.NoError(t, err)

	_, err = handlerTest.Signup(context.Background(), &testSignUpValid[0])
	require.NoError(t, err)

	_, err = handlerTest.Signup(context.Background(), &testSignUpValid[0])
	require.Error(t, err)
}

func TestUserClassic_Login(t *testing.T) {
	var response *pr.LoginResponse
	var err error

	_, err = postgresPool.Exec(context.Background(), "delete from users where login=$1", testSignUpValid[0].Login)
	require.NoError(t, err)

	_, err = handlerTest.Signup(context.Background(), &testSignUpValid[0])
	require.NoError(t, err)

	response, err = handlerTest.Login(context.Background(), &pr.LoginRequest{
		Login:    testSignUpValid[0].Login,
		Password: testSignUpValid[0].Password,
	})
	require.NoError(t, err)

	_, err = verify(response.AccessToken)
	require.NoError(t, err)

	_, err = postgresPool.Exec(context.Background(), "delete from users where login=$1", testSignUpValid[0].Login)
	require.NoError(t, err)

	_, err = handlerTest.Login(context.Background(), &pr.LoginRequest{
		Login:    testSignUpValid[0].Login,
		Password: testSignUpValid[0].Password,
	})
	require.Error(t, err)

	_, err = postgresPool.Exec(context.Background(), "delete from users where login=$1", testSignUpValid[0].Login)
	require.NoError(t, err)

	_, err = handlerTest.Signup(context.Background(), &testSignUpValid[0])
	require.NoError(t, err)

	_, err = handlerTest.Login(context.Background(), &pr.LoginRequest{
		Login:    testSignUpValid[0].Login,
		Password: "wrong password",
	})
	require.Error(t, err)

}

func TestUserClassic_Refresh(t *testing.T) {
	var response *pr.RefreshResponse
	var signupResponse *pr.SignupResponse
	var err error

	_, err = postgresPool.Exec(context.Background(), "delete from users where login=$1", testSignUpValid[0].Login)
	require.NoError(t, err)

	signupResponse, err = handlerTest.Signup(context.Background(), &testSignUpValid[0])
	require.NoError(t, err)

	response, err = handlerTest.Refresh(context.Background(), &pr.RefreshRequest{
		Login:        testSignUpValid[0].Login,
		RefreshToken: signupResponse.RefreshToken,
	})
	require.NoError(t, err)

	_, err = verify(response.AccessToken)
	require.NoError(t, err)

	_, err = postgresPool.Exec(context.Background(), "delete from users where login=$1", testSignUpValid[0].Login)
	require.NoError(t, err)

	_, err = handlerTest.Refresh(context.Background(), &pr.RefreshRequest{
		Login:        testSignUpValid[0].Login,
		RefreshToken: "",
	})
	require.Error(t, err)

	_, err = postgresPool.Exec(context.Background(), "delete from users where login=$1", testSignUpValid[0].Login)
	require.NoError(t, err)

	_, err = handlerTest.Signup(context.Background(), &testSignUpValid[0])
	require.NoError(t, err)

	_, err = handlerTest.Refresh(context.Background(), &pr.RefreshRequest{
		Login:        testSignUpValid[0].Login,
		RefreshToken: "ad",
	})
	require.Error(t, err)
}

func TestUserClassic_Update(t *testing.T) {
	var response *pr.UpdateResponse
	var signupResponse *pr.SignupResponse
	var err error

	_, err = postgresPool.Exec(context.Background(), "delete from users where login=$1", testSignUpValid[0].Login)
	require.NoError(t, err)

	signupResponse, err = handlerTest.Signup(context.Background(), &testSignUpValid[0])
	require.NoError(t, err)

	var claims *service.CustomClaims
	claims, err = verify(signupResponse.AccessToken)
	ctx := context.WithValue(context.Background(), "user", claims)

	response, err = handlerTest.Update(ctx, &pr.UpdateRequest{
		Email: testSignUpValid[0].Email,
		Name:  testSignUpValid[0].Name,
		Age:   testSignUpValid[0].Age,
	})
	require.NoError(t, err)
	require.Equal(t, response.Login, testSignUpValid[0].Login)

	for i, user := range testSignUpInvalid { //nolint:govet //all ok
		if i == 4 {
			break
		}
		_, err = postgresPool.Exec(context.Background(), "delete from users where login=$1", testSignUpValid[0].Login)
		require.NoError(t, err)

		_, err = handlerTest.Signup(context.Background(), &testSignUpValid[0])
		require.NoError(t, err)

		claims, err = verify(signupResponse.AccessToken)
		ctx = context.WithValue(context.Background(), "user", claims)
		_, err = handlerTest.Update(ctx, &pr.UpdateRequest{
			Email: user.Email,
			Name:  user.Name,
			Age:   user.Age,
		})
		require.Error(t, err)
	}
}

func server(ctx context.Context) (pr.UserServiceClient, func()) {
	buffer := 101024 * 1024
	lis := bufconn.Listen(buffer)

	baseServer := grpc.NewServer()
	pr.RegisterUserServiceServer(baseServer, handlerTest)
	go func() {
		if err := baseServer.Serve(lis); err != nil {
			logrus.Printf("error serving server: %v", err)
		}
	}()

	conn, err := grpc.DialContext(ctx, "",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return lis.Dial()
		}), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logrus.Printf("error connecting to server: %v", err)
	}

	closer := func() {
		err := lis.Close()
		if err != nil {
			logrus.Printf("error closing listener: %v", err)
		}
		baseServer.Stop()
	}

	client := pr.NewUserServiceClient(conn)

	return client, closer
}

func fileToChunks(file *os.File) ([][]byte, int) {

	fileInfo, _ := file.Stat()

	var fileSize int64 = fileInfo.Size()

	const fileChunk = 1 * (1 << 20)

	totalPartsNum := uint64(math.Ceil(float64(fileSize) / float64(fileChunk)))

	fmt.Printf("Splitting to %d pieces.\n", totalPartsNum)

	var chunks = make([][]byte, totalPartsNum)

	for i := uint64(0); i < totalPartsNum; i++ {

		partSize := int(math.Min(fileChunk, float64(fileSize-int64(i*fileChunk))))
		chunks[i] = make([]byte, partSize)

		file.Read(chunks[i])
	}
	return chunks, int(totalPartsNum)
}

func TestUserClassic_Upload(t *testing.T) {
	client, closer := server(context.Background())
	defer closer()

	type expectation struct {
		out *pr.UploadResponse
		err error
	}

	t.Parallel()

	testImageFolder := "../../fileStore"

	imagePath := fmt.Sprintf("%s/img1.avif", testImageFolder)
	file, err := os.Open(imagePath)
	require.NoError(t, err)
	defer file.Close()

	outClient, err := client.Upload(context.Background())

	outClient.Send(&pr.UploadRequest{
		Data: &pr.UploadRequest_Info{
			Info: &pr.FileInfo{
				FileType: filepath.Ext(imagePath),
			}}})

	chunks, _ := fileToChunks(file)

	for _, c := range chunks {
		outClient.Send(&pr.UploadRequest{
			Data: &pr.UploadRequest_Chunk{
				Chunk: c,
			}})
	}

	out, err := outClient.CloseAndRecv()
	require.NoError(t, err)

	logrus.Info(out)
}
