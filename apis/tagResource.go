package apis

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/gofiber/fiber/v2"
)

type TagResourcesWrapper struct {
	LambdaClient *lambda.Client
}

type TagResourceRequest struct {
	ResourceARN string            `json:"resourceARN"`
	Tags        map[string]string `json:"tags"`
}

// Helper function to convert map[string]string to []types.Tag

func (wrapper TagResourcesWrapper) TagResource(resourceArn string, tags map[string]string) error {

	input := &lambda.TagResourceInput{
		Resource: &resourceArn,
		Tags:     tags,
	}

	_, err := wrapper.LambdaClient.TagResource(context.Background(), input)

	return err
}

func HandleTagResource(c *fiber.Ctx) error {
	var request TagResourceRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return err
	}
	lambdaClient := lambda.NewFromConfig(cfg)

	wrapper := TagResourcesWrapper{
		LambdaClient: lambdaClient,
	}

	err = wrapper.TagResource(request.ResourceARN, request.Tags)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"message": "Tags added to resource successfully",
	})
}
