package apis

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/gofiber/fiber/v2"
	
)

type PublishVersionWrapper struct {
	LambdaClient *lambda.Client
}

// PublishVersion publishes a new version of an AWS Lambda function.
func (wrapper PublishVersionWrapper) PublishVersion(functionName string) (string, error) {
	input := &lambda.PublishVersionInput{
		FunctionName: &functionName,
		
	}

	result, err := wrapper.LambdaClient.PublishVersion(context.Background(), input)
	if err != nil {
		return "", err
	}

	return *result.Version, nil
}

func HandlePublishVersion(c *fiber.Ctx) error {
	functionName := c.Query("functionName")

	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}
	lambdaClient := lambda.NewFromConfig(cfg)

	// Instantiate FunctionWrapper with the Lambda client
	wrapper := PublishVersionWrapper{
		LambdaClient: lambdaClient,
	}

	newVersion, err := wrapper.PublishVersion(functionName)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message":    "Function version published successfully",
		"newVersion": newVersion,
	})
}
