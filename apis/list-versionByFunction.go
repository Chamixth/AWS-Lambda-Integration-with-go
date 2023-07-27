package apis

import (
	"context"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
	"github.com/gofiber/fiber/v2"
)

type ListVersionByFunctionWrapper struct {
	LambdaClient *lambda.Client
}

func (wrapper ListVersionByFunctionWrapper) ListVersionByFunction(functionName string, maxItems int) ([]types.FunctionConfiguration, error) {
	input := &lambda.ListVersionsByFunctionInput{
		FunctionName: aws.String(functionName),
		MaxItems:     aws.Int32(int32(maxItems)),
	}
	result, err := wrapper.LambdaClient.ListVersionsByFunction(context.Background(), input)

	if err != nil {
		return nil, err
	}

	return result.Versions, err
}
func HandleListVersionByFunction(c *fiber.Ctx) error {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return err
	}
	lambdaClient := lambda.NewFromConfig(cfg)

	wrapper := ListVersionByFunctionWrapper{
		LambdaClient: lambdaClient,
	}

	functionName := c.Query("functionName")
	maxItems := 60

	result, err := wrapper.ListVersionByFunction(functionName, maxItems)

	if err != nil {
		return err
	}

	return c.Status(http.StatusOK).JSON(result)
}
