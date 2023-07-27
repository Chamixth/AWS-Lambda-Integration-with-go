package apis

import (
	"context"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
	"github.com/gofiber/fiber/v2"
)

type CreateCodeSigningConfigWrapper struct {
	LambdaClient *lambda.Client
}

// CreateCodeSigningConfigRequest represents the JSON payload for creating a code signing configuration.
type CreateCodeSigningConfigRequest struct {
	AllowedPublishers   *types.AllowedPublishers   `json:"allowedPublishers"`
	CodeSigningPolicies *types.CodeSigningPolicies `json:"codeSigningPolicies"`
	Description         string                     `json:"description"`
}

// CreateCodeSigningConfig creates a code signing configuration for an AWS Lambda function.
func (wrapper CreateCodeSigningConfigWrapper) CreateCodeSigningConfig(request CreateCodeSigningConfigRequest) error {
	input := &lambda.CreateCodeSigningConfigInput{
		AllowedPublishers:   request.AllowedPublishers,
		CodeSigningPolicies: request.CodeSigningPolicies,
		Description:         &request.Description,
	}

	_, err := wrapper.LambdaClient.CreateCodeSigningConfig(context.Background(), input)
	return err
}

func HandleCreateSiginingConfig(c *fiber.Ctx) error {
	cfg, err := config.LoadDefaultConfig(context.TODO())

	if err != nil {
		return err
	}

	lambdaClient := lambda.NewFromConfig(cfg)

	wrapper := CreateCodeSigningConfigWrapper{
		LambdaClient: lambdaClient,
	}

	var request CreateCodeSigningConfigRequest

	if err := c.BodyParser(&request); err != nil {
		return err
	}

	if err := wrapper.CreateCodeSigningConfig(request); err != nil {
		return err
	}

	return c.Status(http.StatusOK).JSON(request)

}
