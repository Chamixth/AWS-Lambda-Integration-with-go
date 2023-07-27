package apis

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/gofiber/fiber/v2"
)

type ListTagsWrapper struct {
	LambdaClient *lambda.Client
}

// ListTagsResponse represents the JSON response for listing the tags of a function.
type ListTagsResponse struct {
	FunctionName string            `json:"functionName"`
	Tags         map[string]string `json:"tags"`
}

// ListTags lists the tags associated with an AWS Lambda function.
func (wrapper ListTagsWrapper) ListTags(functionName string) (*ListTagsResponse, error) {
	input := &lambda.ListTagsInput{
		Resource: &functionName,
	}

	output, err := wrapper.LambdaClient.ListTags(context.Background(), input)
	if err != nil {
		return nil, err
	}

	return &ListTagsResponse{
		FunctionName: functionName,
		Tags:         output.Tags,
	}, nil
}

func HandleListTagsFunction(c *fiber.Ctx) error {
	functionName := c.Query("functionArn")

	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}
	lambdaClient := lambda.NewFromConfig(cfg)

	// Instantiate FunctionWrapper with the Lambda client
	wrapper := ListTagsWrapper{
		LambdaClient: lambdaClient,
	}

	response, err := wrapper.ListTags(functionName)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(response)
}
