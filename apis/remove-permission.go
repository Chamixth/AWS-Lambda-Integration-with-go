package apis

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/gofiber/fiber/v2"
)

type RemovePermissionWrapper struct {
	LambdaClient *lambda.Client
}

// RemovePermission removes an existing permission from an AWS Lambda function.
func (wrapper RemovePermissionWrapper) RemovePermission(functionName, statementId string) error {
	input := &lambda.RemovePermissionInput{
		FunctionName: &functionName,
		StatementId:  &statementId,
	}

	_, err := wrapper.LambdaClient.RemovePermission(context.Background(), input)
	return err
}

func HandleRemovePermission(c *fiber.Ctx) error {
	functionName := c.Query("functionName")
	statementId := c.Query("statementId")

	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}
	lambdaClient := lambda.NewFromConfig(cfg)

	// Instantiate FunctionWrapper with the Lambda client
	wrapper := RemovePermissionWrapper{
		LambdaClient: lambdaClient,
	}

	err = wrapper.RemovePermission(functionName, statementId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to remove permission: %v", err),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Permission removed successfully",
	})
}
