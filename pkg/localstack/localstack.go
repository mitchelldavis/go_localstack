/*

go_localstack

This package was written to help writing tests with Localstack.  
(https://github.com/localstack/localstack)  It uses libraries that help create
and manage a Localstack docker container for your go tests.

Requirements

    Go v1.12.0 or higher
    Docker (Tested on version 19.03.0-rc Community Edition)

Example

Within a test file:

    func TestMain(t *testing.M) {
        // Here we define a S3 Localstack Service Definition
        s3, err := localstack.NewLocalstackService("s3")
        if err != nil {
            log.Fatal(fmt.Sprintf("Unable to create the s3 Service: %s", err))
        }

        // Here we define a SQS Localstack Service Definition
        sqs, err := localstack.NewLocalstackService("sqs")
        if err != nil {
            log.Fatal(fmt.Sprintf("Unable to create the sqs Service: %s", err))
        }

        // Combine all the services we're requesting
        LOCALSTACK_SERVICES := &localstack.LocalstackServiceCollection {
            *s3,
            *sqs,
        }

        // Initialize the services
        LOCALSTACK, err := localstack.NewLocalstack(LOCALSTACK_SERVICES)
        if err != nil {
            log.Fatal(fmt.Sprintf("Unable to create the localstack instance: %s", err))
        }
        if LOCALSTACK == nil {
            log.Fatal("LOCALSTACK was nil.")
        }

        // If you need to initialize s3 or sqs, do it here.

        // RUN TESTS HERE
        result := t.Run()

        // We can't defer this because os.Exit terminates the application in place
        // and the defered function won't be called.  We need to call os.Exit because
        // we need to correctly report the test results.
        LOCALSTACK.Destroy()

        os.Exit(result)
    }
*/
package localstack

import (
	"errors"
	"fmt" 
	"strings"
	"bytes"
	"bufio"
	"github.com/ory/dockertest"
	"github.com/ory/dockertest/docker"
)

// Localstack_Repository is the Localstack Docker repository
const Localstack_Repository string = "localstack/localstack"
// Localstack_Tag is the last tested version of the Localstack Docker repository
const Localstack_Tag string = "0.9.1"

// Localstack is a structure used to control the lifecycle of the Localstack 
// Docker container.
type Localstack struct {
    // Resource is a pointer to the dockertest.Resource 
    // object that is the localstack docker container.
    // (https://godoc.org/github.com/ory/dockertest#Resource)
	Resource *dockertest.Resource
    // Services is a pointer to a collection of service definitions
    // that are being requested from this particular instance of Localstack.
	Services *LocalstackServiceCollection
}

// Destroy simply shuts down and cleans up the Localstack container out of docker.
func (ls *Localstack) Destroy() error {
	
	pool, err := dockertest.NewPool("")
	if err != nil {
		return errors.New(fmt.Sprintf("Could not connect to docker: %s", err))
	}

	// You can't defer this because os.Exit doesn't care for defer
	if err := pool.Purge(ls.Resource); err != nil {
		return errors.New(fmt.Sprintf("Could not purge resource: %s", err))
	}

	return nil
}

// NewLocalstack creates a new Localstack docker container based on the latest version.
func NewLocalstack(services *LocalstackServiceCollection) (*Localstack, error) {
	return NewSpecificLocalstack(services, Localstack_Repository, "latest")
}

// NewSpecificLocalstack creates a new Localstack docker container based on
// the given repository and tag given.  NOTE:  The Docker image used should be a 
// Localstack image.  The behavior is unknown otherwise.  This method is provided
// to allow special situations like using a tag other than latest or when referencing 
// an internal Localstack image.
func NewSpecificLocalstack(services *LocalstackServiceCollection, repository, tag string) (*Localstack, error) {
	return newLocalstack(services, &_DockerWrapper{ }, repository, tag)
}

func getLocalstack(services *LocalstackServiceCollection, dockerWrapper DockerWrapper, repository, tag string) (*dockertest.Resource, error) {
	
	containers, err := dockerWrapper.ListContainers(docker.ListContainersOptions { All: true })
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Unable to retrieve docker containers: %s", err))
	}
	for _, c := range containers {
		if c.Image == fmt.Sprintf("%s:%s", repository, tag) {
			container, err := dockerWrapper.InspectContainer(c.ID)
			if err != nil { return nil, errors.New(fmt.Sprintf("Unable to inspect container %s: %s", c.ID, err))
			}
			for _, env := range container.Config.Env {
				if env == fmt.Sprintf("SERVICES=%s", services.GetServiceMap()) {
					return &dockertest.Resource{ Container: container }, nil
				}
			}

			return nil, errors.New("We're only supporting one Localstack instance at a time.")
		}
	}

	return nil, nil
}

func newLocalstack(services *LocalstackServiceCollection, wrapper DockerWrapper, repository, tag string) (*Localstack, error) {

	localstack, err := getLocalstack(services, wrapper, repository, tag)
	if err != nil {
		return nil, err	
	}

	if localstack == nil {

		// Fifth, If we didn't find a running container before, we spin one up now.
		localstack, err = wrapper.RunWithOptions(&dockertest.RunOptions{
			Repository: repository,
			Tag: tag,
			Env: []string{
				fmt.Sprintf("SERVICES=%s", services.GetServiceMap()),
			},
		})
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Could not start resource: %s", err))
		}
	}

	// Sixth, we wait for the services to be ready before we allow the tests
	// to be run.
	for _, service := range *services {
		//fmt.Println(fmt.Sprintf("Testing connectivity with %s", service.Name))
		if err := wrapper.Retry(func() error {

			// We have to use a method that checks the output
			// of the docker container here because simply checking for
			// connetivity on the ports doesn't work.
			client, err := docker.NewClientFromEnv()
			if err != nil {
				fmt.Println(err)
				return errors.New(fmt.Sprintf("Unable to create a docker client: %s", err))
			}

			buffer := new(bytes.Buffer)

			logsOptions := docker.LogsOptions {
				Container: localstack.Container.ID,
				OutputStream: buffer,
				RawTerminal: true,
				Stdout: true,
				Stderr: true,
			}
			err = client.Logs(logsOptions)
			if err != nil {
				fmt.Println(err)
				return errors.New(fmt.Sprintf("Unable to retrieve logs for container %s: %s", localstack.Container.ID, err))
			}

			scanner := bufio.NewScanner(buffer)
			for scanner.Scan() {
				token := strings.TrimSpace(scanner.Text())
				expected := "Ready."
				if strings.Contains(strings.TrimSpace(token),expected) {
					fmt.Println(token)
					return nil
				}
			}
			if err := scanner.Err(); err != nil {
				fmt.Println(err)
				return errors.New(fmt.Sprintf("Reading input: %s", err))
			}
			return errors.New("Not Ready")
		}); err != nil {
			return nil, errors.New(fmt.Sprintf("Unable to connect to %s: %s", service.Name, err))
		}
	}

	return &Localstack{
		Resource: localstack,
		Services: services,
	}, nil
}

