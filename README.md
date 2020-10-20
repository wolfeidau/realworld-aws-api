# realworld-aws-api

This project illustrates how to build an API in [Amazon Web Services (AWS)](https://aws.amazon.com/) using [Go](https://golang.org).

# Goals

The main goal of this project is to illustrate how to build a maintainable real world REST API hosted in lambda using Go. To enable this I have added examples of:

* Contract first [OpenAPI](https://swagger.io/specification/) using code generation
* Validation of inputs in both API Gateway and the server
* The ability to run the lambda service locally without the need for SAM, this is mainly to enable compile on change for local development
* Logging with meta data including lambda request identifiers 

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

# Libraries

* [github.com/aws/aws-lambda-go](https://github.com/aws/aws-lambda-go)
* [github.com/apex/gateway](https://github.com/apex/gateway)
* [github.com/rs/zerolog](https://github.com/rs/zerolog)
* [github.com/deepmap/oapi-codegen](https://github.com/deepmap/oapi-codegen)
* [github.com/labstack/echo/v4](https://github.com/labstack/echo/v4)

# License

This application is released under Apache 2.0 license and is copyright [Mark Wolfe](https://www.wolfe.id.au).