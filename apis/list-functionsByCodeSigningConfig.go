package apis

import (
	"context"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/gofiber/fiber/v2"
)

type ListFunctionsByCodeSigningConfigWrapper struct {
	LambdaClient *lambda.Client
}

type ListFunctionsByCodeSigningConfigResponse struct {
	FunctionArns *[]string `json:"functionArns"`
	NextMarker   *string   `json:"nextMarker"`
}

func (wrapper ListFunctionsByCodeSigningConfigWrapper) ListFunctionByCodeSigningConfig(codeSigningConfigArn string) (*ListFunctionsByCodeSigningConfigResponse, error) {
	input := &lambda.ListFunctionsByCodeSigningConfigInput{
		CodeSigningConfigArn: &codeSigningConfigArn,
	}

	result, err := wrapper.LambdaClient.ListFunctionsByCodeSigningConfig(context.Background(), input)

	if err != nil {
		return nil, err
	}

	response := &ListFunctionsByCodeSigningConfigResponse{
		FunctionArns: &result.FunctionArns,
		NextMarker:   result.NextMarker,
	}

	return response, nil
}

func HandleListFunctionByCodeSigningConfig(c *fiber.Ctx) error {
	cfg, err := config.LoadDefaultConfig(context.TODO())

	if err != nil {
		return err
	}

	lambdaClient := lambda.NewFromConfig(cfg)

	wrapper := ListFunctionsByCodeSigningConfigWrapper{
		LambdaClient: lambdaClient,
	}

	codeSigningConfigArn := c.Query("codeSigningConfigArn")

	result, err := wrapper.ListFunctionByCodeSigningConfig(codeSigningConfigArn)

	if err != nil {
		return err
	}

	return c.Status(http.StatusOK).JSON(result)
}
