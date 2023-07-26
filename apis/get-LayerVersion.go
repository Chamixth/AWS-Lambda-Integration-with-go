package apis

import (
	"context"
	"log"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
	"github.com/gofiber/fiber/v2"
)

type GetLayerVersionWrapper struct {
	LambdaClient *lambda.Client
}

// GetLayerVersionResponse represents the response payload for getting information about a layer version.
type GetLayerVersionResponse struct {
	VersionDetails types.LayerVersionsListItem `json:"versionDetails"`
}

// GetLayerVersion retrieves information about a specific version of an AWS Lambda layer.
func (wrapper GetLayerVersionWrapper) GetLayerVersion(layerName string, versionNumber int64) (*types.LayerVersionContentOutput, error) {
	input := &lambda.GetLayerVersionInput{
		LayerName:     &layerName,
		VersionNumber: versionNumber,
	}

	result, err := wrapper.LambdaClient.GetLayerVersion(context.Background(), input)
	if err != nil {
		return nil, err
	}

	return result.Content, nil
}

func HandleGetLayerVersion(c *fiber.Ctx) error {
	layerName := c.Query("layerName")
	versionNumberStr := c.Query("versionNumber")

	// Convert the versionNumber parameter to an int64
	versionNumber, err := strconv.ParseInt(versionNumberStr, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid version number"})
	}

	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}
	lambdaClient := lambda.NewFromConfig(cfg)

	// Instantiate FunctionWrapper with the Lambda client
	wrapper := GetLayerVersionWrapper{
		LambdaClient: lambdaClient,
	}

	response, err := wrapper.GetLayerVersion(layerName, versionNumber)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(response)
}
