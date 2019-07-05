package localstack

import (
	"fmt"
	"errors"
	"log"
	"testing"
	"github.com/ory/dockertest"
	"github.com/ory/dockertest/docker"
	"github.com/mitchelldavis/go_localstack/pkg/mock_localstack"
	"github.com/golang/mock/gomock"
)

func getLocalstack_Found(services *LocalstackServiceCollection, ctrl *gomock.Controller) (*mock_localstack.MockDockerWrapper, *docker.Container) {
	m := mock_localstack.NewMockDockerWrapper(ctrl)
	container := &docker.Container {
		Config: &docker.Config {
			Env: []string {
				fmt.Sprintf("SERVICES=%s", services.GetServiceMap()),
			},
		},
	}

	m.
	EXPECT().
	ListContainers(gomock.Any()).
	Times(1).
	Return([]docker.APIContainers {
		docker.APIContainers {Image: fmt.Sprintf("%s:%s", Localstack_Repository, Localstack_Tag)},
	}, nil)

	m.
	EXPECT().
	InspectContainer(gomock.Any()).
	Times(1).
	Return(container, nil)

	return m, container
}

func getLocalstack_Empty(services *LocalstackServiceCollection, ctrl *gomock.Controller) *mock_localstack.MockDockerWrapper {
	m := mock_localstack.NewMockDockerWrapper(ctrl)

	m.
	EXPECT().
	ListContainers(gomock.Any()).
	Times(1).
	Return(nil, nil)

	return m
}

func Test_getLocalstack_ErrorWithListContainers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()	

	m := mock_localstack.NewMockDockerWrapper(ctrl)

	m.
	EXPECT().
	ListContainers(gomock.Any()).
	Times(1).
	Return(nil, errors.New("Dummy Error"))

	m.
	EXPECT().
	InspectContainer(gomock.Any()).
	Times(0).
	Return(nil, nil)

	sqs, _ := NewLocalstackService("sqs")
	services := &LocalstackServiceCollection {
		*sqs,
	}
	actual, err := getLocalstack(services, m, Localstack_Repository, Localstack_Tag)

	if actual != nil {
		log.Fatal("We're expecting the localstack result to be nil.")
	}

	if err == nil {
		log.Fatal("We're expecting the error returned to be populated.")	
	}
}

func Test_getLocalstack_UnknownImage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()	

	m := mock_localstack.NewMockDockerWrapper(ctrl)

	m.
	EXPECT().
	ListContainers(gomock.Any()).
	Times(1).
	Return([]docker.APIContainers {
		docker.APIContainers {Image: "DummyImage:1.0.0"},
	}, nil)

	m.
	EXPECT().
	InspectContainer(gomock.Any()).
	Times(0).
	Return(nil, nil)

	sqs, _ := NewLocalstackService("sqs")
	services := &LocalstackServiceCollection {
		*sqs,
	}
	actual, err := getLocalstack(services, m, Localstack_Repository, Localstack_Tag)

	if actual != nil || err != nil {
		log.Fatal("We're expecting both the localstack and error return results to be nil.")
	}
}

func Test_getLocalstack_ErrorWithInspectContainer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()	

	m := mock_localstack.NewMockDockerWrapper(ctrl)

	m.
	EXPECT().
	ListContainers(gomock.Any()).
	Times(1).
	Return([]docker.APIContainers {
		docker.APIContainers {Image: fmt.Sprintf("%s:%s", Localstack_Repository, Localstack_Tag)},
	}, nil)

	m.
	EXPECT().
	InspectContainer(gomock.Any()).
	Times(1).
	Return(nil, errors.New("Dummy Error"))

	sqs, _ := NewLocalstackService("sqs")
	services := &LocalstackServiceCollection {
		*sqs,
	}
	actual, err := getLocalstack(services, m, Localstack_Repository, Localstack_Tag)

	if actual != nil {
		log.Fatal("We're expecting the localstack result to be nil.")
	}

	if err == nil {
		log.Fatal("We're expecting the error returned to be populated.")	
	}
}

