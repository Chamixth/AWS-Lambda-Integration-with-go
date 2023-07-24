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

type ListFunctionWrapper struct {
	LambdaClient *lambda.Client
}

func (wrapper ListFunctionWrapper) ListFunctions(maxItems int) []types.FunctionConfiguration {
	var functions []types.FunctionConfiguration
	paginator := lambda.NewListFunctionsPaginator(wrapper.LambdaClient, &lambda.ListFunctionsInput{
		MaxItems: aws.Int32(int32(maxItems)),
	})
	for paginator.HasMorePages() && len(functions) < maxItems {
		pageOutput, err := paginator.NextPage(context.Background())
		if err != nil {
			log.Panicf("Couldn't list functions for your account. Here's why: %v\n", err)
		}
		functions = append(functions, pageOutput.Functions...)
	}
	return functions
}

func HandleListFunction(c *fiber.Ctx) error {
	// Initialize AWS config and create the Lambda client
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}
	lambdaClient := lambda.NewFromConfig(cfg)

	// Create the FunctionWrapper
	wrapper := ListFunctionWrapper{
		LambdaClient: lambdaClient,
	}
	maxItems := 1000 // You can change this value to limit the number of items returned
	functions := wrapper.ListFunctions(maxItems)
	return c.JSON(functions)
}
