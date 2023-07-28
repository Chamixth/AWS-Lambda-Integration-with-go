package apis

import (
	"context"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/gofiber/fiber/v2"
)

type GetPolicyWrapper struct {
	LambdaClient *lambda.Client
}

// GetPolicyResponse represents the JSON response for retrieving the resource policy of a function.
type GetPolicyResponse struct {
	FunctionName string `json:"functionName"`
	Policy       string `json:"policy"`
}

func (wrapper GetPolicyWrapper) GetPolicy(functionName, aliasName string) (*GetPolicyResponse, error) {
	input := &lambda.GetPolicyInput{
		FunctionName: &functionName,
		
	}
	result, err := wrapper.LambdaClient.GetPolicy(context.Background(), input)

	if err != nil {
		return nil, err
	}

	return &GetPolicyResponse{
		FunctionName: functionName,
		Policy:       *result.Policy,
	}, nil
}

func HandleGetPolicy(c *fiber.Ctx) error {
	cfg, err := config.LoadDefaultConfig(context.TODO())

	if err != nil {
		return err
	}

	lambdaClient := lambda.NewFromConfig(cfg)

	wrapper := GetPolicyWrapper{
		LambdaClient: lambdaClient,
	}
	functionName := c.Query("functionName")
	aliasName := c.Query("aliasName")

	result, err := wrapper.GetPolicy(functionName, aliasName)

	if err != nil {
		return err
	}

	return c.Status(http.StatusOK).JSON(result)
}
