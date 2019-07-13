package localstack

import (
    "errors"
    "fmt"
    "time"
	"github.com/ory/dockertest"
	"github.com/ory/dockertest/docker"
)

// DockerWrapper is used to abstract docker to make testing easier.
// Each method of this interface simply wraps functionality that already
// exists in the Client object of the github.com/ory/dockertest/docker library.
type DockerWrapper interface {
    // See https://godoc.org/github.com/ory/dockertest/docker#Client.InspectContainer
	InspectContainer(string) (*docker.Container, error)
    // See https://godoc.org/github.com/ory/dockertest/docker#Client.ListContainers
	ListContainers(docker.ListContainersOptions) ([]docker.APIContainers, error)
    // See https://godoc.org/github.com/ory/dockertest/docker#Client.RunWithOptions
	RunWithOptions(*dockertest.RunOptions, ...func(*docker.HostConfig)) (*dockertest.Resource, error)
    // See https://godoc.org/github.com/ory/dockertest/docker#Client.Retry
	Retry(func() error) error
}

type _DockerWrapper struct { }

func (dw *_DockerWrapper) InspectContainer(id string) (*docker.Container, error) {
	client, err := docker.NewClientFromEnv()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Unable to create a docker client: %s", err))
	}
	return client.InspectContainer(id)
}

func (dw *_DockerWrapper) ListContainers(options docker.ListContainersOptions) ([]docker.APIContainers, error) {
	client, err := docker.NewClientFromEnv()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Unable to create a docker client: %s", err))
	}
	
	return client.ListContainers(options)
}

func (dw *_DockerWrapper) RunWithOptions(opts *dockertest.RunOptions, hcOpts ...func(*docker.HostConfig)) (*dockertest.Resource, error) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Could not connect to docker: %s", err))
	}
	return pool.RunWithOptions(opts, hcOpts...)
}

func (dw *_DockerWrapper) Retry(op func() error) error {
	pool, err := dockertest.NewPool("")
    pool.MaxWait = time.Minute * 5
	if err != nil {
		return errors.New(fmt.Sprintf("Could not connect to docker: %s", err))
	}
	return pool.Retry(op)
}
