package localstack

import (
    "log"
    "fmt"
    "testing"
    "os"
    "strings"
    "io/ioutil"

    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/awserr"
    "github.com/aws/aws-sdk-go/service/s3"
    "github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// LOCALSTACK: A global reference to the Localstack object
var LOCALSTACK *Localstack

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
    s3Service, _ := NewLocalstackService("s3")

    // Gather up all service definitions in a single collection.
    // (Only one in this case.)
    LOCALSTACK_SERVICES := &LocalstackServiceCollection {
        *s3Service,
    }

    // Initialize the service
    var err error
    LOCALSTACK, err = NewLocalstack(LOCALSTACK_SERVICES)
    if err != nil {
        log.Fatal(fmt.Sprintf("Unable to create the instance: %s", err))
    }
    if LOCALSTACK == nil {
        log.Fatal("LOCALSTACK was nil.")
    }

    // Make sure we Destroy Localstack.  This method handles
    // stopping and removing the docker container.
    defer LOCALSTACK.Destroy()

    // If you need to initialize s3 or sqs, do it here.
    err = InitS3()
    if err != nil {
        if aerr, ok := err.(awserr.Error); ok {
            switch aerr.Code() {
            case s3.ErrCodeBucketAlreadyExists:
                log.Fatal(s3.ErrCodeBucketAlreadyExists, aerr.Error())
            case s3.ErrCodeBucketAlreadyOwnedByYou:
                log.Fatal(s3.ErrCodeBucketAlreadyOwnedByYou, aerr.Error())
            default:
                log.Fatal(aerr.Error())
            }
        } else {
            // Print the error, cast err to awserr.Error to get the Code and
            // Message from an error.
            log.Fatal(err.Error())
        }
    }

    // RUN TESTS HERE
    return t.Run()
}

func InitS3() error {
    svc := s3.New(LOCALSTACK.CreateAWSSession())

    // Create Bucket
    input := &s3.CreateBucketInput{
        Bucket: aws.String("examplebucket"),
        CreateBucketConfiguration: &s3.CreateBucketConfiguration{
            LocationConstraint: aws.String("us-east-1"),
        },
    }

    _, err := svc.CreateBucket(input)
    if err != nil {
        return err
    }

    //Upload File
    uploader := s3manager.NewUploader(LOCALSTACK.CreateAWSSession())
    _, err = uploader.Upload(&s3manager.UploadInput{
        Bucket: aws.String("examplebucket"),
        Key: aws.String("examplefile"),
        Body: strings.NewReader("Hello World"),
    })

    if err != nil {
        return err
    }

    return nil
}

// Download and check the content of the file.
func Example_s3() {
    svc := s3.New(LOCALSTACK.CreateAWSSession())
    input := &s3.GetObjectInput{
        Bucket: aws.String("examplebucket"),
        Key:    aws.String("examplefile"),
    }

    result, err := svc.GetObject(input)
    if err != nil {
        log.Fatal(err)
    }

    text, err := ioutil.ReadAll(result.Body)
    
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(string(text))
    // Output: Hello World
}


