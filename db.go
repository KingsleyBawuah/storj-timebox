package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

const (
	fileKey       = "fileKey"
	downloadCount = "downloadCount"
)

func initDynamoDB(region, endpoint string) *dynamodb.DynamoDB {
	sess, err := session.NewSession(&aws.Config{
		Region:   aws.String(region),
		Endpoint: aws.String(endpoint)})
	if err != nil {
		log.Fatalf("could not create dynamodb session %s", err)
	}

	log.Println("Initializing Dynamodb Session")

	dbSvc := dynamodb.New(sess)

	return dbSvc
}

func (s *server) ensureTables(tableName string) {
	// List to ensure the table used for our app exists. Create it if it doesn't. Set a timeout to prevent waiting forever.
	timeoutCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	result, err := s.DB.ListTablesWithContext(timeoutCtx, &dynamodb.ListTablesInput{})
	if err != nil {
		log.Println(err)
		return
	}

	log.Println(result)

	var exists bool

	for _, name := range result.TableNames {
		if aws.StringValue(name) == tableName {
			exists = true
		}
	}

	if exists == false {
		input := &dynamodb.CreateTableInput{
			AttributeDefinitions: []*dynamodb.AttributeDefinition{
				{
					AttributeName: aws.String(fileKey),
					AttributeType: aws.String("S"),
				},
			},
			KeySchema: []*dynamodb.KeySchemaElement{
				{
					AttributeName: aws.String(fileKey),
					KeyType:       aws.String("HASH"),
				},
			},
			ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
				ReadCapacityUnits:  aws.Int64(10),
				WriteCapacityUnits: aws.Int64(10),
			},
			TableName: aws.String(tableName),
		}

		_, err := s.DB.CreateTable(input)
		if err != nil {
			log.Fatalf("Got error calling CreateTable: %s", err)
		}

		fmt.Println("Created the table", tableName)
	}

}

// Retrieves the amount of times a file has been downloaded
func (s *server) GetDownloadCount(key, tableName string) (int, error) {
	itemOut, err := s.DB.GetItem(&dynamodb.GetItemInput{
		ProjectionExpression:     aws.String(downloadCount),
		ConsistentRead:           nil,
		ExpressionAttributeNames: nil,
		Key: map[string]*dynamodb.AttributeValue{
			fileKey: {S: aws.String(key)},
		},
		TableName: aws.String(tableName),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			return 0, aerr
		}
	}

	countVal := aws.StringValue(itemOut.Item[downloadCount].N) // TODO: Handle this warning.
	if countVal == "" {
		return 0, errors.New("error file doesn't have download count set")
	}

	count, err := strconv.Atoi(countVal)
	log.Println("DOWNLOAD COUNT", count, "AND countVal", countVal)

	return count, nil
}

// Increments the download count for a specific file
func (s *server) IncrementDownloadCount(key, tableName string) error {
	// TODO: This is creating duplicate values somehow.
	newCount, err := s.DB.UpdateItem(&dynamodb.UpdateItemInput{
		ConditionExpression: nil,
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":inc": {
				N: aws.String("1"),
			},
			":start": {
				N: aws.String("0"),
			},
		},
		Key: map[string]*dynamodb.AttributeValue{
			fileKey: {S: aws.String(key)},
		},
		ReturnValues:     aws.String("UPDATED_NEW"),
		TableName:        aws.String(tableName),
		UpdateExpression: aws.String(fmt.Sprintf("SET %s = if_not_exists(%s, :start) + :inc", downloadCount, downloadCount)),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			return aerr
		}
	}

	fmt.Println("NEW UPDATED COUNT FROM DYNAMO", aws.StringValue(newCount.Attributes[downloadCount].N))

	return nil
}

// Create a row representing the download count for a specific file.
func (s *server) CreateDownloadCount(key string, tableName string) error {
	_, err := s.DB.PutItem(&dynamodb.PutItemInput{

		Item: map[string]*dynamodb.AttributeValue{
			fileKey:       {S: aws.String(key)},
			downloadCount: {N: aws.String("0")},
		},
		TableName: aws.String(tableName),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			return aerr
		}
	}

	return nil
}
