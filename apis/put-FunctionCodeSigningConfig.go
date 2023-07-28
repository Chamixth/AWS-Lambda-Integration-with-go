package apis

import (
	"context"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/gofiber/fiber/v2"
)

type PutFunctionCodeSigningConfigWrapper struct {
	LambdaClient *lambda.Client
}

type PutFunctionCodeSigningConfigRequest struct {
	CodeSigningConfigArn string `json:"codeSigningConfigArn"`
	FunctionName         string `json:"functionName"`
}

func (wrapper PutFunctionCodeSigningConfigWrapper) PutFunctionCodeSigningConfig(request PutFunctionCodeSigningConfigRequest) error {
	input := &lambda.PutFunctionCodeSigningConfigInput{
		CodeSigningConfigArn: &request.CodeSigningConfigArn,
		FunctionName:         &request.FunctionName,
	}
	_, err := wrapper.LambdaClient.PutFunctionCodeSigningConfig(context.Background(), input)

	if err != nil {
		return err
	}

	return err
}

func HandlePutFunctionCodeSigningConfig(c *fiber.Ctx) error {
	var request PutFunctionCodeSigningConfigRequest
	cfg, err := config.LoadDefaultConfig(context.TODO())

	if err != nil {
		return err
	}

	lambdaClient := lambda.NewFromConfig(cfg)

	wrapper := PutFunctionCodeSigningConfigWrapper{
		LambdaClient: lambdaClient,
	}

	if err := c.BodyParser(&request); err != nil {
		return err
	}

	if err := wrapper.PutFunctionCodeSigningConfig(request); err != nil {
		return err
	}

	return c.Status(http.StatusOK).JSON("PutFunctionCodeSigningConfig is successfull")
}
