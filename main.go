package main

import (
	"context"
	"aws-lambda-integration-with-go/apis"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/gofiber/fiber/v2"
)

func main() {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatalf("failed to load the configurations: %v", err)
	}
	lambdaClient := lambda.NewFromConfig(cfg)

	wrapper := apis.FunctionWrapper{
		LambdaClient: lambdaClient,
	}

	app := fiber.New(fiber.Config{
		BodyLimit: 40 * 1024 * 1024, // Set the body limit to 40MB
	})

	app.Post("/createFunction", wrapper.CreateFunction)
	app.Get("/getFunction", apis.HandleGetFunction)
	app.Delete("/deleteFunction", apis.HandleDeleteFunction)
	app.Put("/updateFunction", apis.HandleUpdateFunction)
	app.Put("/updateConfing",apis.HandleUpdateConfigFunction)
	app.Get("listFunction",apis.HandleListFunction)
	app.Post("/addPermission",apis.HandleAddPermissionFunction)
	app.Post("/invokeFunction",apis.HandleInvokeFunction)
	app.Get("/getAlias",apis.HandleGetAliasFunction)
	app.Get("/listAliases",apis.HandleListAliasesFunction)
	app.Post("/createAlias",apis.HandleCreateAliasFunction)
	app.Delete("/deleteAlias",apis.HandleDeleteAlias)
	app.Put("updateAlias",apis.HandleUpdateAliasFunction)
	if err := app.Listen(":8080"); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}

}
