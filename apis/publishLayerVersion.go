package apis

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
	"github.com/gofiber/fiber/v2"
)

type PublishLayerVersionWrapper struct {
	LambdaClient *lambda.Client
}

type PublishLayerVersionRequest struct {
	LayerName               string                          `json:"layerName"`
	Content                 string `json:"content"`
	Description             string                          `json:"description"`
	LicenseInfo             string                          `json:"licenseInfo"`
}

type PublishLayerVersionResponse struct {
	VersionNumber int64  `json:"versionNumber"`
	LayerArn      string `json:"layerArn"`
}

func HandlePublishLayerVersion(c *fiber.Ctx) error {
	cfg, err := config.LoadDefaultConfig(context.Background())

	if err != nil {
		return err
	}

	lambdaClient := lambda.NewFromConfig(cfg)

	var request PublishLayerVersionRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	zipContent, err := ioutil.ReadFile(request.Content)
	if err != nil {
		c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to load Lambda function code package"})
		return err
	}

	zipPackage := bytes.NewBuffer(zipContent)

	wrapper := PublishLayerVersionWrapper{
		LambdaClient: lambdaClient,
	}

	input := &lambda.PublishLayerVersionInput{
		LayerName:               &request.LayerName,
		Content:                 &types.LayerVersionContentInput{ZipFile: zipPackage.Bytes()},
		Description:             &request.Description,
		LicenseInfo:             &request.LicenseInfo,
	}

	result, err := wrapper.LambdaClient.PublishLayerVersion(context.Background(), input)

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Couldn't publish layer version"})
	}

	response := &PublishLayerVersionResponse{
		VersionNumber: *&result.Version,
		LayerArn:      *result.LayerArn,
	}

	return c.Status(http.StatusOK).JSON(response)
}
