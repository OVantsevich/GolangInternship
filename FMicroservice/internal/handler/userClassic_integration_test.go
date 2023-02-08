package handler

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/sirupsen/logrus"
)

var (
	TestPgUser          = "postgres"
	TestPgPassword      = "postgres"
	TestPgDB            = "postgres"
	TestPgPort          = "11111"
	TestPgContainerName = "postgres"
)
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

	postgresPool, err = pgxpool.New(context.Background(), fmt.Sprintf("-url=jdbc:postgresql://%s:%s/%s", TestPgContainerName, "5432", TestPgDB))
	if err != nil {
		logrus.Fatalf("Could not connect to db %s", err)
	}

	code := m.Run()

	if err := pool.Purge(resource); err != nil {
		logrus.Fatalf("Could not purge postgres: %s", err)
	}

	if err := pool.Purge(flyway); err != nil {
		logrus.Fatalf("Could not purge flyway: %s", err)
	}

	if err := pool.Client.RemoveNetwork(network.ID); err != nil {
		logrus.Fatalf("Could not remove network: %s", err)
	}

	os.Exit(code)
}

func TestUserClassic_Signup(t *testing.T) {

}
