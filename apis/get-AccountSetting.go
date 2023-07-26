package apis

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/gofiber/fiber/v2"
)

// FunctionWrapper encapsulates function actions used in the examples.
// It contains an AWS STS service client that is used to perform user actions.
type GeAccountSettingsFunctionWrapper struct {
	STSClient *sts.Client
}

// GetCallerIdentityInfo retrieves the caller identity information.
func (wrapper GeAccountSettingsFunctionWrapper) GetAccountSetting() (*sts.GetCallerIdentityOutput, error) {
	input := &sts.GetCallerIdentityInput{}

	output, err := wrapper.STSClient.GetCallerIdentity(context.Background(), input)
	if err != nil {
		return nil, err
	}

	return output, nil
}
func HandleGetAccountSeetings(c *fiber.Ctx) error {

	// Load the AWS configuration from environment variables or AWS configuration files
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatalf("Failed to load AWS configuration: %v", err)
	}

	// Create a new STS service client
	stsClient := sts.NewFromConfig(cfg)

	// Initialize the FunctionWrapper with the STS service client
	wrapper := GeAccountSettingsFunctionWrapper{
		STSClient: stsClient,
	}

	output, err := wrapper.GetAccountSetting()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Format and display the AWS account ID
	response := map[string]string{
		"AccountID":  aws.ToString(output.Account),
		"UserID":     aws.ToString(output.UserId),
		"ARN":        aws.ToString(output.Arn),
	
		
	}

	return c.JSON(response)

}
