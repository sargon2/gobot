package gobot

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

type MapquestApiKey string
type WolframAlphaKey string
type SlackBotToken string

func ProvideMapquestApiKey() (*MapquestApiKey, error) {
	apiKey, err := GetSecret("mapquest_api_key")
	if err != nil {
		return nil, err
	}
	result := MapquestApiKey(*apiKey)
	return &result, nil
}

func ProvideWolframAlphaKey() (*WolframAlphaKey, error) {
	apiKey, err := GetSecret("wolfram_alpha_key")
	if err != nil {
		return nil, err
	}
	result := WolframAlphaKey(*apiKey)
	return &result, nil
}

func ProvideSlackBotToken() (*SlackBotToken, error) {
	apiKey, err := GetSecret("slack_bot_token")
	if err != nil {
		return nil, err
	}
	result := SlackBotToken(*apiKey)
	return &result, nil
}

func GetSecret(key string) (*string, error) {
	svc := secretsmanager.New(session.Must(session.NewSession()), aws.NewConfig().WithRegion("us-east-1"))
	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(key),
	}

	result, err := svc.GetSecretValue(input)
	if err != nil {
		fmt.Printf("Error in GetSecret: %v\n", err)
		return nil, err

	}
	return result.SecretString, nil
}
