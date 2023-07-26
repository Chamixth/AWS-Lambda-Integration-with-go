package apis

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/gofiber/fiber/v2"
)

type AddLayerVersionPermissionWrapper struct {
	LambdaClient *lambda.Client
}

// AddLayerVersionPermissionRequest represents the JSON payload for adding permissions to a layer version.
type AddLayerVersionPermissionRequest struct {
	LayerName      *string `json:"layerName"`
	VersionNumber  int64   `json:"versionNumber"`
	StatementID    *string `json:"statementId"`
	Action         *string `json:"action"`
	Principal      *string `json:"principal"`
	OrganizationID *string `json:"organizationId,omitempty"`
	RevisionID     *string `json:"revisionId,omitempty"`
}

// AddLayerVersionPermission adds permissions to a specific version of an AWS Lambda layer.
func (wrapper AddLayerVersionPermissionWrapper) AddLayerVersionPermission(request AddLayerVersionPermissionRequest) error {
	input := &lambda.AddLayerVersionPermissionInput{
		LayerName:     request.LayerName,
		VersionNumber: request.VersionNumber,
		StatementId:   request.StatementID,
		Action:        request.Action,
		Principal:     request.Principal,
	}

	// Optional fields
	if request.OrganizationID != nil {
		input.OrganizationId = request.OrganizationID
	}

	if request.RevisionID != nil {
		input.RevisionId = request.RevisionID
	}

	_, err := wrapper.LambdaClient.AddLayerVersionPermission(context.Background(), input)

	return err
}

func HandleAddLayerVersionPermission(c *fiber.Ctx) error {
	var request AddLayerVersionPermissionRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}
	lambdaClient := lambda.NewFromConfig(cfg)

	// Instantiate FunctionWrapper with the Lambda client
	wrapper := AddLayerVersionPermissionWrapper{
		LambdaClient: lambdaClient,
	}

	if err := wrapper.AddLayerVersionPermission(request); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Layer version permission added successfully",
	})
}
