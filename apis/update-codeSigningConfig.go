package apis

import (
	"context"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
	"github.com/gofiber/fiber/v2"
)

type UpdateCodeSigningConfigWrapper struct {
	LambdaClient *lambda.Client
}

type UpdateCodeSigningConfigRequest struct {
	CodeSigningConfigArn string                     `json:"codeSigningConfigArn"`
	AllowedPublishers    *types.AllowedPublishers   `json:"allowedPublishers"`
	CodeSigningPolicies  *types.CodeSigningPolicies `json:"codeSigningPolicies"`
	Description          string                     `json:"description"`
}

func (wrapper UpdateCodeSigningConfigWrapper) UpdateCodeSigningConfig(request UpdateCodeSigningConfigRequest) error {
	input := &lambda.UpdateCodeSigningConfigInput{
		CodeSigningConfigArn: &request.CodeSigningConfigArn,
		AllowedPublishers:    request.AllowedPublishers,
		CodeSigningPolicies:  request.CodeSigningPolicies,
		Description:          &request.Description,
	}
	_, err := wrapper.LambdaClient.UpdateCodeSigningConfig(context.Background(), input)

	if err != nil {
		return err
	}

	return err
}

func HandleUpdateCodeSigningConfig(c *fiber.Ctx) error {
	var request UpdateCodeSigningConfigRequest
	cfg, err := config.LoadDefaultConfig(context.Background())

	if err != nil {
		return err
	}

	lambdaClient := lambda.NewFromConfig(cfg)

	wrapper := UpdateCodeSigningConfigWrapper{
		LambdaClient: lambdaClient,
	}

	if err := c.BodyParser(&request); err != nil {
		return err
	}

	if err := wrapper.UpdateCodeSigningConfig(request); err != nil {
		return err
	}

	return c.Status(http.StatusOK).JSON("Update CodeSigningConfig is Successful")
}
