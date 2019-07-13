package examples

import (
    "log"
    "fmt"
    "testing"
    "os"
    "github.com/mitchelldavis/go_localstack/pkg/localstack"

    "github.com/aws/aws-sdk-go/service/dynamodb"
    "github.com/aws/aws-sdk-go/service/dynamodbstreams"
)

// LOCALSTACK: A global reference to the Localstack object
var LOCALSTACK *localstack.Localstack

// In order to setup a single Localstack instance for all tests in a
// test suite, the TestMain function allows a single place to wrap all
// tests in setup and teardown logic.  
// https://golang.org/pkg/testing/#hdr-Main
func TestMain(t *testing.M) {
    os.Exit(InitializeLocalstack(t))
}

// We create a seperate iniitalize function so we can call
// `defer LOCALSTACK.Destroy()`
func InitializeLocalstack(t *testing.M) int {
    // Create the S3 Service definition
    dynamodb, _ := localstack.NewLocalstackService("dynamodb")
    dynamodbstreams, _ := localstack.NewLocalstackService("dynamodbstreams")

    // Gather up all service definitions in a single collection.
    // (Only one in this case.)
    LOCALSTACK_SERVICES := &localstack.LocalstackServiceCollection {
        *dynamodb,
        *dynamodbstreams,
    }

    // Initialize the service
    var err error
    LOCALSTACK, err = localstack.NewSpecificLocalstack(
        LOCALSTACK_SERVICES, 
        "dynamotest", 
        localstack.Localstack_Repository, 
        localstack.Localstack_Tag)
    if err != nil {
        log.Fatal(fmt.Sprintf("Unable to create the localstack instance: %s", err))
    }
    if LOCALSTACK == nil {
        log.Fatal("LOCALSTACK was nil.")
    }

    // Make sure we Destroy Localstack.  This method handles
    // stopping and removing the docker container.
    defer LOCALSTACK.Destroy()

    // RUN TESTS HERE
    return t.Run()
}
func Test_Dynamodb(t *testing.T) {
    svc := dynamodb.New(LOCALSTACK.CreateAWSSession())
    result, err := svc.ListTables(&dynamodb.ListTablesInput{ })
    if err != nil {
        t.Error(err)
    }

    if len(result.TableNames) != 0 {
        t.Error("The number of returned table names should be zero.")
    }
}
func Test_DynamoDBStreams(t *testing.T) {
    svc := dynamodbstreams.New(LOCALSTACK.CreateAWSSession())
    result, err := svc.ListStreams(&dynamodbstreams.ListStreamsInput{ })
    if err != nil {
        t.Error(err)
    }

    if len(result.Streams) != 0 {
        t.Error("The number of returned streams should be zero.")
    }
}
