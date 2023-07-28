package apis

import (
	"context"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
	"github.com/gofiber/fiber/v2"
)

type PutRuntimeManagementConfigWrapper struct {
	LambdaClient *lambda.Client
}

type PutRuntimeManagementConfigRequest struct {
	FunctionName string `json:"functionName"`
	Runtime      string `json:"runtime"`
	RuntimeVersionArn string `json:"runTimeVersionArn"`
}


func (wrapper PutRuntimeManagementConfigWrapper) PutRuntimeManagementConfig(request PutRuntimeManagementConfigRequest) error {
	input := &lambda.PutRuntimeManagementConfigInput{
		FunctionName:    &request.FunctionName,
		UpdateRuntimeOn: types.UpdateRuntimeOn(request.Runtime),
		RuntimeVersionArn: &request.RuntimeVersionArn,
	}

	_, err := wrapper.LambdaClient.PutRuntimeManagementConfig(context.Background(), input)
	if err != nil {
		return err
	}

	return err
}

func HandlePutRuntimeManagementConfig(c *fiber.Ctx) error {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return err
	}

	lambdaClient := lambda.NewFromConfig(cfg)

	var request PutRuntimeManagementConfigRequest
	if err := c.BodyParser(&request); err != nil {
		return err
	}

	wrapper := PutRuntimeManagementConfigWrapper{
		LambdaClient: lambdaClient,
	}

	if err := wrapper.PutRuntimeManagementConfig(request); err != nil {
		return err
	}

	return c.Status(http.StatusOK).JSON("PutRuntimeManagementConfig is successfull")
}
