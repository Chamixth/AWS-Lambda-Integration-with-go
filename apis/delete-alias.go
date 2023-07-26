package apis

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/gofiber/fiber/v2"
)

type DeleteAliasWrapper struct {
	LambdaCLient *lambda.Client
}

func (wrapper DeleteAliasWrapper) DeleteAlias(functionName, aliasName string) error {
	input := &lambda.DeleteAliasInput{
		FunctionName: aws.String(functionName),
		Name:         aws.String(aliasName),
	}

	_, err := wrapper.LambdaCLient.DeleteAlias(context.Background(), input)

	return err
}
func HandleDeleteAlias(c *fiber.Ctx) error {
	cfg, err := config.LoadDefaultConfig(context.Background())

	if err != nil {
		log.Fatalf("Failed to load AWS configuration: %v", err)
	}
	// Create a new Lambda service client
	lambdaClient := lambda.NewFromConfig(cfg)

	wrapper := DeleteAliasWrapper{
		LambdaCLient: lambdaClient,
	}

	functionName := c.Query("functionName")
	aliasName := c.Query("aliasName")

	err_ := wrapper.DeleteAlias(functionName,aliasName)

	if err_ != nil {
		return nil
	}

	return c.JSON(fiber.Map{
		"message": "Alias deleted successfully",
	})

}
