package apis

import (
	"context"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/gofiber/fiber/v2"
)

type GetRuntimeManagementConfigWrapper struct {
	LambdaClient *lambda.Client
}

type GetRuntimeManagementConfigResponse struct {
	FunctionArn       string `json:"functionArn"`
	UpdateRuntimeOn   string `json:"updateRuntimeOn"`
}

func (wrapper GetRuntimeManagementConfigWrapper) GetRuntimeManagementConfig(functionName string) (*GetRuntimeManagementConfigResponse, error) {
	input := &lambda.GetRuntimeManagementConfigInput{
		FunctionName: &functionName,
	}

	result, err := wrapper.LambdaClient.GetRuntimeManagementConfig(context.Background(), input)

	if err != nil {
		return nil, err
	}



	response := &GetRuntimeManagementConfigResponse{
		FunctionArn:       *result.FunctionArn,
		UpdateRuntimeOn:   string(result.UpdateRuntimeOn),
	}

	return response, nil
}

func HandleGetRuntimeManagementConfig(c *fiber.Ctx) error {
	cfg, err := config.LoadDefaultConfig(context.Background())

	if err != nil {
		return err
	}

	lambdaClient := lambda.NewFromConfig(cfg)

	wrapper := GetRuntimeManagementConfigWrapper{
		LambdaClient: lambdaClient,
	}

	functionName := c.Query("functionName")

	result, err := wrapper.GetRuntimeManagementConfig(functionName)

	if err != nil {
		return err
	}

	return c.Status(http.StatusOK).JSON(result)
}
