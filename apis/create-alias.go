package apis

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/gofiber/fiber/v2"
)

type CreateAlias struct {
	FunctionName    string `json:functionName`
	AliasName       string `json:aliasName`
	FunctionVersion string `json:"functionVersion"`
}

type CreateAliasWrapper struct {
	LambdaClient *lambda.Client
}

func (wrapper CreateAliasWrapper) CreateAlias(request CreateAlias) (*string, error) {
	input := &lambda.CreateAliasInput{
		FunctionName:    aws.String(request.FunctionName),
		Name:            aws.String(request.AliasName),
		FunctionVersion: aws.String(request.FunctionVersion),
	}

	aliasOutput, err := wrapper.LambdaClient.CreateAlias(context.Background(), input)
	if err != nil {
		return nil, err
	}

	return aliasOutput.AliasArn, nil

}
func HandleCreateAliasFunction(c *fiber.Ctx) error {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatalf("Failed to load AWS configuration: %v", err)
	}
	lambdaClient := lambda.NewFromConfig(cfg)

	wrapper := CreateAliasWrapper{
		LambdaClient: lambdaClient,
	}

	var request CreateAlias

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid JSON payload",
		})
	}
	aliasArn, err := wrapper.CreateAlias(request)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"message":  "Alias created successfully",
		"aliasArn": aliasArn,
	})
}
