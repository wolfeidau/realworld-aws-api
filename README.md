# realworld-aws-api

This project illustrates how to build an API in [Amazon Web Services (AWS)](https://aws.amazon.com/) using [Go](https://golang.org).

# Goals

The main goal of this project is to illustrate how to build a maintainable real world REST API hosted in lambda using Go. To enable this I have added examples of:

* Contract first [OpenAPI](https://swagger.io/specification/) using code generation
* Validation of inputs in both API Gateway and the server
* Logging with meta data including lambda request identifiers
* Provide an example client using [sigv4](https://docs.aws.amazon.com/general/latest/gr/signature-version-4.html) with Go.

# Conventions

In this example I use a few conventions when deploying the software, this is done to support multiple environments, and branch based deploys which are common when building and testing.

* `AppName` - Label given service(s) with some collective role in a system.
* `Stage` - The stage where the application is running in, e.g., dev, prod.
* `Branch` - The branch this release is deployed from, typically something other than `main` or `master` is only used when testing in parallel.

# Deployment

Create an `.envrc` using the `.envrc.example` and update it with your settings, this is used with [direnv](https://direnv.net/).

```
cp .envrc.example .envrc
```

Run make to deploy the stack.

```
make
```

# Client

To invoke a simple API which is authenticated using [sigv4](https://docs.aws.amazon.com/general/latest/gr/signature-version-4.html) I have included a generated client, to use this you will need `AWS_REGION` and `AWS_PROFILE` exported as environment variables.

```
go run cmd/customer-cli/main.go --url=https://xxxxxxxxxx.execute-api.us-west-2.amazonaws.com/Prod
```

# Libraries

* [github.com/aws/aws-lambda-go](https://github.com/aws/aws-lambda-go)
* [github.com/apex/gateway](https://github.com/apex/gateway)
* [github.com/rs/zerolog](https://github.com/rs/zerolog)
* [github.com/deepmap/oapi-codegen](https://github.com/deepmap/oapi-codegen)
* [github.com/labstack/echo/v4](https://github.com/labstack/echo/v4)

# TODO

* Add the ability to run the lambda service locally without the need for SAM, this is mainly to enable compile on change for local development


# License

This application is released under Apache 2.0 license and is copyright [Mark Wolfe](https://www.wolfe.id.au).