package apis

import (
	"context"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/gofiber/fiber/v2"
)

type DeleteFunctionEventInvokeConfigWrapper struct {
	LambdaClient *lambda.Client
}

func (wrapper DeleteFunctionEventInvokeConfigWrapper) DeleteFunctionEventInvokeConfig(functionName string) error {
	input := &lambda.DeleteFunctionEventInvokeConfigInput{
		FunctionName: &functionName,
	}
	_, err := wrapper.LambdaClient.DeleteFunctionEventInvokeConfig(context.Background(), input)

	if err != nil {
		return err
	}
	return err
}

func HandleDeleteFunctionEventInvokeConfig(c *fiber.Ctx) error {
	cfg, err := config.LoadDefaultConfig(context.TODO())

	if err != nil {
		return err
	}

	lambdaClient := lambda.NewFromConfig(cfg)

	wrapper := DeleteFunctionEventInvokeConfigWrapper{
		LambdaClient: lambdaClient,
	}

	functionName := c.Query("functionName")

	if err := wrapper.DeleteFunctionEventInvokeConfig(functionName); err != nil {
		return err
	}

	return c.Status(http.StatusOK).JSON("FunctionEventInvokeConfig deleted Successfully")
}
