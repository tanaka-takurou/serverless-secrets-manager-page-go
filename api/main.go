package main

import (
	"os"
	"fmt"
	"log"
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

type APIResponse struct {
	Message string `json:"message"`
}

type SecretData struct {
	SecretName   string `json:"SecretName"`
	SecretString string `json:"SecretString"`
	CreatedDate  string `json:"CreatedDate"`
}

type Response events.APIGatewayProxyResponse

const layout  string = "2006-01-02 15:04"
var secretsmanagerClient *secretsmanager.Client

func HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (Response, error) {
	var jsonBytes []byte
	var err error
	d := make(map[string]string)
	json.Unmarshal([]byte(request.Body), &d)
	if v, ok := d["action"]; ok {
		switch v {
		case "getsecret" :
			res, e := getSecret(ctx)
			if e != nil {
				err = e
			} else {
				jsonBytes, _ = json.Marshal(APIResponse{Message: res})
			}
		}
	}
	log.Print(request.RequestContext.Identity.SourceIP)
	if err != nil {
		log.Print(err)
		jsonBytes, _ = json.Marshal(APIResponse{Message: fmt.Sprint(err)})
		return Response{
			StatusCode: 500,
			Body: string(jsonBytes),
		}, nil
	}
	return Response {
		StatusCode: 200,
		Body: string(jsonBytes),
	}, nil
}

func getSecret(ctx context.Context)(string, error) {
	if secretsmanagerClient == nil {
		secretsmanagerClient = secretsmanager.NewFromConfig(getConfig(ctx))
	}
	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(os.Getenv("SECRET_ID")),
	}

	result, err := secretsmanagerClient.GetSecretValue(ctx, input)
	if err != nil {
		log.Print(err)
		return "", err
	}
	resultJson, err := json.Marshal(SecretData{
		SecretName: aws.ToString(result.Name),
		SecretString: aws.ToString(result.SecretString),
		CreatedDate: result.CreatedDate.Format(layout),
	})
	if err != nil {
		log.Print(err)
		return "", err
	}
	return string(resultJson), nil
}

func getConfig(ctx context.Context) aws.Config {
	var err error
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(os.Getenv("REGION")))
	if err != nil {
		log.Print(err)
	}
	return cfg
}

func main() {
	lambda.Start(HandleRequest)
}
