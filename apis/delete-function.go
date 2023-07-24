package apis

import (
	"context"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/gofiber/fiber/v2"
)

// FunctionWrapper encapsulates function actions used in the examples.
// It contains an AWS Lambda service client that is used to perform user actions.

// DeleteFunction deletes the Lambda function specified by functionName.
// DeleteFunction deletes the Lambda function specified by functionName.
// DeleteFunction deletes the Lambda function specified by functionName.
// DeleteFunction deletes the Lambda function specified by functionName.
type DeleteFunctionWrapper struct {
	LambdaClient *lambda.Lambda
}

func (wrapper DeleteFunctionWrapper) DeleteFunction(functionName string) {
	_, err := wrapper.LambdaClient.DeleteFunctionWithContext(context.TODO(), &lambda.DeleteFunctionInput{
		FunctionName: aws.String(functionName),
	})
	if err != nil {
		log.Panicf("Couldn't delete function %v. Here's why: %v\n", functionName, err)
	}
}

func HandleDeleteFunction(c *fiber.Ctx) error {
	functionName := c.Query("functionName")
	if functionName == "" {
		c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Function name is invalid"})
		return fiber.ErrBadRequest
	}

	// Create a new AWS session
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create an AWS Lambda client
	lambdaClient := lambda.New(sess)
	wrapper := DeleteFunctionWrapper{
		LambdaClient: lambdaClient,
	}

	wrapper.DeleteFunction(functionName)

	return c.Status(http.StatusOK).JSON("Function Deleted")
}
