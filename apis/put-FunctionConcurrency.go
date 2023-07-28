package apis

import (
	"context"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/gofiber/fiber/v2"
)

type PutFunctionConcurrencyWrapper struct {
	LambdaClient *lambda.Client
}

func (wrapper PutFunctionConcurrencyWrapper) PutFunctionConcurrency(functionName string, reservedConcurrentExecutions int32) error {
	input := &lambda.PutFunctionConcurrencyInput{
		FunctionName:                 &functionName,
		ReservedConcurrentExecutions: &reservedConcurrentExecutions,
	}

	_, err := wrapper.LambdaClient.PutFunctionConcurrency(context.Background(), input)

	if err != nil {
		return err
	}

	return err
}

func HandlePutFunctionConcurrency(c *fiber.Ctx) error {
	cfg, err := config.LoadDefaultConfig(context.TODO())

	if err != nil {
		return err
	}

	lambdaClient := lambda.NewFromConfig(cfg)

	wrapper := PutFunctionConcurrencyWrapper{
		LambdaClient: lambdaClient,
	}

	functionName := c.Query("functionName")
	reservedConcurrentExecutions := c.QueryInt("reservedConcurrentExecutions")

	if err := wrapper.PutFunctionConcurrency(functionName, int32(reservedConcurrentExecutions)); err != nil {
		return err
	}

	return c.Status(http.StatusOK).JSON("PutFunctionConcurrency is successfull")
}
