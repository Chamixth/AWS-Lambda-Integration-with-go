package apis

import (
	"context"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/gofiber/fiber/v2"
)

type GetProvisionedConcurrencyConfigWrapper struct {
	LambdaClient *lambda.Client
}

// GetProvisionedConcurrencyConfigResponse represents the JSON response for retrieving the provisioned concurrency configuration of a Lambda function.
type GetProvisionedConcurrencyConfigResponse struct {
	FunctionName                             string `json:"functionName"`
	RequestedConcurrentExecutions            int64  `json:"requestedConcurrentExecutions"`
	AvailableProvisionedConcurrentExecutions int64  `json:"availableProvisionedConcurrentExecutions"`
	AllocatedProvisionedConcurrentExecutions int64  `json:"allocatedProvisionedConcurrentExecutions"`
	Status                                   string `json:"status"`
}

// GetProvisionedConcurrencyConfig retrieves the provisioned concurrency configuration for an AWS Lambda function.
func (wrapper GetProvisionedConcurrencyConfigWrapper) GetProvisionedConcurrencyConfig(functionName,qualifier string) (*GetProvisionedConcurrencyConfigResponse, error) {
	input := &lambda.GetProvisionedConcurrencyConfigInput{
		FunctionName: &functionName,
		Qualifier: &qualifier ,
	}

	output, err := wrapper.LambdaClient.GetProvisionedConcurrencyConfig(context.Background(), input)
	if err != nil {
		return nil, err
	}

	return &GetProvisionedConcurrencyConfigResponse{
		FunctionName:                             functionName,
		RequestedConcurrentExecutions:            int64(*output.RequestedProvisionedConcurrentExecutions),
		AvailableProvisionedConcurrentExecutions: int64(*output.AvailableProvisionedConcurrentExecutions),
		AllocatedProvisionedConcurrentExecutions: int64(*output.AllocatedProvisionedConcurrentExecutions),
		Status:                                   string(output.Status),
	}, nil
}

func HandleGetProvisionedConcurrencyConfig(c *fiber.Ctx) error {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return err
	}

	lambdaClient := lambda.NewFromConfig(cfg)

	wrapper := GetProvisionedConcurrencyConfigWrapper{
		LambdaClient: lambdaClient,
	}

	functionName := c.Query("functionName")
	qualifier := c.Query("qualifier")

	result, err := wrapper.GetProvisionedConcurrencyConfig(functionName,qualifier)
	if err != nil {
		return err
	}

	return c.Status(http.StatusOK).JSON(result)
}
