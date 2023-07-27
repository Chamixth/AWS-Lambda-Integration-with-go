package apis

import (
	"context"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/gofiber/fiber/v2"
)

type DeleteFunctionCodeSigningConfigWrapper struct {
	LambdaClient *lambda.Client
}

func (wrapper DeleteFunctionCodeSigningConfigWrapper) DeleteFunctionCodeSigningConfig(functionName string) error {
	input := &lambda.DeleteFunctionCodeSigningConfigInput{
		FunctionName: &functionName,
	}

	_, err := wrapper.LambdaClient.DeleteFunctionCodeSigningConfig(context.Background(), input)

	if err != nil {
		return err
	}

	return err
}

func HandleDeleteFunctionCodeSigningConfig(c *fiber.Ctx) error {
	cfg, err := config.LoadDefaultConfig(context.TODO())

	if err != nil {
		return err
	}

	lambdaClient := lambda.NewFromConfig(cfg)

	wrapper := DeleteFunctionCodeSigningConfigWrapper{
		LambdaClient: lambdaClient,
	}

	functionName := c.Query("functionName")
	if err := wrapper.DeleteFunctionCodeSigningConfig(functionName); err != nil {
		return err
	}

	return c.Status(http.StatusOK).JSON("FunctionCodeSigningConfig deleted Successfully")
}
