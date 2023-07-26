package apis

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/gofiber/fiber/v2"
)

type GetFunctionUrlConfig struct {
	LambdaClient *lambda.Client
}

// GetFunctionUrlConfigResponse represents the JSON response for getting the function URL configuration.
type GetFunctionUrlConfigResponse struct {
	FunctionName string `json:"functionName"`
	AliasName      string `json:"version"`
	URL          string `json:"url"`
}

func (wrapper GetFunctionUrlConfig) GetFunctionUrlConfig(functionName, aliasName string) (*GetFunctionUrlConfigResponse, error) {
	input := &lambda.GetFunctionUrlConfigInput{
		FunctionName: aws.String(functionName),
		Qualifier:    aws.String(aliasName),
	}
	urlConfig, err := wrapper.LambdaClient.GetFunctionUrlConfig(context.Background(), input)
	if err != nil {
		return nil, err
	}
	response := &GetFunctionUrlConfigResponse{
		FunctionName: functionName,
		AliasName:      aliasName,
		URL:          *urlConfig.FunctionUrl,
	}

	return response, nil
}

func HandleGetFunctionUrlConfig(c *fiber.Ctx) error {
	// Load the AWS configuration from environment variables or AWS configuration files
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatalf("Failed to load AWS configuration: %v", err)
	}

	// Create a new Lambda service client
	lambdaClient := lambda.NewFromConfig(cfg)

	wrapper := GetFunctionUrlConfig{
		LambdaClient: lambdaClient,
	}
	functionName := c.Query("functionName")
	aliasName := c.Query("aliasName")

	response, err := wrapper.GetFunctionUrlConfig(functionName, aliasName)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(response)
}
