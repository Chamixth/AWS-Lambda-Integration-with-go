package apis

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/gofiber/fiber/v2"
	// Other import statements as needed
)

type UpdateFunctionConfigurationWrapper struct {
	LambdaClient *lambda.Client
}

// UpdateFunctionConfigurationRequest represents the JSON payload for updating a function's configuration.
type UpdateFunctionConfigurationRequest struct {
	FunctionName string `json:"functionName"`
	Description  string `json:"description"`
	Timeout      int32  `json:"timeout"`
	MemorySize   int32  `json:"memorySize"`
	// Add other fields as needed to update the function's configuration.
}

// UpdateFunctionConfiguration updates the configuration of an AWS Lambda function.
func (wrapper UpdateFunctionConfigurationWrapper) UpdateFunctionConfiguration(request UpdateFunctionConfigurationRequest) error {
	input := &lambda.UpdateFunctionConfigurationInput{
		FunctionName: &request.FunctionName,
		Description:  &request.Description,
		Timeout:      &request.Timeout,
		MemorySize:   &request.MemorySize,
		// Add other fields as needed to update the function's configuration.
	}

	_, err := wrapper.LambdaClient.UpdateFunctionConfiguration(context.Background(), input)
	return err
}

func HandleUpdateFunctionConfiguration(c *fiber.Ctx) error {
	var request UpdateFunctionConfigurationRequest
	if err := c.BodyParser(&request); err != nil {
		return err
	}

	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return err
	}
	lambdaClient := lambda.NewFromConfig(cfg)

	// Instantiate FunctionWrapper with the Lambda client
	wrapper := UpdateFunctionConfigurationWrapper{
		LambdaClient: lambdaClient,
	}

	err = wrapper.UpdateFunctionConfiguration(request)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Function configuration updated successfully",
	})
}
