package gobot

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

func GetSecret(key string) *string {
	svc := secretsmanager.New(session.Must(session.NewSession()), aws.NewConfig().WithRegion("us-east-1"))
	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(key),
	}

	result, err := svc.GetSecretValue(input)
	if err != nil {
		fmt.Printf("Error in GetSecret: %v\n", err)
		panic(err)

	}
	return result.SecretString
}