func Test_getLocalstack_ContainerExistsButHasDifferentServices(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()	

	m := mock_localstack.NewMockDockerWrapper(ctrl)

	// Setup call to ListContainers
	m.
	EXPECT().
	ListContainers(gomock.Any()).
	Times(1).
	Return([]docker.APIContainers {
		docker.APIContainers {Image: fmt.Sprintf("%s:%s", Localstack_Repository, Localstack_Tag)},
	}, nil)

	m.
	EXPECT().
	InspectContainer(gomock.Any()).
	Times(1).
	Return(&docker.Container {
		Config: &docker.Config {
			Env: []string {
				"NOTSERVICES=DUMMY",
			},
		},
	}, nil)

	sqs, _ := NewLocalstackService("sqs")
	services := &LocalstackServiceCollection {
		*sqs,
	}
	actual, err := getLocalstack(services, m, Localstack_Repository, Localstack_Tag)

	if actual != nil {
		log.Fatal("We're expecting the localstack result to be nil.")
	}

	if err == nil {
		log.Fatal("We're expecting the error returned to be populated.")	
	}
}

func Test_getLocalstack_UnableToFindContainer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()	

	m := mock_localstack.NewMockDockerWrapper(ctrl)

	m.
	EXPECT().
	ListContainers(gomock.Any()).
	Times(1).
	Return(nil, nil)

	m.
	EXPECT().
	InspectContainer(gomock.Any()).
	Times(0).
	Return(nil, nil)

	sqs, _ := NewLocalstackService("sqs")
	services := &LocalstackServiceCollection {
		*sqs,
	}
	actual, err := getLocalstack(services, m, Localstack_Repository, Localstack_Tag)

	if actual != nil || err != nil {
		log.Fatal("We're expecting both the localstack and error return results to be nil.")
	}
}

func Test_getLocalstack_ContainerExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()	

	sqs, _ := NewLocalstackService("sqs")
	services := &LocalstackServiceCollection {
		*sqs,
	}
	m, c := getLocalstack_Found(services, ctrl)
	
	actual, err := getLocalstack(services, m, Localstack_Repository, Localstack_Tag)

	if err != nil {
		log.Fatal("We're expecting the error returned to be nil.")	
	}

	if actual.Container != c {
		log.Fatal("The actual result doesn't match what was expected.")
	}
}

func Test_NewLocalstack_GetLocalstackReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()	

	sqs, _ := NewLocalstackService("sqs")
	services := &LocalstackServiceCollection {
		*sqs,
	}
	m := mock_localstack.NewMockDockerWrapper(ctrl)

	m.
	EXPECT().
	ListContainers(gomock.Any()).
	Times(1).
	Return(nil, errors.New("Dummy Error"))

	result, err := newLocalstack(services, m, Localstack_Repository, Localstack_Tag)

	if result != nil {
		log.Fatal("We were expecting the returned container to be nil.")
	}
	
	if err == nil {
		log.Fatal("We were expecting the returned error to be populated.")
	}
}

func Test_NewLocalstack_GetLocalstackReturnsResult(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()	

	sqs, _ := NewLocalstackService("sqs")
	s3, _ := NewLocalstackService("s3")
	services := &LocalstackServiceCollection {
		*sqs,
		*s3,
	}
	m, c := getLocalstack_Found(services, ctrl)

	m.
	EXPECT().
	RunWithOptions(gomock.Any()).
	Times(0)

	m.
	EXPECT().
	Retry(gomock.Any()).
	Times(2).
	Return(nil)

	result, err := newLocalstack(services, m, Localstack_Repository, Localstack_Tag)

	if err != nil {
		log.Fatal("We were expecting the returned error to be nil.")
	}

	if result.Resource.Container != c {
		log.Fatal("The actual result doesn't match what was expected.")
	}
}

