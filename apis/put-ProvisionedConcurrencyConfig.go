package apis

import (
	"context"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/gofiber/fiber/v2"
)

type PutProvisionedConcurrencyConfigWrapper struct {
	LambdaClient *lambda.Client
}

type PutProvisionedConcurrencyConfigRequest struct {
	FunctionName       string `json:"functionName"`
	ProvisionedConcurrentExecutions *int32  `json:"provisionedConcurrentExecutions"`
	Qualifier string `json:"qualifier"`
}

type PutProvisionedConcurrencyConfigResponse struct {
	FunctionName string `json:"functionName"`
}

func (wrapper PutProvisionedConcurrencyConfigWrapper) PutProvisionedConcurrencyConfig(request PutProvisionedConcurrencyConfigRequest) (*PutProvisionedConcurrencyConfigResponse, error) {
	input := &lambda.PutProvisionedConcurrencyConfigInput{
		FunctionName:                &request.FunctionName,
		ProvisionedConcurrentExecutions: request.ProvisionedConcurrentExecutions,
		Qualifier: &request.Qualifier,
	}

	_, err := wrapper.LambdaClient.PutProvisionedConcurrencyConfig(context.Background(), input)
	if err != nil {
		return nil, err
	}

	response := &PutProvisionedConcurrencyConfigResponse{
		FunctionName: request.FunctionName,
	}

	return response, nil
}

func HandlePutProvisionedConcurrencyConfig(c *fiber.Ctx) error {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return err
	}

	lambdaClient := lambda.NewFromConfig(cfg)

	var request PutProvisionedConcurrencyConfigRequest
	if err := c.BodyParser(&request); err != nil {
		return err
	}

	wrapper := PutProvisionedConcurrencyConfigWrapper{
		LambdaClient: lambdaClient,
	}

	result, err := wrapper.PutProvisionedConcurrencyConfig(request)
	if err != nil {
		return err
	}

	return c.Status(http.StatusOK).JSON(result)
}
