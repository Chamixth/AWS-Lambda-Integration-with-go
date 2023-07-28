package apis

import (
	"context"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/gofiber/fiber/v2"
)

type ListEventSourceMappingWrapper struct {
	LambdaClient *lambda.Client
}

// ListEventSourceMappingResponse represents the JSON response for listing event source mappings.
type ListEventSourceMappingResponse struct {
	EventSourceMappings []EventSourceMappingDetails `json:"eventSourceMappings"`
}

// EventSourceMappingDetails represents the details of an event source mapping.
type EventSourceMappingDetails struct {
	UUID               string `json:"uuid"`
	FunctionName       string `json:"functionName"`
	EventSourceArn     string `json:"eventSourceArn"`
	LastModified       string `json:"lastModified"`
	LastProcessingTime string `json:"lastProcessingTime"`
	State              string `json:"state"`
	StateTransition    string `json:"stateTransition"`
}

func (wrapper ListEventSourceMappingWrapper) ListEventSourceMapping() (*ListEventSourceMappingResponse, error) {
	input := &lambda.ListEventSourceMappingsInput{}

	result, err := wrapper.LambdaClient.ListEventSourceMappings(context.Background(), input)
	if err != nil {
		return nil, err
	}

	var mappings []EventSourceMappingDetails

	for _, mapping := range result.EventSourceMappings {
		mappings = append(mappings, EventSourceMappingDetails{
			UUID:               *mapping.UUID,
			FunctionName:       *mapping.FunctionArn,
			EventSourceArn:     *mapping.EventSourceArn,
			LastModified:       mapping.LastModified.String(),
			LastProcessingTime: *mapping.LastProcessingResult,
			State:              *mapping.State,
			StateTransition:    *mapping.StateTransitionReason,
		})
	}

	response := &ListEventSourceMappingResponse{
		EventSourceMappings: mappings,
	}

	return response, nil
}

func HandleListEventSourceMapping(c *fiber.Ctx) error {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return err
	}

	lambdaClient := lambda.NewFromConfig(cfg)

	wrapper := ListEventSourceMappingWrapper{
		LambdaClient: lambdaClient,
	}

	result, err := wrapper.ListEventSourceMapping()
	if err != nil {
		return err
	}

	return c.Status(http.StatusOK).JSON(result)
}
