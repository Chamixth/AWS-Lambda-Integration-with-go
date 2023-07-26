package apis

import (
	"context"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
	"github.com/gofiber/fiber/v2"
)

type GetLayerVersionByArnwrapper struct {
	LambdaClient *lambda.Client
}

type LayerArn struct {
	Arn string `json:"arn"`
}

func (wrapper GetLayerVersionByArnwrapper) GetLayerVersionByArn(layerArn LayerArn) (*types.LayerVersionContentOutput, error) {
	input := &lambda.GetLayerVersionByArnInput{
		Arn: aws.String(layerArn.Arn),
	}

	result, err := wrapper.LambdaClient.GetLayerVersionByArn(context.Background(), input)

	if err != nil {
		return nil, err
	}

	return result.Content, nil
}

func HandleGetLayerVersionByArn(c *fiber.Ctx) error {
	cfg, err := config.LoadDefaultConfig(context.TODO())

	if err != nil {
		return err
	}

	lambdaClient := lambda.NewFromConfig(cfg)

	wrapper := GetLayerVersionByArnwrapper{
		LambdaClient: lambdaClient,
	}

	var layerArn LayerArn
	if err := c.BodyParser(&layerArn); err != nil {
		return err
	}

	 result, err := wrapper.GetLayerVersionByArn(layerArn)
	
	 if err != nil {
		return err
	}

	return c.Status(http.StatusOK).JSON(result)

}
