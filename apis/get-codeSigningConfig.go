package apis

import (
	"context"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
	"github.com/gofiber/fiber/v2"
)

type GetCodeSigningConfigWrapper struct {
	LambdaClient *lambda.Client
}

func (wrapper GetCodeSigningConfigWrapper) GetCodeSigingConfig(codeSigningConfigArn string) (*types.CodeSigningConfig, error) {
	input := &lambda.GetCodeSigningConfigInput{
		CodeSigningConfigArn: &codeSigningConfigArn,
	}
	output, err := wrapper.LambdaClient.GetCodeSigningConfig(context.Background(), input)

	if err != nil {
		return nil, err
	}
	return output.CodeSigningConfig, nil
}

func HandleGetCodeSigningConfig(c *fiber.Ctx) error {
	cfg, err := config.LoadDefaultConfig(context.TODO())

	if err != nil {
		return err
	}

	lambdaClient := lambda.NewFromConfig(cfg)

	wrapper := GetCodeSigningConfigWrapper{
		LambdaClient: lambdaClient,
	}

	codeSigningConfigArn := c.Query("codeSigningConfigArn")

	result, err := wrapper.GetCodeSigingConfig(codeSigningConfigArn)

	if err != nil {
		return err
	}

	return c.Status(http.StatusOK).JSON(result)

}
