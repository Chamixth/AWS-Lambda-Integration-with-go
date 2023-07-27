package apis

import (
	"context"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/gofiber/fiber/v2"
)

type GetFunctionCodeSigningConfigWrapper struct {
	LambdaClient *lambda.Client
}

// GetFunctionCodeSigningConfigResponse represents the JSON response for retrieving the code signing configuration of a function.
type GetFunctionCodeSigningConfigResponse struct {
	FunctionName         string `json:"functionName"`
	CodeSigningConfigArn string `json:"codeSigningConfigArn"`
}

// GetFunctionCodeSigningConfig retrieves the code signing configuration of an AWS Lambda function.
func (wrapper GetFunctionCodeSigningConfigWrapper) GetFunctionCodeSigningConfig(functionName string) (*GetFunctionCodeSigningConfigResponse, error) {
	input := &lambda.GetFunctionCodeSigningConfigInput{
		FunctionName: &functionName,
	}

	output, err := wrapper.LambdaClient.GetFunctionCodeSigningConfig(context.Background(), input)
	if err != nil {
		return nil, err
	}

	return &GetFunctionCodeSigningConfigResponse{
		FunctionName:         *output.FunctionName,
		CodeSigningConfigArn: *output.CodeSigningConfigArn,
	}, nil

}

func HandleGetFunctionCodeSigningConfig(c *fiber.Ctx) error {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return err
	}
	lambdaClient := lambda.NewFromConfig(cfg)
	functionName := c.Query("functionName")
	wrapper := GetFunctionCodeSigningConfigWrapper{
		LambdaClient: lambdaClient,
	}
	result, err := wrapper.GetFunctionCodeSigningConfig(functionName)
	if err != nil {
		return err
	}
	return c.Status(http.StatusOK).JSON(result)
}
