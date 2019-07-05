package localstack

import (
	"errors"
	"fmt" 
	"strings"
    "sort"
)

// LocalstackService defines a particular AWS service requested for a Localstack
// instances.  Note: You shold not create an instance of LocastackService directly.
// See: NewLocalstackService
type LocalstackService struct {
    // The name of the AWS Service. (I.E. "s3" or "apigateway")
	Name string
    // Protocol is the network protocol used for communication.
	Protocol string
    // Port is the port used when communicating with the service in the
    // Localstack instance.
	Port int
}

// Equals returns wether two pointers to a LocalstackService are equal.
func (service *LocalstackService) Equals(rhs *LocalstackService) bool {
	if service == nil && rhs == nil {
		return true
	} else if service == nil || rhs == nil {
		return false
	} else {
		return	service.Name == rhs.Name &&
				service.Protocol == rhs.Protocol &&
				service.Port == rhs.Port
	}
}

// GetPortProtocol returns the protocol string (eg. 1234/tcp) used by Docker.
func (service *LocalstackService) GetPortProtocol() string {
	return fmt.Sprintf("%d/%s", service.Port, service.Protocol)	
}

// GetNameProtocol returns the protocol string (eg. s3:1234) used by Docker.
func (service *LocalstackService) GetNamePort() string {
	return fmt.Sprintf("%s:%d", service.Name, service.Port) 
}

// NewLocalstackService returns a new pointer to an instance of LocalstackService
// given the name of the service provided.  Note: The name must match an aws service
// from this list (https://docs.aws.amazon.com/cli/latest/reference/#available-services)
// and be a supported service by Localstack.
func NewLocalstackService(name string) (*LocalstackService, error) {

	switch name {
	case "apigateway":
		return &LocalstackService {
			Name: "apigateway",
			Protocol: "tcp",
			Port: 4567,
		}, nil
	case "kinesis":
		return &LocalstackService {
			Name: "kinesis",
			Protocol: "tcp",
			Port: 4568,
		}, nil
	case "dynamodb":
		return &LocalstackService {
			Name: "dynamodb",
			Protocol: "tcp",
			Port: 4569,
		}, nil
	case "dynamodbstreams":
		return &LocalstackService {
			Name: "dynamodb",
			Protocol: "tcp",
			Port: 4570,
		}, nil
	case "es":
		return &LocalstackService {
			Name: "es",
			Protocol: "tcp",
			Port: 4571,
		}, nil
	case "s3":
		return &LocalstackService {
			Name: "s3",
			Protocol: "tcp",
			Port: 4572,
		}, nil
	case "firehose":
		return &LocalstackService {
			Name: "firehose",
			Protocol: "tcp",
			Port: 4573,
		}, nil
	case "lambda":
		return &LocalstackService {
			Name: "lambda",
			Protocol: "tcp",
			Port: 4574,
		}, nil
	case "sns":
		return &LocalstackService {
			Name: "sns",
			Protocol: "tcp",
			Port: 4575,
		}, nil
	case "sqs":
		return &LocalstackService {
			Name: "sqs",
			Protocol: "tcp",
			Port: 4576,
		}, nil
	case "redshift":
		return &LocalstackService {
			Name: "redshift",
			Protocol: "tcp",
			Port: 4577,
		}, nil
	case "ses":
		return &LocalstackService {
			Name: "ses",
			Protocol: "tcp",
			Port: 4579,
		}, nil
	case "route53":
		return &LocalstackService {
			Name: "route53",
			Protocol: "tcp",
			Port: 4580,
		}, nil
	case "cloudformation":
		return &LocalstackService {
			Name: "cloudformation",
			Protocol: "tcp",
			Port: 4581,
		}, nil
	case "cloudwatch":
		return &LocalstackService {
			Name: "cloudwatch",
			Protocol: "tcp",
			Port: 4582,
		}, nil
	case "ssm":
		return &LocalstackService {
			Name: "ssm",
			Protocol: "tcp",
			Port: 4583,
		}, nil
	case "secretsmanager":
		return &LocalstackService {
			Name: "secretsmanager",
			Protocol: "tcp",
			Port: 4584,
		}, nil
	case "stepfunctions":
		return &LocalstackService {
			Name: "stepfunctions",
			Protocol: "tcp",
			Port: 4585,
		}, nil
	case "logs":
		return &LocalstackService {
			Name: "logs",
			Protocol: "tcp",
			Port: 4586,
		}, nil
	case "sts":
		return &LocalstackService {
			Name: "sts",
			Protocol: "tcp",
			Port: 4592,
		}, nil
	case "iam":
		return &LocalstackService {
			Name: "iam",
			Protocol: "tcp",
			Port: 4593,
		}, nil
	default:
		return nil, errors.New(fmt.Sprintf("Unknown Localstack Service: %s", name))
	}
}

// LocalstackServiceCollection represents a collection of LocalstackService objects.
type LocalstackServiceCollection []LocalstackService

// GetServiceMap returns a comma delimited string of all the AWS service
// names in the collection.
func (collection *LocalstackServiceCollection) GetServiceMap() string {
	var maps []string
	for _, element := range *collection {
		maps = append(maps, element.GetNamePort())
	}

	return strings.Join(maps, ",")
}
// Len returns the number of items in the collection.
func (a LocalstackServiceCollection) Len() int { 
	return len(a) 
}
// Swap will swap two items in the collection.
func (a LocalstackServiceCollection) Swap(i, j int) { 
	a[i], a[j] = a[j], a[i]
}
// Less compares two items in the collection.  This returns true if the instance
// at i is less than the instance at j.  Otherwise it will return false.
func (a LocalstackServiceCollection) Less(i, j int) bool { 
	return a[i].Name < a[j].Name
}
// Sort simply sorts the collection based on the names of the defined services.
// The collection returned is a pointer to the calling collection.
func (a *LocalstackServiceCollection) Sort() *LocalstackServiceCollection { 
	sort.Sort(a)
	return a
}
