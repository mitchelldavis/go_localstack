package localstack

/*
go_localstack

This package was written to help writing tests with Localstack.  
(https://github.com/localstack/localstack)  It uses libraries that help create
and manage a Localstack docker container for your go tests.

Requirements

- Go v1.12.0 or higher
- Docker (Tested on version 19.03.0-rc Community Edition)

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
