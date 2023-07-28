package apis

import (
	"context"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/gofiber/fiber/v2"
)

type DeleteProvisionedConcurrencyConfigWrapper struct {
	LambdaClient *lambda.Client
}

// DeleteProvisionedConcurrencyConfig removes the provisioned concurrency configuration for an AWS Lambda function.
func (wrapper DeleteProvisionedConcurrencyConfigWrapper) DeleteProvisionedConcurrencyConfig(functionName string, qualifier string) error {
	input := &lambda.DeleteProvisionedConcurrencyConfigInput{
		FunctionName: &functionName,
		Qualifier:    &qualifier,
	}

	_, err := wrapper.LambdaClient.DeleteProvisionedConcurrencyConfig(context.Background(), input)
	return err
}

func HandleDeleteProvisionedConcurrencyConfig(c *fiber.Ctx) error {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return err
	}

	lambdaClient := lambda.NewFromConfig(cfg)

	wrapper := DeleteProvisionedConcurrencyConfigWrapper{
		LambdaClient: lambdaClient,
	}

	functionName := c.Query("functionName")
	qualifier := c.Query("qualifier")

	err = wrapper.DeleteProvisionedConcurrencyConfig(functionName, qualifier)
	if err != nil {
		return err
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Provisioned concurrency configuration removed successfully",
	})
}
