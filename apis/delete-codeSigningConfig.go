package apis

import (
	"context"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/gofiber/fiber/v2"
)

type DeleteCodeSigningConfigWrapper struct {
	LambdaClient *lambda.Client
}

func (wrapper DeleteCodeSigningConfigWrapper) DeleteCodeSigningConfig(codeSigningConfigArn string) error {
	input := &lambda.DeleteCodeSigningConfigInput{
		CodeSigningConfigArn: &codeSigningConfigArn,
	}
	_, err := wrapper.LambdaClient.DeleteCodeSigningConfig(context.Background(), input)

	if err != nil {
		return err
	}

	return err
}

func HandleDeleteCodeSigningConfig(c *fiber.Ctx) error {
	cfg, err := config.LoadDefaultConfig(context.TODO())

	if err != nil {
		return err
	}

	lambdaClient := lambda.NewFromConfig(cfg)

	wrapper := DeleteCodeSigningConfigWrapper{
		LambdaClient: lambdaClient,
	}

	codeSigningConfigArn := c.Query("codeSigningConfigArn")

	if err := wrapper.DeleteCodeSigningConfig(codeSigningConfigArn); err != nil {
		return err
	}

	return c.Status(http.StatusOK).JSON("CodeSigningConfigArn deleted Successfully")
}
