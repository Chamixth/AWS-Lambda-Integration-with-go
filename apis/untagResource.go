package apis

import (
	"context"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/gofiber/fiber/v2"
)

type UntagResourceWrapper struct {
	LambdaClient *lambda.Client
}

type UntagResourceRequest struct {
	FunctionArn string   `json:"functionArn"`
	TagKey      []string `json:"tagKey"`
}

func (wrapper UntagResourceWrapper) UntagResource(request UntagResourceRequest) error {
	input := &lambda.UntagResourceInput{
		Resource: &request.FunctionArn,
		TagKeys:  request.TagKey,
	}
	_, err := wrapper.LambdaClient.UntagResource(context.Background(), input)

	return err
}

func HandleUntagResource(c *fiber.Ctx) error {
	var request UntagResourceRequest
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil
	}
	lambdaClient := lambda.NewFromConfig(cfg)
	wrapper := UntagResourceWrapper{
		LambdaClient: lambdaClient,
	}
	if err := c.BodyParser(&request); err != nil {
		return err
	}

	if err := wrapper.UntagResource(request); err != nil {
		return err
	}

	return c.Status(http.StatusOK).JSON("Untag Resource is successful")
}
