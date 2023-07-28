package apis

import (
	"context"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
	"github.com/gofiber/fiber/v2"
)

type ListProvisionedConcurrencyConfigsWrapper struct {
	LambdaClient *lambda.Client
}

type ListProvisionedConcurrencyConfigsResponse struct {
	NextMarker                    *string                                      `json:"nextMarker"`
	ProvisionedConcurrencyConfigs []types.ProvisionedConcurrencyConfigListItem `json:"provisionedConcurrencyConfigs"`
}

func (wrapper ListProvisionedConcurrencyConfigsWrapper) ListProvisionedConcurrencyConfigs(functionName string) (*ListProvisionedConcurrencyConfigsResponse, error) {
	input := &lambda.ListProvisionedConcurrencyConfigsInput{
		FunctionName: &functionName,
	}

	result, err := wrapper.LambdaClient.ListProvisionedConcurrencyConfigs(context.Background(), input)

	if err != nil {
		return nil, err
	}

	response := &ListProvisionedConcurrencyConfigsResponse{
		NextMarker:                    result.NextMarker,
		ProvisionedConcurrencyConfigs: result.ProvisionedConcurrencyConfigs,
	}

	return response, nil
}

func HandleListProvisionedConcurrencyConfigs(c *fiber.Ctx) error {
	cfg, err := config.LoadDefaultConfig(context.TODO())

	if err != nil {
		return err
	}

	lambdaClient := lambda.NewFromConfig(cfg)

	wrapper := ListProvisionedConcurrencyConfigsWrapper{
		LambdaClient: lambdaClient,
	}

	functionName := c.Query("functionName")

	result, err := wrapper.ListProvisionedConcurrencyConfigs(functionName)

	if err != nil {
		return err
	}

	return c.Status(http.StatusOK).JSON(result)
}
