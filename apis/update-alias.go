package apis

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/gofiber/fiber/v2"
)

type UpdateAlias struct {
	FunctionName    string `json:functionName`
	AliasName       string `json:aliasName`
	FunctionVersion string `json:"functionVersion"`
}

type UpdateAliasWrapper struct {
	LambdaClient *lambda.Client
}

func (wrapper UpdateAliasWrapper) UpdateAlias(request UpdateAlias) error {
	input := &lambda.UpdateAliasInput{
		FunctionName:    aws.String(request.FunctionName),
		Name:            aws.String(request.AliasName),
		FunctionVersion: aws.String(request.FunctionVersion),
	}
	_, err := wrapper.LambdaClient.UpdateAlias(context.Background(), input)

	return err
}

func HandleUpdateAliasFunction(c *fiber.Ctx) error {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}
	lambdaClient := lambda.NewFromConfig(cfg)

	// Instantiate FunctionWrapper with the Lambda client
	wrapper := UpdateAliasWrapper{
		LambdaClient: lambdaClient,
	}

	var request UpdateAlias

	if err := c.BodyParser(&request); err != nil{
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Invalid JSON payload",
		})
	}
	if err := wrapper.UpdateAlias(request); err != nil{
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message":  "Alias updated successfully",
	})
}
