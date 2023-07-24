package apis

import (
	"context"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
	"github.com/gofiber/fiber/v2"
)

type GetFunctionWrapper struct {
	LambdaClient *lambda.Client
}

func (wrapper GetFunctionWrapper) GetFunction(functionName string) (*types.FunctionConfiguration, error) {
	funcOutput, err := wrapper.LambdaClient.GetFunction(context.TODO(), &lambda.GetFunctionInput{
		FunctionName: aws.String(functionName),
	})
	if err != nil {
		return nil, err
	}
	return funcOutput.Configuration, nil
}

func HandleGetFunction(c *fiber.Ctx) error {
	functionName := c.Query("functionName")
	if functionName == "" {
		c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Function name is invalid"})
		return fiber.ErrBadRequest
	}

	// Create AWS Lambda service client using default configuration
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatalf("failed to load AWS configuration: %v", err)
	}
	lambdaClient := lambda.NewFromConfig(cfg)

	wrapper := GetFunctionWrapper{
		LambdaClient: lambdaClient,
	}

	functionConfig, err := wrapper.GetFunction(functionName)
	if err != nil {
		c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "Function not found"})
		return fiber.ErrBadRequest
	}

	// Function found, handle functionConfig as needed
	// For example, you can extract function details such as ARN, state, etc.

	return c.JSON(fiber.Map{
		"functionName": functionName,
		"state":        functionConfig.State,
		// ... other function details you want to include in the response
	})
}
