package apis

import (
	"context"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
	"github.com/gofiber/fiber/v2"
)

type GetFunctionEventInvokeConfigWrapper struct {
	LambdaClient *lambda.Client
}

// GetFunctionEventInvokeConfigResponse represents the JSON response for retrieving the event invoke configuration of a function.
type GetFunctionEventInvokeConfigResponse struct {
	FunctionName             string                   `json:"functionName"`
	MaximumRetryAttempts     int32                    `json:"maximumRetryAttempts"`
	MaximumEventAgeInSeconds int32                    `json:"maximumEventAgeInSeconds"`
	DestinationConfig        *types.DestinationConfig `json:"destinationConfig"`
	Qualifier                string                   `json:"qualifier"`
}

func (wrapper GetFunctionEventInvokeConfigWrapper) GetFunctionEventInvokeConfig(functionName string) (*GetFunctionEventInvokeConfigResponse, error) {
	input := &lambda.GetFunctionEventInvokeConfigInput{
		FunctionName: &functionName,
	}
	output, err := wrapper.LambdaClient.GetFunctionEventInvokeConfig(context.Background(), input)

	if err != nil {
		return nil, err
	}

	response := GetFunctionEventInvokeConfigResponse{
		FunctionName:             *output.FunctionArn,
		MaximumRetryAttempts:     *output.MaximumRetryAttempts,
		MaximumEventAgeInSeconds: *output.MaximumEventAgeInSeconds,
		DestinationConfig:        output.DestinationConfig,
		Qualifier:                functionName, // Use the provided functionName instead of dereferencing the input.Qualifier
	}

	if output.DestinationConfig != nil {
		response.DestinationConfig = output.DestinationConfig
	}

	return &response, nil
}

func HandleGetFunctionEventInvokeConfig(c *fiber.Ctx) error {
	cfg, err := config.LoadDefaultConfig(context.Background())

	if err != nil {
		return err
	}

	lambdaClient := lambda.NewFromConfig(cfg)

	wrapper := GetFunctionEventInvokeConfigWrapper{
		LambdaClient: lambdaClient,
	}

	functionName := c.Query("functionName")

	result, err := wrapper.GetFunctionEventInvokeConfig(functionName)

	if err != nil {
		return err
	}

	return c.Status(http.StatusOK).JSON(result)
}
