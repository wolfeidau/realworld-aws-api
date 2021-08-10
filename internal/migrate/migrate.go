package migrate

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/rs/zerolog/log"
)

// Table migrate a table
func Table(awscfg *aws.Config, tableName string) error {
	log.Info().Msg("run migrations")

	dbSvc := dynamodb.New(session.Must(session.NewSession(awscfg)))

	return EnsureTable(dbSvc, tableName)
}

// EnsureTable ensure the table exists
func EnsureTable(dbSvc dynamodbiface.DynamoDBAPI, tableName string) error {
	log.Info().Str("tableName", tableName).Msg("run CreateTable")

	_, err := dbSvc.CreateTable(&dynamodb.CreateTableInput{
		TableName: aws.String(tableName),
		KeySchema: []*dynamodb.KeySchemaElement{
			{AttributeName: aws.String("id"), KeyType: aws.String(dynamodb.KeyTypeHash)},
			{AttributeName: aws.String("name"), KeyType: aws.String(dynamodb.KeyTypeRange)},
		},
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{AttributeName: aws.String("id"), AttributeType: aws.String(dynamodb.ScalarAttributeTypeS)},
			{AttributeName: aws.String("name"), AttributeType: aws.String(dynamodb.ScalarAttributeTypeS)},
			{AttributeName: aws.String("created"), AttributeType: aws.String(dynamodb.ScalarAttributeTypeS)},
			{AttributeName: aws.String("pk1"), AttributeType: aws.String(dynamodb.ScalarAttributeTypeS)},
			{AttributeName: aws.String("sk1"), AttributeType: aws.String(dynamodb.ScalarAttributeTypeS)},
		},
		LocalSecondaryIndexes: []*dynamodb.LocalSecondaryIndex{
			{
				IndexName: aws.String("idx_created"),
				KeySchema: []*dynamodb.KeySchemaElement{
					{AttributeName: aws.String("id"), KeyType: aws.String(dynamodb.KeyTypeHash)},
					{AttributeName: aws.String("created"), KeyType: aws.String(dynamodb.KeyTypeRange)},
				},
				Projection: &dynamodb.Projection{ProjectionType: aws.String(dynamodb.ProjectionTypeAll)},
			},
		},
		GlobalSecondaryIndexes: []*dynamodb.GlobalSecondaryIndex{
			{
				IndexName: aws.String("idx_global_1"),
				KeySchema: []*dynamodb.KeySchemaElement{
					{AttributeName: aws.String("pk1"), KeyType: aws.String(dynamodb.KeyTypeHash)},
					{AttributeName: aws.String("sk1"), KeyType: aws.String(dynamodb.KeyTypeRange)},
				},
				Projection: &dynamodb.Projection{ProjectionType: aws.String(dynamodb.ProjectionTypeKeysOnly)},
				ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
					ReadCapacityUnits:  aws.Int64(1),
					WriteCapacityUnits: aws.Int64(1),
				},
			},
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
			if aerr.Code() == dynamodb.ErrCodeResourceInUseException {
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
