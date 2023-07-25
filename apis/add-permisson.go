package apis

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/gofiber/fiber/v2"
)

type PermissionRequest struct {
	FunctionName string `json:"functionName"`
	StatementID  string `json:"statementId"`
	Action       string `json:"action"`
	Principal    string `json:"principal"`
	SourceArn    string `json:"sourceArn"`
}

func AddPermissionToLambdaFunction(functionName, statementID, action, principal, sourceArn string) error {
	// Load the AWS configuration from environment variables or AWS configuration files
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return err
	}

	// Create a new Lambda service client
	client := lambda.NewFromConfig(cfg)

	// Set up the input parameters for the AddPermission request
	input := &lambda.AddPermissionInput{
		FunctionName: &functionName,
		StatementId:  &statementID,
		Action:       &action,
		Principal:    &principal,
		SourceArn:    &sourceArn,
	}

	// Make the API call to add permission
	_, err = client.AddPermission(context.Background(), input)
	return err
}

func HandleAddPermissionFunction(c *fiber.Ctx) error {
	// Parse the JSON payload from the request into a PermissionRequest struct
	var request PermissionRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid JSON payload",
		})
	}

	// Call the function to add permission to the Lambda function
	err := AddPermissionToLambdaFunction(request.FunctionName, request.StatementID, request.Action, request.Principal, request.SourceArn)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.SendString("Permission added to Lambda function successfully.")
}
