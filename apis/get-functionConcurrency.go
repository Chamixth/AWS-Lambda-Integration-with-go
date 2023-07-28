package apis

import (
	"context"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/gofiber/fiber/v2"
)

type GetFunctionConcurrencyWrapper struct {
	LambdaClient *lambda.Client
}

// GetFunctionConcurrencyResponse represents the JSON response for retrieving the function concurrency settings.
type GetFunctionConcurrencyResponse struct {
	FunctionName string `json:"functionName"`
	Reserved     int32  `json:"reserved"`
	
}

// GetFunctionConcurrency retrieves the function concurrency settings for an AWS Lambda function.
func (wrapper GetFunctionConcurrencyWrapper) GetFunctionConcurrency(functionName string) (*GetFunctionConcurrencyResponse, error) {
	input := &lambda.GetFunctionConcurrencyInput{
		FunctionName: &functionName,
	}

	output, err := wrapper.LambdaClient.GetFunctionConcurrency(context.Background(), input)
	if err != nil {
		return nil, err
	}

	reserved := int32(0) // Default value if ReservedConcurrentExecutions is nil
	if output.ReservedConcurrentExecutions != nil {
		reserved = *output.ReservedConcurrentExecutions
	}

	return &GetFunctionConcurrencyResponse{
		FunctionName: functionName,
		Reserved:     reserved,
	}, nil
}


func HandleGetFunctionConcurrency(c *fiber.Ctx) error {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return err
	}

	lambdaClient := lambda.NewFromConfig(cfg)

	wrapper := GetFunctionConcurrencyWrapper{
		LambdaClient: lambdaClient,
	}

	functionName := c.Query("functionName")

	result, err := wrapper.GetFunctionConcurrency(functionName)
	if err != nil {
		return err
	}

	return c.Status(http.StatusOK).JSON(result)
}
