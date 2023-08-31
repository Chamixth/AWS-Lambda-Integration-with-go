package apis

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/gofiber/fiber/v2"
)

// State represents the state of a Lambda function.
type State string

const (
	// StateActive indicates that the Lambda function is active and ready to be invoked.
	StateActive State = "Active"
	// StateInactive indicates that the Lambda function is inactive and cannot be invoked.
	StateInactive State = "Inactive"
)

// FunctionWrapper encapsulates function actions used in the examples.
// It contains an AWS Lambda service client that is used to perform user actions.
type FunctionWrapper struct {
	LambdaClient *lambda.Client
}

// CreateFunctionRequest represents the request body for creating a Lambda function.
type CreateFunctionRequest struct {
	FunctionName string `json:"functionName"`
	HandlerName  string `json:"handlerName"`
	IAMRoleARN   string `json:"iamRoleArn"`
	S3Bucket     string `json:"s3Bucket"`     // Use just the bucket name
	S3ObjectKey  string `json:"s3ObjectKey"`  // Add this field for the object key
}

// CreateFunctionResponse represents the response for creating a Lambda function.
type CreateFunctionResponse struct {
	State State `json:"state"`
}

// CreateFunction creates a new Lambda function from code contained in the zipPackage
// buffer. The specified handlerName must match the name of the file and function
// contained in the uploaded code. The role specified by iamRoleArn is assumed by
// Lambda and grants specific permissions.
// When the function already exists, types.StateActive is returned.
// When the function is created, a lambda.FunctionActiveV2Waiter is used to wait until the
// function is active.
func (wrapper FunctionWrapper) CreateFunction(c *fiber.Ctx) error {
	var req CreateFunctionRequest
	if err := c.BodyParser(&req); err != nil {
		c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
		return err
	}

	iamRoleArn := aws.String(req.IAMRoleARN)
	
	var state State
	_, err := wrapper.LambdaClient.CreateFunction(context.TODO(), &lambda.CreateFunctionInput{
		Code: &types.FunctionCode{
			S3Bucket: aws.String(req.S3Bucket),
			S3Key:    aws.String(req.S3ObjectKey),
		},
		FunctionName: aws.String(req.FunctionName),
		Role:         iamRoleArn,
		Handler:      aws.String(req.HandlerName),
		Publish:      true,
		Runtime:    types.RuntimeProvidedal2  ,
	})
	if err != nil {
		var resConflict *types.ResourceConflictException
		if errors.As(err, &resConflict) {
			log.Printf("Function %v already exists.\n", req.FunctionName)
			state = StateActive
		} else {
			c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Couldn't create function"})
			return err
		}
	} else {
		waiter := lambda.NewFunctionActiveV2Waiter(wrapper.LambdaClient)
		funcOutput, err := waiter.WaitForOutput(context.TODO(), &lambda.GetFunctionInput{
			FunctionName: aws.String(req.FunctionName),
		}, 1*time.Minute)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Couldn't wait for function to be active"})

		} else {
			state = State(funcOutput.Configuration.State)
		}
	}

	return c.Status(http.StatusOK).JSON(CreateFunctionResponse{State: state})
}
