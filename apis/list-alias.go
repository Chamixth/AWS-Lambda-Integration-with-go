package apis

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/gofiber/fiber/v2"
)

type ListAliasFunctionWrapper struct {
	LambdaClient *lambda.Client
}

func (wrapper ListAliasFunctionWrapper) ListAlias(functionName string) ([]types.AliasConfiguration, error) {
	input := &lambda.ListAliasesInput{
		FunctionName: aws.String(functionName),
	}
	aliasOutput, err := wrapper.LambdaClient.ListAliases(context.Background(), input)
	if err != nil {
		return nil, err
	}
	return aliasOutput.Aliases, nil
}

func HandleListAliasesFunction(c *fiber.Ctx) error {
	// Load the AWS configuration from environment variables or AWS configuration files
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatalf("Failed to load AWS configuration: %v", err)
	}

	// Create a new Lambda service client
	lambdaClient := lambda.NewFromConfig(cfg)

	// Initialize the FunctionWrapper with the Lambda service client
	wrapper := ListAliasFunctionWrapper{
		LambdaClient: lambdaClient,
	}
	functionName := c.Query("functionName")

	listAliasesOutput,err := wrapper.ListAlias(functionName)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(listAliasesOutput)
}
