go_localstack
===

[![Build Status](https://travis-ci.org/mitchelldavis/go_localstack.svg?branch=master)](https://travis-ci.org/mitchelldavis/go_localstack)
[![License: Unlicense](https://img.shields.io/badge/license-Unlicense-blue.svg)](http://unlicense.org/)
[![GoDoc](https://godoc.org/github.com/mitchelldavis/go_localstack/pkg/localstack?status.svg)](https://godoc.org/github.com/mitchelldavis/go_localstack/pkg/localstack)

This project makes [localstack](https://github.com/localstack/localstack) available to golang tests.

Requirements
---

- Go v1.12.0 or higher
- Docker (Tested on version 19.03.0-rc Community Edition)

Example
---

```go

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
```

Build
---

```sh
git clone git@github.com:mitchelldavis/go_localstack.git
cd go_localstack
make
```

Contribute
---

The more the merrier!

License
---

This is free and unencumbered software released into the public domain.

Anyone is free to copy, modify, publish, use, compile, sell, or
distribute this software, either in source code form or as a compiled
binary, for any purpose, commercial or non-commercial, and by any
means.

In jurisdictions that recognize copyright laws, the author or authors
of this software dedicate any and all copyright interest in the
software to the public domain. We make this dedication for the benefit
of the public at large and to the detriment of our heirs and
successors. We intend this dedication to be an overt act of
relinquishment in perpetuity of all present and future rights to this
software under copyright law.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
IN NO EVENT SHALL THE AUTHORS BE LIABLE FOR ANY CLAIM, DAMAGES OR
OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE,
ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
OTHER DEALINGS IN THE SOFTWARE.

For more information, please refer to [http://unlicense.org](http://unlicense.org)
