package apis

import (
	"context"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
	"github.com/gofiber/fiber/v2"
	
)

type FunctionEventInvokeConfigRequest struct {
	FunctionName      string                  `json:"functionName"`
	AliasName string `json:"aliasName`
	MaxAgeOfEvent     int32                   `json:"maxAgeOfEvent"`
	MaxRetryAttempts  int32                   `json:"maxRetryAttempts"`
	DestinationConfig types.DestinationConfig `json:"destinationConfig"`
}

type UpdateFunctionEventInvokeConfigWrapper struct {
	LambdaClient *lambda.Client
}

func (wrapper UpdateFunctionEventInvokeConfigWrapper) UpdateFunctionEventInvokeConfig(request FunctionEventInvokeConfigRequest) error {
	input := &lambda.UpdateFunctionEventInvokeConfigInput{
		FunctionName:             &request.FunctionName,
		Qualifier: &request.AliasName ,
		MaximumEventAgeInSeconds: &request.MaxAgeOfEvent,
		MaximumRetryAttempts:     &request.MaxRetryAttempts,
		DestinationConfig:        &request.DestinationConfig,
	}
	_, err := wrapper.LambdaClient.UpdateFunctionEventInvokeConfig(context.Background(), input)

	return err
}

func HandleUpdateFunctionEventInvokeConfig(c *fiber.Ctx) error {
	var request FunctionEventInvokeConfigRequest

	cfg, err := config.LoadDefaultConfig(context.TODO())

	if err != nil {
		return nil
	}

	lambdaClient := lambda.NewFromConfig(cfg)

	wrapper := UpdateFunctionEventInvokeConfigWrapper{
		LambdaClient: lambdaClient,
	}

	if err := c.BodyParser(&request); err != nil {
		return err
	}

	if err := wrapper.UpdateFunctionEventInvokeConfig(request); err != nil {
		return err
	}

	return c.Status(http.StatusOK).JSON("Update FunctionEventInvoke is Successful")
}
