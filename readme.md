go_localstack
===

[![Build Status](https://travis-ci.org/mitchelldavis/go_localstack.svg?branch=master)](https://travis-ci.org/mitchelldavis/go_localstack)
[![License: Unlicense](https://img.shields.io/badge/license-Unlicense-blue.svg)](http://unlicense.org/)
[![GoDoc](https://godoc.org/github.com/mitchelldavis/go_localstack/pkg/localstack?status.svg)](https://godoc.org/github.com/mitchelldavis/go_localstack/pkg/localstack)

This project makes [localstack](https://github.com/localstack/localstack) available to golang tests.

Requirements
---

- Go v1.11.0 or higher
- Docker (Tested on version 19.03.0-rc Community Edition)

Examples
---

- [All Services](/examples/allservices/allservices_test.go)
- [S3](/examples/s3/s3_test.go)

Build
---

```sh
git clone git@github.com:mitchelldavis/go_localstack.git
cd go_localstack
make
```

Contribute
---

*The more the merrier!*

Pull requests are welcome!

TODO
---

- [ ] Finish Integration Tests
  - [ ] Is the ES service something this will support?
  - [ ] Is ses something this will support?
  - [ ] The secretsmanager service is returning weird errors.
- [x] Update Readme with updated examples
- [ ] Update Documentation with updated examples
- [x] Update Localstack.EndpointFor to only redirect calls to requested services
- [ ] Fix timing on container destruction.  In rare cases (usually when a test is canceled then another one is immediatly started again) a container may not be destroyed.
- [ ] Fix the Dynamodb issue that pops up in the integration tests.
  - There is a weird error message if you look at the logs of the docker container.

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
