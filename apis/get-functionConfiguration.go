package apis

import (
	"context"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/gofiber/fiber/v2"
)

type GetFunctionConfigurationWrapper struct {
	LambdaClient *lambda.Client
}

// GetFunctionConfigurationResponse represents the JSON response for retrieving the configuration of a function.
type GetFunctionConfigurationResponse struct {
	FunctionName string `json:"functionName"`
	FunctionArn  string `json:"functionArn"`
	Role         string `json:"role"`
	Runtime      string `json:"runtime"`
	Handler      string `json:"handler"`
	Description  string `json:"description"`
	MemorySize   int32  `json:"memorySize"`
	Timeout      int32  `json:"timeout"`
	LastModified string `json:"lastModified"`
	CodeSize     int64  `json:"codeSize"`
}

func (wrapper GetFunctionConfigurationWrapper) GetFunctionConfiguration(functionName string) (*GetFunctionConfigurationResponse, error) {
	input := &lambda.GetFunctionConfigurationInput{
		FunctionName: &functionName,
	}

	output, err := wrapper.LambdaClient.GetFunctionConfiguration(context.Background(), input)

	if err != nil {
		return nil, err
	}

	return &GetFunctionConfigurationResponse{
		FunctionName: *output.FunctionName,
		FunctionArn:  *output.FunctionArn,
		Role:         *output.Role,
		Runtime:      string(output.Runtime),
		Handler:      *output.Handler,
		Description:  *output.Description,
		MemorySize:   *output.MemorySize,
		Timeout:      *output.Timeout,
		LastModified: *output.LastModified,
		CodeSize:     *&output.CodeSize,
	}, nil
}

func HandleGetFunctionConfiguration(c *fiber.Ctx) error {
	cfg, err := config.LoadDefaultConfig(context.TODO())

	if err != nil {
		return err
	}

	lambdaClient := lambda.NewFromConfig(cfg)

	wrapper := GetFunctionConfigurationWrapper{
		LambdaClient: lambdaClient,
	}

	functionName := c.Query("functionName")

	result, err := wrapper.GetFunctionConfiguration(functionName)

	if err != nil {
		return err
	}

	return c.Status(http.StatusOK).JSON(result)
}
