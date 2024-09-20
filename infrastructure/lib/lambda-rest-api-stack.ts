import * as cdk from "aws-cdk-lib";
import { Construct } from "constructs";
import { Function, Runtime, Code } from "aws-cdk-lib/aws-lambda";
import {
  CorsHttpMethod,
  HttpApi,
  HttpMethod,
} from "aws-cdk-lib/aws-apigatewayv2";
import { HttpLambdaIntegration } from "aws-cdk-lib/aws-apigatewayv2-integrations";

export class HttpApiStack extends cdk.Stack {
  private static readonly ROUTES = [
    "/ping",
    "/groceries",
    "/groceries/magic",
    "/groceries/batchDelete",
    "/groceries/{id+}",
    "/users",
    "/users/{id+}",
    "/households",
    "/households/join/{householdId}/{userId+}",
    "/households/leave/{householdId}/{userId+}",
    "/catalog",
    "/receipt/upload",
  ];

  public readonly lambdaFunction: Function;

  constructor(scope: Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    this.lambdaFunction = new Function(this, "GoFunction", {
      runtime: Runtime.PROVIDED_AL2023,
      handler: "bootstrap",
      code: Code.fromAsset("../api/dist"),
      timeout: cdk.Duration.seconds(60),
      memorySize: 512,
    });

    const lambdaIntegration = new HttpLambdaIntegration(
      "GoFunctionIntegration",
      this.lambdaFunction,
      {}
    );

    const httpApi = new HttpApi(this, "GoApi", {
      corsPreflight: {
        allowOrigins: ["*"],
        allowMethods: [CorsHttpMethod.ANY],
        allowHeaders: ["*"],
      },
    });

    HttpApiStack.ROUTES.map((path) => {
      httpApi.addRoutes({
        methods: [HttpMethod.ANY],
        path,
        integration: lambdaIntegration,
      });
    });

    // Output the API endpoint URL
    new cdk.CfnOutput(this, "ApiEndpoint", {
      value: httpApi.apiEndpoint,
    });
  }
}
