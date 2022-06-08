import * as lambda from "@aws-cdk/aws-lambda-go-alpha";
import { Construct } from "constructs";
import { aws_iam, CfnOutput, Stack, StackProps } from "aws-cdk-lib";

import { LambdaIntegration, RestApi } from "aws-cdk-lib/aws-apigateway";
import * as path from "path";

export class GoLambdaBooksStack extends Stack {
  constructor(scope: Construct, id: string, props?: StackProps) {
    super(scope, id, props);

    const booksLambdaFn = new lambda.GoFunction(this, "go-lambda-books-fn", {
      entry: path.join(__dirname, "../../go-lambda-app/cmd"),
    });

    booksLambdaFn.addToRolePolicy(
      new aws_iam.PolicyStatement({
        actions: [
          "dynamodb:Scan",
          "dynamodb:GetItem",
          "dynamodb:PutItem",
          "dynamodb:UpdateItem",
          "dynamodb:DeleteItem",
        ],
        resources: ["arn:aws:dynamodb:eu-west-3:463496343972:table/books"],
      })
    );

    const api = new RestApi(this, "books-api", {
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
    new CfnOutput(this, "apiUrl", { value: api.url });

    // add a /books resource
    const booksApi = api.root.addResource("books");

    // integrate GET /books with lambdaFn
    booksApi.addMethod(
      "GET",
      new LambdaIntegration(booksLambdaFn, { proxy: true })
    );

    booksApi.addMethod(
      "POST",
      new LambdaIntegration(booksLambdaFn, { proxy: true })
    );

    booksApi.addMethod(
      "PUT",
      new LambdaIntegration(booksLambdaFn, { proxy: true })
    );

    booksApi.addMethod(
      "DELETE",
      new LambdaIntegration(booksLambdaFn, { proxy: true })
    );
  }
}
