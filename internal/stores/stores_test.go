package stores

import (
	"fmt"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/ory/dockertest/v3"
	"github.com/rs/zerolog/log"
)

const (
	defaultRegion = "us-east-1"
)

var (
	dbSvc    dynamodbiface.DynamoDBAPI
	endpoint string
)

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatal().Msgf("Could not connect to docker: %s", err)
	}

	// pulls an image, creates a container based on it and runs it
	resource, err := pool.Run("amazon/dynamodb-local", "latest", []string{})
	if err != nil {
		log.Fatal().Msgf("Could not start resource: %s", err)
	}

	endpoint = fmt.Sprintf("http://localhost:%s", resource.GetPort("8000/tcp"))

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	if err := pool.Retry(func() error {

		dbSvc = dynamodb.New(session.Must(session.NewSession(mustConfig(endpoint))))

		_, err := dbSvc.ListTables(&dynamodb.ListTablesInput{})
		if err != nil {
			log.Info().Msgf("Failed to create dynamodb client: %v", err)
			return err
		}
		log.Printf("%#v\n", endpoint)
		return nil
	}); err != nil {
		log.Fatal().Msgf("Could not connect to docker: %s", err)
	}

	code := m.Run()

	// You can't defer this because os.Exit doesn't care for defer
	if err := pool.Purge(resource); err != nil {
		log.Fatal().Msgf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}

func mustConfig(endpoint string) *aws.Config {

	creds := credentials.NewStaticCredentials("123", "test", "test")
	return &aws.Config{
		Region:      aws.String(defaultRegion),
		Endpoint:    aws.String(endpoint),
		Credentials: creds,
	}
}

func ensureVersionTable(dbSvc dynamodbiface.DynamoDBAPI, tableName string) error {

	_, err := dbSvc.CreateTable(&dynamodb.CreateTableInput{
		TableName: aws.String(tableName),
		KeySchema: []*dynamodb.KeySchemaElement{
			{AttributeName: aws.String("id"), KeyType: aws.String(dynamodb.KeyTypeHash)},
			{AttributeName: aws.String("name"), KeyType: aws.String(dynamodb.KeyTypeRange)},
		},
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{AttributeName: aws.String("id"), AttributeType: aws.String(dynamodb.ScalarAttributeTypeS)},
			{AttributeName: aws.String("name"), AttributeType: aws.String(dynamodb.ScalarAttributeTypeS)},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(1),
			WriteCapacityUnits: aws.Int64(1),
		},
		SSESpecification: &dynamodb.SSESpecification{
			Enabled: aws.Bool(true),
			SSEType: aws.String(dynamodb.SSETypeAes256),
		},
	})

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeResourceInUseException:
				return nil
			}
		}
		return err
	}

	err = dbSvc.WaitUntilTableExists(&dynamodb.DescribeTableInput{
		TableName: aws.String(tableName),
	})
	if err != nil {
		return err
	}

	_, err = dbSvc.UpdateTimeToLive(&dynamodb.UpdateTimeToLiveInput{
		TableName: aws.String(tableName),
		TimeToLiveSpecification: &dynamodb.TimeToLiveSpecification{
			AttributeName: aws.String("expires"),
			Enabled:       aws.Bool(true),
		},
	})
	if err != nil {
		return err
	}

	return nil
}
