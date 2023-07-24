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



// FunctionWrapper encapsulates function actions used in the examples.
// It contains an AWS Lambda service client that is used to perform user actions.
type UpdateConfigFunctionWrapper struct {
	LambdaClient *lambda.Client
}

type Payload struct {
	FunctionName string            `json:"functionName"`
	EnvVars      map[string]string `json:"envVars"`
}


// UpdateFunctionConfiguration updates a map of environment variables configured for
// the Lambda function specified by functionName.
func (wrapper UpdateConfigFunctionWrapper) UpdateFunctionConfiguration(functionName string, envVars map[string]string) {
	_, err := wrapper.LambdaClient.UpdateFunctionConfiguration(context.TODO(), &lambda.UpdateFunctionConfigurationInput{
		FunctionName: aws.String(functionName),
		Environment: &types.Environment{
			Variables: envVars,
		},
	})
	if err != nil {
		log.Panicf("Couldn't update configuration for %v. Here's why: %v", functionName, err)
	}
}

func HandleUpdateConfigFunction(c *fiber.Ctx)error{
	// Create an AWS Lambda service client
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}
	lambdaClient := lambda.NewFromConfig(cfg)

	// Instantiate FunctionWrapper with the Lambda client
	wrapper := UpdateConfigFunctionWrapper{
		LambdaClient: lambdaClient,
	}
	

	var payload Payload
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	wrapper.UpdateFunctionConfiguration(payload.FunctionName, payload.EnvVars)
	return c.JSON(fiber.Map{"message": "function configuration updated successfully"})
}

