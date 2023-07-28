package apis

import (
	"context"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/gofiber/fiber/v2"
)

type RemoveLayerVersionPermissionWrapper struct {
	LambdaClient *lambda.Client
}

type RemoveLayerVersionPermissionRequest struct {
	LayerName       string `json:"layerName"`
	VersionNumber   *int64  `json:"versionNumber"`
	StatementId     string `json:"statementId"`
	
}

func (wrapper RemoveLayerVersionPermissionWrapper) RemoveLayerVersionPermission(request RemoveLayerVersionPermissionRequest) error {
	input := &lambda.RemoveLayerVersionPermissionInput{
		LayerName:     &request.LayerName,
		VersionNumber: *request.VersionNumber,
		StatementId:   &request.StatementId,
		
	}

	_, err := wrapper.LambdaClient.RemoveLayerVersionPermission(context.Background(), input)
	return err
}

func HandleRemoveLayerVersionPermission(c *fiber.Ctx) error {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return err
	}

	lambdaClient := lambda.NewFromConfig(cfg)

	var request RemoveLayerVersionPermissionRequest
	if err := c.BodyParser(&request); err != nil {
		return err
	}

	wrapper := RemoveLayerVersionPermissionWrapper{
		LambdaClient: lambdaClient,
	}

	if err := wrapper.RemoveLayerVersionPermission(request); err != nil {
		return err
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{"message": "Layer version permission removed successfully"})
}
