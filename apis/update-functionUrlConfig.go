package apis

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/gofiber/fiber/v2"
)

type UpdateFunctionUrlConfig struct {
	LambdaClient *lambda.Client
}

type CorsSettings struct {
	AllowOrigins   []string `json:"allowOrigins"`
	AllowHeaders   []string `json:"allowHeaders"`
	AllowMethods   []string `json:"allowMethods"`
	MaxAgeSeconds  int32    `json:"maxAgeSeconds"`
}
// UpdateFunctionUrlConfigRequest represents the JSON payload for updating a function URL configuration.
type UpdateFunctionUrlConfigRequest struct {
	FunctionName string `json:"functionName"`
	AliasName    string `json:"aliasName"`
	AuthType     string `json:"authType"`
	Cors CorsSettings `json:"cors"`
}

func (wrapper UpdateFunctionUrlConfig) UpdateFunctionUrlConfig(request UpdateFunctionUrlConfigRequest) error {
	input := &lambda.UpdateFunctionUrlConfigInput{
		FunctionName: aws.String(request.FunctionName),
		Qualifier:    aws.String(request.AliasName),
		AuthType:     types.FunctionUrlAuthType(*aws.String(request.AuthType)),
		Cors: &types.Cors{},
	}
	_, err := wrapper.LambdaClient.UpdateFunctionUrlConfig(context.Background(), input)

	if err != nil {
		return err
	}
	return err
}

func HandleUpdateFunctionUrlConfig(c *fiber.Ctx) error {
	var request UpdateFunctionUrlConfigRequest
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}
	lambdaClient := lambda.NewFromConfig(cfg)

	// Instantiate FunctionWrapper with the Lambda client
	wrapper := UpdateFunctionUrlConfig{
		LambdaClient: lambdaClient,
	}

	if err := c.BodyParser(&request); err != nil{
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}
	if err := wrapper.UpdateFunctionUrlConfig(request); err != nil{
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"message":"Function Url Config updated Successfully",
	})

}