func Test_NewLocalstack_GetLocalstackReturnsResult_RetryFailsOnFirstService(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()	

	sqs, _ := NewLocalstackService("sqs")
	s3, _ := NewLocalstackService("s3")
	services := &LocalstackServiceCollection {
		*sqs,
		*s3,
	}
	m, _ := getLocalstack_Found(services, ctrl)

	m.
	EXPECT().
	RunWithOptions(gomock.Any()).
	Times(0)

	gomock.InOrder(
		m.
		EXPECT().
		Retry(gomock.Any()).
		Return(errors.New("DummyError")),
		m.
		EXPECT().
		Retry(gomock.Any()).
		Times(0).
		Return(nil),
	)

	result, err := newLocalstack(services, m, Localstack_Repository, Localstack_Tag)
	
	if result != nil {
		log.Fatal("We were expecting the returned container to be nil.")
	}

	if err == nil {
		log.Fatal("We were expecting the returned error to be populated.")
	}
}

func Test_NewLocalstack_GetLocalstackReturnsResult_RetryFailsOnSecondtService(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()	

	sqs, _ := NewLocalstackService("sqs")
	s3, _ := NewLocalstackService("s3")
	services := &LocalstackServiceCollection {
		*sqs,
		*s3,
	}
	m, _ := getLocalstack_Found(services, ctrl)

	m.
	EXPECT().
	RunWithOptions(gomock.Any()).
	Times(0)

	gomock.InOrder(
		m.
		EXPECT().
		Retry(gomock.Any()).
		Return(nil),
		m.
		EXPECT().
		Retry(gomock.Any()).
		Return(errors.New("DummyError")),
	)

	result, err := newLocalstack(services, m, Localstack_Repository, Localstack_Tag)
	
	if result != nil {
		log.Fatal("We were expecting the returned container to be nil.")
	}

	if err == nil {
		log.Fatal("We were expecting the returned error to be populated.")
	}
}

func Test_NewLocalstack_StartFreshContainer_RunWithOptionsReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()	

	sqs, _ := NewLocalstackService("sqs")
	s3, _ := NewLocalstackService("s3")
	services := &LocalstackServiceCollection {
		*sqs,
		*s3,
	}
	m := getLocalstack_Empty(services, ctrl)

	m.
	EXPECT().
	RunWithOptions(gomock.Any()).
	Times(1).
	Return(nil, errors.New("Dummy Error"))

	m.
	EXPECT().
	Retry(gomock.Any()).
	Times(0).
	Return(nil)

	result, err := newLocalstack(services, m, Localstack_Repository, Localstack_Tag)
	
	if result != nil {
		log.Fatal("We were expecting the returned container to be nil.")
	}

	if err == nil {
		log.Fatal("We were expecting the returned error to be populated.")
	}
}

func Test_NewLocalstack_StartFreshContainer_RunWithOptions(t *testing.T) {
	ctrl := gomock.NewController(t)

	defer ctrl.Finish()	

	sqs, _ := NewLocalstackService("sqs")
	s3, _ := NewLocalstackService("s3")
	services := &LocalstackServiceCollection {
		*sqs,
		*s3,
	}
	m := getLocalstack_Empty(services, ctrl)
	resource := &dockertest.Resource{ }

	m.
	EXPECT().
	RunWithOptions(gomock.Any()).
	Times(1).
	Return(resource, nil)

	m.
	EXPECT().
	Retry(gomock.Any()).
	Times(2).
	Return(nil)

	result, err := newLocalstack(services, m, Localstack_Repository, Localstack_Tag)
	
	if err != nil {
		log.Fatal("We were expecting the returned error to be nil.")
	}

	if result == nil {
		log.Fatal("We were expecting the returned container to be populated.")
	}

	if result.Services != services {
		log.Fatal("The returned collection of services doesn't match what we sent in.")
	}

	if result.Resource != resource {
		log.Fatal("The returned resource is not what is expected.")
	}
}
