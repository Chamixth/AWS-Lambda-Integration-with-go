package apis

import (
	"context"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
	"github.com/gofiber/fiber/v2"
)

type ListCodeSigningConfigWrapper struct {
	LambdaClient *lambda.Client
}

func (wrapper ListCodeSigningConfigWrapper) ListCodeSigningConfig() (*[]types.CodeSigningConfig, error) {
	input := &lambda.ListCodeSigningConfigsInput{}

	output, err := wrapper.LambdaClient.ListCodeSigningConfigs(context.Background(), input)

	if err != nil {
		return nil, err
	}

	return &output.CodeSigningConfigs, nil
}

func HandleListCodeSigningConfig(c *fiber.Ctx) error {
	cfg, err := config.LoadDefaultConfig(context.TODO())

	if err != nil {
		return err
	}

	lambdaClient := lambda.NewFromConfig(cfg)

	wrapper := ListCodeSigningConfigWrapper{
		LambdaClient: lambdaClient,
	}

	result, err := wrapper.ListCodeSigningConfig()

	if err != nil {
		return err
	}

	return c.Status(http.StatusOK).JSON(result)
}
