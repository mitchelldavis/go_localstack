package examples

import (
    "log"
    "fmt"
    "testing"
    "os"
    "github.com/mitchelldavis/go_localstack/pkg/localstack"

    "github.com/aws/aws-sdk-go/service/apigateway"
	"github.com/aws/aws-sdk-go/service/s3"
)

var LOCALSTACK *localstack.Localstack

func TestMain(t *testing.M) {
    api, _ := localstack.NewLocalstackService("apigateway")
    s3, _ := localstack.NewLocalstackService("s3")
    sqs, _ := localstack.NewLocalstackService("sqs")

    // Combine all the services we're requesting
    LOCALSTACK_SERVICES := &localstack.LocalstackServiceCollection {
        *api,
        *s3,
        *sqs,
    }

    // Initialize the services
    var err error
    LOCALSTACK, err = localstack.NewLocalstack(LOCALSTACK_SERVICES)
    if err != nil {
        log.Fatal(fmt.Sprintf("Unable to create the localstack instance: %s", err))
    }
    if LOCALSTACK == nil {
        log.Fatal("LOCALSTACK was nil.")
    }

    // If you need to initialize s3 or sqs, do it here.

    // RUN TESTS HERE
    result := t.Run()

    // We can't use defer this because os.Exit terminates the application in place
    // and the defered function won't be called.  We need to call os.Exit because
    // we need to correctly report the test results.
    LOCALSTACK.Destroy()

    os.Exit(result)
}

func Test_APIGateway(t *testing.T) {
	svc := apigateway.New(LOCALSTACK.CreateAWSSession())
	result, err := svc.GetRestApis(&apigateway.GetRestApisInput{})
	if err != nil {
		t.Error(err)
	}

    if len(result.Items) != 0 {
        t.Error("The number of Rest Apis returned should be zero.")
    }
}
// TODO:
func Test_Kinesis(t *testing.T) { }
// TODO:
func Test_Dynamodb(t *testing.T) { }
// TODO:
func Test_DynamoDBStreams(t *testing.T) { }
// TODO:
func Test_ES(t *testing.T) { }
func Test_S3(t *testing.T) {
	svc := s3.New(LOCALSTACK.CreateAWSSession())
	result, err := svc.ListBuckets(&s3.ListBucketsInput{})
	if err != nil {
		t.Error(err)
	}

    if len(result.Buckets) != 0 {
        t.Error("The number of buckets returned should be zero.")
    }
}
// TODO:
func Test_Firehose(t *testing.T) { }
// TODO:
func Test_Lambda(t *testing.T) { }
// TODO:
func Test_Sns(t *testing.T) { }
// TODO:
func Test_Sqs(t *testing.T) { }
// TODO:
func Test_Redshift(t *testing.T) { }
// TODO:
func Test_Email(t *testing.T) { }
// TODO:
func Test_Route53(t *testing.T) { }
// TODO:
func Test_Cloudformation(t *testing.T) { }
// TODO:
func Test_Monitoring(t *testing.T) { }
// TODO:
func Test_Ssm(t *testing.T) { }
// TODO:
func Test_Secretsmanager(t *testing.T) { }
// TODO:
func Test_States(t *testing.T) { }
// TODO:
func Test_Logs(t *testing.T) { }
// TODO:
func Test_Sts(t *testing.T) { }
// TODO:
func Test_Iam(t *testing.T) { }
