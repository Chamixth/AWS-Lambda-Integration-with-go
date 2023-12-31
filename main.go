package main

import (
	"aws-lambda-integration-with-go/apis"
	"context"
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
	app.Put("/updateConfing", apis.HandleUpdateConfigFunction)
	app.Get("listFunction", apis.HandleListFunction)
	app.Post("/addPermission", apis.HandleAddPermissionFunction)
	app.Post("/invokeFunction", apis.HandleInvokeFunction)
	app.Get("/getAlias", apis.HandleGetAliasFunction)
	app.Get("/listAliases", apis.HandleListAliasesFunction)
	app.Post("/createAlias", apis.HandleCreateAliasFunction)
	app.Delete("/deleteAlias", apis.HandleDeleteAlias)
	app.Put("updateAlias", apis.HandleUpdateAliasFunction)
	app.Get("/getAccountSetting", apis.HandleGetAccountSeetings)
	app.Post("/createFunctionUrlConfig", apis.HandleCreateFunctionUrlConfig)
	app.Get("/getFunctionUrlConfig", apis.HandleGetFunctionUrlConfig)
	app.Delete("deleteFunctionUrlConfig", apis.HandleDeleteFunctionUrlConfig)
	app.Put("/updateFunctionUrlConfig", apis.HandleUpdateFunctionUrlConfig)
	app.Get("/listFunctionUrlConfigs", apis.HandleListFunctionUrlConfig)
	app.Post("/addLayerVersionPermission",apis.HandleAddLayerVersionPermission)
	app.Get("/getLayerVersion",apis.HandleGetLayerVersion)
	app.Get("/getLayerVersionByArn",apis.HandleGetLayerVersionByArn)
	app.Get("/listLayers",apis.HandleListLayers)
	app.Get("/listLayerVersions",apis.HandleListLayerVersion)
	app.Get("/listVersionByFunction",apis.HandleListVersionByFunction)
	app.Post("/publishVersion",apis.HandlePublishVersion)
	app.Delete("/removePermission",apis.HandleRemovePermission)
	app.Post("/tagResources",apis.HandleTagResource)
	app.Delete("/unTagResource",apis.HandleUntagResource)
	app.Put("/updateFunctionConfiguration",apis.HandleUpdateFunctionConfiguration)
	app.Post("/updateFunctionEventInvokeConfig",apis.HandleUpdateFunctionEventInvokeConfig)
	app.Get("listTags",apis.HandleListTagsFunction)
	app.Post("/createCodeSigningConfig",apis.HandleCreateSiginingConfig)
	app.Get("/getFunctionCodeSigningConfig",apis.HandleGetFunctionCodeSigningConfig)
	app.Get("/listCodeSigningConfig",apis.HandleListCodeSigningConfig)
	app.Delete("/deleteCodeSigningConfig",apis.HandleDeleteCodeSigningConfig)
	app.Get("/getCodeSigningConfig",apis.HandleGetCodeSigningConfig)
	app.Delete("/deleteFunctionCodeSigningConfig",apis.HandleDeleteFunctionCodeSigningConfig)
	app.Put("/updateCodeSigninConfig",apis.HandleUpdateCodeSigningConfig)
	app.Get("/getPolicy",apis.HandleGetPolicy)
	app.Get("/getFunctionConfiguration",apis.HandleGetFunctionConfiguration)
	app.Get("/getFunctionEventInvokeConfig",apis.HandleGetFunctionEventInvokeConfig)
	app.Get("getFunctionConcurrency",apis.HandleGetFunctionConcurrency)
	app.Delete("/deleteFunctionEventInvokeConfig",apis.HandleDeleteFunctionEventInvokeConfig)
	app.Delete("deleteFunctionConcurrency",apis.HandleDeleteFunctionConcurrency)
	app.Delete("/deleteLayerVersion",apis.HandleDeleteLayerVersion)
	app.Delete("/deleteProvisionedConcurrencyConfig",apis.HandleDeleteProvisionedConcurrencyConfig)
	app.Get("/getLayerVersionPolicy",apis.HandleGetLayerVersionPolicy)
	app.Get("/getProvisionedConcurrencyConfig",apis.HandleGetProvisionedConcurrencyConfig)
	app.Get("/getRuntimeManagementConfig",apis.HandleGetRuntimeManagementConfig)
	app.Get("/listEventSourceMappings",apis.HandleListEventSourceMapping)
	app.Get("/listFunctionEventInvokeConfigs",apis.HandleListFunctionEventInvokeConfigs)
	app.Get("/listFunctionByCodeSigningConfig",apis.HandleListFunctionByCodeSigningConfig)
	app.Get("/listProvisionedConcurrencyConfigs",apis.HandleListProvisionedConcurrencyConfigs)
	app.Post("/publishLayerVersion",apis.HandlePublishLayerVersion)
	app.Put("/putFunctionCodeSigningConfig",apis.HandlePutFunctionCodeSigningConfig)
	app.Put("/putFunctionConcurrency",apis.HandlePutFunctionConcurrency)
	app.Put("putFunctionEventInvokeConfig",apis.HandlePutFunctionEventInvokeConfig)
	app.Put("/putProvisionedConcurrencyConfig",apis.HandlePutProvisionedConcurrencyConfig)
	app.Put("/putRuntimeManagementConfig",apis.HandlePutRuntimeManagementConfig)
	app.Delete("/removeLayerVersionPermission",apis.HandleRemoveLayerVersionPermission)
	if err := app.Listen(":8000"); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}

}
