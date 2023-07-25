package apis

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/gofiber/fiber/v2"
)

// FunctionWrapper encapsulates function actions used in the examples.
// It contains an AWS Lambda service client that is used to perform user actions.
type GetAliasFunctionWrapper struct {
	LambdaClient *lambda.Client
}

// GetAlias retrieves the information about an alias for the specified Lambda function.
func (wrapper GetAliasFunctionWrapper) GetAlias(functionName, aliasName string) (*lambda.GetAliasOutput, error) {
	input := &lambda.GetAliasInput{
		FunctionName: aws.String(functionName),
		Name:         aws.String(aliasName),
	}

	return wrapper.LambdaClient.GetAlias(context.TODO(), input)
}
func HandleGetAliasFunction(c *fiber.Ctx) error {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatalf("Failed to load AWS configuration: %v", err)
	}

	// Create a new Lambda service client
	lambdaClient := lambda.NewFromConfig(cfg)

	// Initialize the FunctionWrapper with the Lambda service client
	wrapper := GetAliasFunctionWrapper{
		LambdaClient: lambdaClient,
	}

	functionName := c.Query("functionName")
	aliasName := c.Query("aliasName")

	aliasOutput, err := wrapper.GetAlias(functionName, aliasName)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(aliasOutput)
}
