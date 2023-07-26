package apis

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/gofiber/fiber/v2"
)

type DeleteFunctionUrlConfigWrapper struct {
	LambdaClient *lambda.Client
}

func (wrapper DeleteFunctionUrlConfigWrapper) DeleteFunctionUrlConfig(functionName, aliasName string) error {
	input := &lambda.DeleteFunctionUrlConfigInput{
		FunctionName: aws.String(functionName),
		Qualifier:    aws.String(aliasName),
	}

	_, err := wrapper.LambdaClient.DeleteFunctionUrlConfig(context.Background(), input)

	if err != nil {
		return err
	}

	return err
}

func HandleDeleteFunctionUrlConfig(c *fiber.Ctx) error {
	// Create AWS Lambda service client using default configuration
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatalf("failed to load AWS configuration: %v", err)
	}
	lambdaClient := lambda.NewFromConfig(cfg)

	wrapper := DeleteFunctionUrlConfigWrapper{
		LambdaClient: lambdaClient,
	}

	functionName := c.Query("functionName")
	aliasName := c.Query("aliasName")

	if err := wrapper.DeleteFunctionUrlConfig(functionName,aliasName); err !=nil{
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"message": "Function URL configuration deleted successfully",
	})
}
