package apis

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
	"github.com/gofiber/fiber/v2"
)

type CreateFunctionUrlConfig struct {
	LambdaClient *lambda.Client
}



// CreateFunctionUrlConfig creates a new function URL configuration for the specified Lambda function.
func (wrapper CreateFunctionUrlConfig) CreateFunctionUrlConfig(functionName,aliasName string) error {
	input := &lambda.CreateFunctionUrlConfigInput{
		FunctionName: aws.String(functionName),
		Qualifier:    aws.String(aliasName),
		AuthType:     types.FunctionUrlAuthTypeNone,

	}

	_, err := wrapper.LambdaClient.CreateFunctionUrlConfig(context.Background(), input)
	return err
}

func HandleCreateFunctionUrlConfig(c *fiber.Ctx) error {
	// Load the AWS configuration from environment variables or AWS configuration files
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatalf("Failed to load AWS configuration: %v", err)
	}

	// Create a new Lambda service client
	lambdaClient := lambda.NewFromConfig(cfg)

	wrapper := CreateFunctionUrlConfig{
		LambdaClient: lambdaClient,
	}
	functionName := c.Query("functionName")
	version := c.Query("aliasName")

	if err := wrapper.CreateFunctionUrlConfig(functionName,version); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Function URL configuration created successfully",
	})

}
