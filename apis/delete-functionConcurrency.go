package apis

import (
	"context"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/gofiber/fiber/v2"
)

type DeleteFunctionConcurrencyWrapper struct {
	LambdaClient *lambda.Client
}

// DeleteFunctionConcurrency removes the function concurrency settings for an AWS Lambda function.
func (wrapper DeleteFunctionConcurrencyWrapper) DeleteFunctionConcurrency(functionName string) error {
	input := &lambda.DeleteFunctionConcurrencyInput{
		FunctionName: &functionName,
	}

	_, err := wrapper.LambdaClient.DeleteFunctionConcurrency(context.Background(), input)
	return err
}

func HandleDeleteFunctionConcurrency(c *fiber.Ctx) error {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return err
	}

	lambdaClient := lambda.NewFromConfig(cfg)

	wrapper := DeleteFunctionConcurrencyWrapper{
		LambdaClient: lambdaClient,
	}

	functionName := c.Query("functionName")

	err = wrapper.DeleteFunctionConcurrency(functionName)
	if err != nil {
		return err
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Function concurrency removed successfully",
	})
}
