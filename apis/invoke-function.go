package apis

import (
	"context"
	"encoding/json"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/gofiber/fiber/v2"
)

// LambdaRequest represents the structure of the incoming JSON payload
type LambdaRequest struct {
	FunctionName string                 `json:"functionName"`
	Params       map[string]interface{} `json:"params"`
	GetLog       bool                   `json:"getLog"`
}

// FunctionWrapper encapsulates function actions used in the examples.
// It contains an AWS Lambda service client that is used to perform user actions.
type InvokeFunctionWrapper struct {
	LambdaClient *lambda.Client
}

// Invoke invokes the Lambda function specified by functionName, passing the parameters
// as a JSON payload. When getLog is true, types.LogTypeTail is specified, which tells
// Lambda to include the last few log lines in the returned result.
func (wrapper InvokeFunctionWrapper) Invoke(request LambdaRequest) *lambda.InvokeOutput {
	logType := types.LogTypeNone
	if request.GetLog {
		logType = types.LogTypeTail
	}

	payload, err := json.Marshal(request.Params)
	if err != nil {
		log.Panicf("Couldn't marshal parameters to JSON. Here's why %v\n", err)
	}

	invokeOutput, err := wrapper.LambdaClient.Invoke(context.TODO(), &lambda.InvokeInput{
		FunctionName: aws.String(request.FunctionName),
		LogType:      logType,
		Payload:      payload,
	})
	if err != nil {
		log.Panicf("Couldn't invoke function %v. Here's why: %v\n", request.FunctionName, err)
	}

	return invokeOutput
}

func HandleInvokeFunction(c *fiber.Ctx) error {
	// Load the AWS configuration from environment variables or AWS configuration files
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatalf("Failed to load AWS configuration: %v", err)
	}

	// Create a new Lambda service client
	lambdaClient := lambda.NewFromConfig(cfg)

	// Initialize the FunctionWrapper with the Lambda service client
	wrapper := InvokeFunctionWrapper{
		LambdaClient: lambdaClient,
	}
	var request LambdaRequest
		if err := c.BodyParser(&request); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid JSON payload",
			})
		}

		// Invoke the Lambda function
		invokeOutput := wrapper.Invoke(request)

		// Parse the Lambda function response and send it as JSON to the client
		var result map[string]interface{}
		if err := json.Unmarshal(invokeOutput.Payload, &result); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to parse Lambda function response",
			})
		}

		return c.JSON(result)
}
