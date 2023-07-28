package apis

import (
	"context"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
	"github.com/gofiber/fiber/v2"
)

type ListFunctionEventInvokeConfigsWrapper struct {
	LambdaClient *lambda.Client
}

type ListFunctionEventInvokeConfigsResponse struct {
	FunctionEventInvokeConfigs *[]types.FunctionEventInvokeConfig `json:functionEventInvokeConfigs`
	NextMarker                 *string                            `json:"nextMarker"`
}

func (wrapper ListFunctionEventInvokeConfigsWrapper) ListFunctionEventInvokeConfigs(functionName string) (*ListFunctionEventInvokeConfigsResponse, error) {
	input := &lambda.ListFunctionEventInvokeConfigsInput{
		FunctionName: &functionName,
	}

	result, err := wrapper.LambdaClient.ListFunctionEventInvokeConfigs(context.Background(), input)

	if err != nil {
		return nil, err
	}

	response := &ListFunctionEventInvokeConfigsResponse{
		FunctionEventInvokeConfigs: &result.FunctionEventInvokeConfigs,
		NextMarker:                 result.NextMarker,
	}

	return response, nil
}

func HandleListFunctionEventInvokeConfigs(c *fiber.Ctx) error {
	cfg, err := config.LoadDefaultConfig(context.TODO())

	if err != nil {
		return err
	}

	lambdaClient := lambda.NewFromConfig(cfg)

	wrapper := ListFunctionEventInvokeConfigsWrapper{
		LambdaClient: lambdaClient,
	}

	functionName := c.Query("functionName")

	result, err := wrapper.ListFunctionEventInvokeConfigs(functionName)

	if err != nil {
		return err
	}

	return c.Status(http.StatusOK).JSON(result)
}
