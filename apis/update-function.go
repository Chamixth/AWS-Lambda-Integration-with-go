package apis

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/service/lambda/lambdaiface"
	"github.com/gofiber/fiber/v2"
)

type UpdateFunctionWrapper struct {
	LambdaClient lambdaiface.LambdaAPI
}

// State represents the state of the Lambda function after the update.
type UpdateState struct {
	State string `json:"state"`
}

func (wrapper UpdateFunctionWrapper) UpdateFunctionCode(functionName string, zipPackage *bytes.Buffer) UpdateState {
	var state UpdateState
	_, err := wrapper.LambdaClient.UpdateFunctionCodeWithContext(context.Background(), &lambda.UpdateFunctionCodeInput{
		FunctionName: aws.String(functionName),
		ZipFile:      zipPackage.Bytes(),
	})
	if err != nil {
		log.Panicf("Couldn't update code for function %v. Here's why: %v\n", functionName, err)
	} else {
		state.State = "UpdateComplete"
	}
	return state
}

func readZipFile(filePath string) (*bytes.Buffer, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	return bytes.NewBuffer(data), nil
}

func HandleUpdateFunction(c *fiber.Ctx) error {
	functionNameToUpdate := c.Query("functionName") // Replace with your function name
	file, err := c.FormFile("zipFile")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Missing or invalid zipFile field"})
	}

	// Open the file
	src, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to read the uploaded file"})
	}
	defer src.Close()

	// Read the uploaded zip file
	data, err := ioutil.ReadAll(src)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to read the uploaded file"})
	}

	// Create a bytes.Buffer from the file data
	zipPackage := bytes.NewBuffer(data)

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create an AWS Lambda client
	lambdaClient := lambda.New(sess)

	// Create a FunctionWrapper instance
	wrapper := UpdateFunctionWrapper{
		LambdaClient: lambdaClient,
	}

	// Call the UpdateFunctionCode method with the new zip file
	newState := wrapper.UpdateFunctionCode(functionNameToUpdate, zipPackage)

	// Print the updated state of the function
	fmt.Printf("Function %v update state: %v\n", functionNameToUpdate, newState.State)

	return c.Status(fiber.StatusOK).JSON(newState)
}
