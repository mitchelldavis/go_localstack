package localstack

import (
    "log"
    "fmt"
    "strings"
    "io/ioutil"

    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/service/s3"
    "github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// In this example we initialize Localstack with the S3 service
// enabled, create a bucket, upload a file to that bucket, then 
// download that file from s3 and output it's content.
//
// For complete testing examples, see the examples in the github.com
// repository. 
// https://github.com/mitchelldavis/go_localstack/tree/master/examples 
func Example_s3() {
    // LOCALSTACK: A reference to the Localstack object
    var LOCALSTACK *Localstack
    
    // Create the S3 Service definition
    s3Service, _ := NewLocalstackService("s3")

    // Gather up all service definitions in a single collection.
    // (Only one in this case.)
    LOCALSTACK_SERVICES := &LocalstackServiceCollection {
        *s3Service,
    }

    // Initialize Localstack.  Here Localstack is created and
    // is ready to go.
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

    // Here we start the code to interact with S3
    svc := s3.New(LOCALSTACK.CreateAWSSession())

    // Create Bucket
    input := &s3.CreateBucketInput{
        Bucket: aws.String("examplebucket"),
        CreateBucketConfiguration: &s3.CreateBucketConfiguration{
            LocationConstraint: aws.String("us-east-1"),
        },
    }

    _, err = svc.CreateBucket(input)
    if err != nil {
        log.Fatal(err)
    }

    //Upload File
    uploader := s3manager.NewUploader(LOCALSTACK.CreateAWSSession())
    _, err = uploader.Upload(&s3manager.UploadInput{
        Bucket: aws.String("examplebucket"),
        Key: aws.String("examplefile"),
        Body: strings.NewReader("Hello World"),
    })
    if err != nil {
        log.Fatal(err)
    }

    // Download the file
    getObjectInput := &s3.GetObjectInput{
        Bucket: aws.String("examplebucket"),
        Key:    aws.String("examplefile"),
    }

    result, err := svc.GetObject(getObjectInput)
    if err != nil {
        log.Fatal(err)
    }

    // Read the contents of the file.
    text, err := ioutil.ReadAll(result.Body)
    
    if err != nil {
        log.Fatal(err)
    }

    // Print the contents of the file out.
    fmt.Println(string(text))
    // Output: Hello World
}
