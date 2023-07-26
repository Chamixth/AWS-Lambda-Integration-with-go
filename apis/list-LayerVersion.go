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

type ListLayerVersionWrapper struct {
	LambdaClient *lambda.Client
}

func (wrapper ListLayerVersionWrapper) ListLayerVersion(layerName string) ([]types.LayerVersionsListItem, error) {
	input := &lambda.ListLayerVersionsInput{
		LayerName: aws.String(layerName),
	}

	result, err := wrapper.LambdaClient.ListLayerVersions(context.Background(), input)

	if err != nil {
		return nil, err
	}

	return result.LayerVersions, err

}

func HandleListLayerVersion(c *fiber.Ctx) error {
	cfg, err := config.LoadDefaultConfig(context.TODO())

	if err != nil {
		return nil
	}

	lambdaClient := lambda.NewFromConfig(cfg)

	wrapper := ListLayerVersionWrapper{
		LambdaClient: lambdaClient,
	}

	layerName := c.Query("layerName")

	result, err := wrapper.ListLayerVersion(layerName)

	if err != nil {
		return err
	}

	return c.Status(http.StatusOK).JSON(result)
}
