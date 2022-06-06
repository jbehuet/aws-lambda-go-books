import * as cdk from "@aws-cdk/core";
// import * as assets from "@aws-cdk/aws-s3-assets";
import * as lambda from "@aws-cdk/aws-lambda";
import * as apigw from "@aws-cdk/aws-apigateway";
import * as iam from "@aws-cdk/aws-iam";
import * as path from "path";

export class GoLambdaStack extends cdk.Stack {
  constructor(scope: cdk.App, id: string, props?: cdk.StackProps) {
    super(scope, id, props);

    // const myLambdaAsset = new assets.Asset(
    //   // @ts-ignore - this expects Construct not cdk.Construct
    //   this,
    //   "GoLambdaFnZip",
    //   {
    //     path: path.join(__dirname, "../../go-lambda-app/cmd"),
    //   }
    // );

    const lambdaFn = new lambda.Function(this, "GoLambdaFn", {
      code: lambda.Code.fromAsset(path.join(__dirname, "../../go-lambda-app"), {
        bundling: {
          image: lambda.Runtime.GO_1_X.bundlingImage,
          user: "root",
          command: [
            "bash",
            "-c",
            [
              "go test -v cmd/main.go",
              "GOOS=linux go build -o /asset-output/main cmd/main.go",
            ].join(" && "),
          ],
        },
      }),
      timeout: cdk.Duration.seconds(300),
      runtime: lambda.Runtime.GO_1_X,
      handler: "main",
    });

    // API Gateway
    const api = new apigw.RestApi(this, "booksApi", {
      description: "API for books",
      deployOptions: {
        stageName: "dev",
      },

      // // enable CORS
      // defaultCorsPreflightOptions: {
      //   allowHeaders: [
      //     'Content-Type',
      //     'X-Amz-Date',
      //     'Authorization',
      //     'X-Api-Key',
      //   ],
      //   allowMethods: ['OPTIONS', 'GET', 'POST', 'PUT', 'PATCH', 'DELETE'],
      //   allowCredentials: true,
      //   allowOrigins: ['http://localhost:3000'],
      // },
    });

    // create an Output for the API URL
    new cdk.CfnOutput(this, "apiUrl", { value: api.url });

    // add a /books resource
    const books = api.root.addResource("books");

    // integrate GET /books with lambdaFn
    books.addMethod(
      "GET",
      new apigw.LambdaIntegration(lambdaFn, { proxy: true })
    );

    lambdaFn.addToRolePolicy(
      new iam.PolicyStatement({
        actions: ["dynamodb:Scan"],
        resources: ["arn:aws:dynamodb:eu-west-3:463496343972:table/Books"],
      })
    );

    // new apigw.LambdaRestApi(
    //   // @ts-ignore - this expects Construct not cdk.Construct
    //   this,
    //   "GoLambdaFnEndpoint",
    //   {
    //     handler: lambdaFn,
    //   }
    // );
  }
}
