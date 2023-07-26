package apis

import (
	"context"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
	"github.com/gofiber/fiber/v2"
)

type ListFunctionUrlConfigWrapper struct {
	LambdaClient *lambda.Client
}

// ListFunctionUrlConfigResponse represents the response payload for listing function URL configurations.
type ListFunctionUrlConfigResponse struct {
	FunctionName string       `json:"functionName"`
	AliasName    string       `json:"aliasName"`
	AuthType     string       `json:"authType"`
	Cors         CorsSettings `json:"cors"`
}

func (wrapper ListFunctionUrlConfigWrapper) ListFunctionUrlConfig(functionName string, maxItem int) ([]types.FunctionUrlConfig, error) {
	var results []types.FunctionUrlConfig
	paginator := lambda.NewListFunctionUrlConfigsPaginator(wrapper.LambdaClient, &lambda.ListFunctionUrlConfigsInput{
		FunctionName: aws.String(functionName),
		MaxItems:     aws.Int32(int32(maxItem)),
	})
	for paginator.HasMorePages() && len(results) < maxItem {
		pageOutput, err := paginator.NextPage(context.Background())
		if err != nil {
			return nil,err
		}
		results = append(results, pageOutput.FunctionUrlConfigs...)
	}
	return results, nil

}

func HandleListFunctionUrlConfig(c *fiber.Ctx) error {
	cfg, err := config.LoadDefaultConfig(context.TODO())

	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}

	lambdaClient := lambda.NewFromConfig(cfg)

	wrapper := ListFunctionUrlConfigWrapper{
		LambdaClient: lambdaClient,
	}

	functionName := c.Query("functionName")
	maxItem := 20

	output, err := wrapper.ListFunctionUrlConfig(functionName, maxItem)

	if err != nil {
		return err
	}

	return c.Status(http.StatusOK).JSON(output)

}
