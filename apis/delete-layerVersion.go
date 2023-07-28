package apis

import (
	"context"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/gofiber/fiber/v2"
)

type DeleteLayerVersionWrapper struct {
	LambdaClient *lambda.Client
}

func (wrapper DeleteLayerVersionWrapper) DeleteLayerVersion(layerName string, versionNumber int) error {
	input := &lambda.DeleteLayerVersionInput{
		LayerName:     &layerName,
		VersionNumber: int64(versionNumber),
	}
	_, err := wrapper.LambdaClient.DeleteLayerVersion(context.Background(), input)

	if err != nil {
		return err
	}

	return err
}

func HandleDeleteLayerVersion(c *fiber.Ctx) error {
	cfg, err := config.LoadDefaultConfig(context.TODO())

	if err != nil {
		return err
	}

	lambdaClient := lambda.NewFromConfig(cfg)

	layerName := c.Query("layerName")
	versionNumber := c.QueryInt("versionName")

	wrapper := DeleteLayerVersionWrapper{
		LambdaClient: lambdaClient,
	}

	if err := wrapper.DeleteLayerVersion(layerName, versionNumber); err != nil {
		return err
	}

	return c.Status(http.StatusOK).JSON("Layer Version Deleted Successfully")
}
