package apis

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
	"github.com/gofiber/fiber/v2"
	// Other import statements as needed
)

type ListLayersWrapper struct {
	LambdaClient *lambda.Client
}

// ListLayersResponse represents the response payload for listing AWS Lambda layers.
type ListLayersResponse struct {
	Layers []types.LayersListItem `json:"layers"`
}

// ListLayers retrieves a list of all AWS Lambda layers available in your account.
func (wrapper ListLayersWrapper) ListLayers() ([]types.LayersListItem, error) {
	input := &lambda.ListLayersInput{}

	result, err := wrapper.LambdaClient.ListLayers(context.Background(), input)
	if err != nil {
		return nil, err
	}

	response := &ListLayersResponse{
		Layers: result.Layers,
	}

	return response.Layers, nil
}

func HandleListLayers(c *fiber.Ctx) error {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}
	lambdaClient := lambda.NewFromConfig(cfg)

	// Instantiate FunctionWrapper with the Lambda client
	wrapper := ListLayersWrapper{
		LambdaClient: lambdaClient,
	}

	response, err := wrapper.ListLayers()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(response)
}
