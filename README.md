# realworld-aws-api

This project illustrates how to build an API in [Amazon Web Services (AWS)](https://aws.amazon.com/) using [Go](https://golang.org).

# Goals

The main goal of this project is to illustrate how to build a maintainable real world REST API hosted in lambda using Go. To enable this I have added examples of:

* Contract first [OpenAPI](https://swagger.io/specification/) using code generation
* Validation of inputs in both API Gateway and the server
* Logging with meta data including lambda request identifiers
* Provide an [API Gateway](https://aws.amazon.com/api-gateway/) client using [sigv4](https://docs.aws.amazon.com/general/latest/gr/signature-version-4.html) with Go.
* Protobuf based versioned storage using [DynamoDB](https://aws.amazon.com/dynamodb/)
* Local testing with [docker](https://www.docker.com) using [github.com/ory/dockertest/v3](https://github.com/ory/dockertest) and [DynamoDB Local](https://hub.docker.com/r/amazon/dynamodb-local/)
* CLI using [github.com/alecthomas/kong](https://github.com/alecthomas/kong) with logging and sub commands
* Linting using [golangci-lint](https://github.com/golangci/golangci-lint)

# CLI

The CLI provides a simple interface to access data via the API.

```
$ customer-cli --help
Usage: customer-cli --url=STRING <command>

Flags:
  -h, --help          Show context-sensitive help.
      --version
      --debug
      --url=STRING

Commands:
  create-customer --url=STRING --name=STRING --labels=LABELS,...
    New Customer.

  get-customer --url=STRING --id=STRING
    Read Customer.

  list-customers --url=STRING
    Read a list of Customers.

Run "customer-cli <command> --help" for more information on a command.

```

Reading a list of customers from the API.

```console
$ customer-cli --url=https://xxxxxxxxxx.execute-api.us-west-2.amazonaws.com/Prod list-customers  | jq .
6:33PM INF cmd/customer-cli/commands/list_customers.go:18 > get a list of customers from the api
6:33PM INF cmd/customer-cli/apigw/apigw.go:25 > signing request host=z0d3zmnwh1.execute-api.us-west-2.amazonaws.com
{
  "customers": [
    {
      "created_at": "2020-10-22T17:38:34.242777542Z",
      "description": "test",
      "id": "01EN8P84M2P9RQJ1XV3XQR4DZM",
      "labels": [
        "test"
      ],
      "name": "test",
      "updated_at": "2020-10-22T17:38:34.242778239Z"
    },
    {
      "created_at": "2020-10-22T17:41:23.701550324Z",
      "description": "test",
      "id": "01EN8PDA3N91HNQEVV16HX6J27",
      "labels": [
        "test"
      ],
      "name": "test2",
      "updated_at": "2020-10-22T17:41:23.701551096Z"
    },
    {
      "created_at": "2020-10-23T02:21:37.259975291Z",
      "description": "test",
      "id": "01EN9M5W3BDKGR3RGCEGNSBYHQ",
      "labels": [
        "test"
      ],
      "name": "test3",
      "updated_at": "2020-10-23T02:21:37.259975968Z"
    }
  ]
}
```

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
* [github.com/wolfeidau/dynastore](https://github.com/wolfeidau/dynastore)

# TODO

* Add the ability to run the lambda service locally without the need for SAM, this is mainly to enable compile on change for local development


# License

This application is released under Apache 2.0 license and is copyright [Mark Wolfe](https://www.wolfe.id.au).