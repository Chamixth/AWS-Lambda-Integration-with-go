package apis

import (
	"context"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/gofiber/fiber/v2"
)

type PutFunctionEventInvokeConfigWrapper struct {
	LambdaClient *lambda.Client
}

type PutFunctionEventInvokeConfigRequest struct {
	FunctionName              string `json:"functionName"`
	MaximumRetryAttemptsEvent *int32 `json:"maximumRetryAttemptsEvent"`
	MaximumEventAgeInSeconds  *int32 `json:"maximumEventAgeInSecondsEvent"`
}



func (wrapper PutFunctionEventInvokeConfigWrapper) PutFunctionEventInvokeConfig(request PutFunctionEventInvokeConfigRequest) ( error) {
	input := &lambda.PutFunctionEventInvokeConfigInput{
		FunctionName:             &request.FunctionName,
		MaximumEventAgeInSeconds: request.MaximumEventAgeInSeconds,
		MaximumRetryAttempts:     request.MaximumRetryAttemptsEvent,
	}

	_, err := wrapper.LambdaClient.PutFunctionEventInvokeConfig(context.Background(), input)
	if err != nil {
		return  err
	}

	

	return err
}

func HandlePutFunctionEventInvokeConfig(c *fiber.Ctx) error {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return err
	}

	lambdaClient := lambda.NewFromConfig(cfg)

	var request PutFunctionEventInvokeConfigRequest
	if err := c.BodyParser(&request); err != nil {
		return err
	}

	wrapper := PutFunctionEventInvokeConfigWrapper{
		LambdaClient: lambdaClient,
	}

	if err := wrapper.PutFunctionEventInvokeConfig(request); err != nil {
		return err
	}
	

	return c.Status(http.StatusOK).JSON("PutFunctionEventInvokeConfig is successfull")
}
