package apis

import (
	"context"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/gofiber/fiber/v2"
)

type GetLayerVersionPolicyWrapper struct {
	LambdaClient *lambda.Client
}

// GetLayerVersionPolicyResponse represents the JSON response for retrieving the policy of a Lambda layer version.
type GetLayerVersionPolicyResponse struct {
	LayerArn string `json:"layerArn"`
	Policy   string `json:"policy"`
}

// GetLayerVersionPolicy retrieves the policy of a Lambda layer version.
func (wrapper GetLayerVersionPolicyWrapper) GetLayerVersionPolicy(layerArn string, versionNumber int) (*GetLayerVersionPolicyResponse, error) {
	input := &lambda.GetLayerVersionPolicyInput{
		LayerName:     &layerArn,
		VersionNumber: int64(versionNumber),
	}

	output, err := wrapper.LambdaClient.GetLayerVersionPolicy(context.Background(), input)
	if err != nil {
		return nil, err
	}

	return &GetLayerVersionPolicyResponse{
		LayerArn: layerArn,
		Policy:   *output.Policy,
	}, nil
}

func HandleGetLayerVersionPolicy(c *fiber.Ctx) error {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return err
	}

	lambdaClient := lambda.NewFromConfig(cfg)

	wrapper := GetLayerVersionPolicyWrapper{
		LambdaClient: lambdaClient,
	}

	layerArn := c.Query("layerArn")
	versionNumber := c.QueryInt("versionNumber")

	

	result, err := wrapper.GetLayerVersionPolicy(layerArn, versionNumber)
	if err != nil {
		return err
	}

	return c.Status(http.StatusOK).JSON(result)
}


