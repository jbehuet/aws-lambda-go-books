import * as lambda from "@aws-cdk/aws-lambda-go-alpha";
import { Construct } from "constructs";
import { aws_iam, CfnOutput, Stack, StackProps } from "aws-cdk-lib";

import { LambdaIntegration, RestApi } from "aws-cdk-lib/aws-apigateway";
import * as path from "path";

export class GoLambdaBooksStack extends Stack {
  constructor(scope: Construct, id: string, props?: StackProps) {
    super(scope, id, props);

    const booksLambdaFn = new lambda.GoFunction(this, "go-lambda-books-fn", {
      entry: path.join(__dirname, "../../go-lambda-app/cmd/books"),
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

    const bookCoversLambdaFn = new lambda.GoFunction(
      this,
      "go-lambda-book-covers-fn",
      {
        entry: path.join(__dirname, "../../go-lambda-app/cmd/covers"),
      }
    );

    bookCoversLambdaFn.addToRolePolicy(
      new aws_iam.PolicyStatement({
        effect: aws_iam.Effect.ALLOW,
        actions: ["s3:ListAllMyBuckets"],
        resources: ["*"],
      })
    );

    bookCoversLambdaFn.addToRolePolicy(
      new aws_iam.PolicyStatement({
        effect: aws_iam.Effect.ALLOW,
        actions: ["s3:ListBucket", "s3:GetBucketLocation"],
        resources: ["arn:aws:s3:::jbehuet-book-covers"],
      })
    );

    bookCoversLambdaFn.addToRolePolicy(
      new aws_iam.PolicyStatement({
        effect: aws_iam.Effect.ALLOW,
        actions: [
          "s3:PutObject",
          "s3:PutObjectAcl",
          "s3:GetObject",
          "s3:GetObjectAcl",
          "s3:DeleteObject",
        ],
        resources: [
          "arn:aws:s3:::jbehuet-book-covers",
          "arn:aws:s3:::jbehuet-book-covers/*",
        ],
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
    const booksRessource = api.root.addResource("books");

    // integrate GET /books with lambdaFn
    booksRessource.addMethod(
      "GET",
      new LambdaIntegration(booksLambdaFn, { proxy: true })
    );

    booksRessource.addMethod(
      "POST",
      new LambdaIntegration(booksLambdaFn, { proxy: true })
    );

    booksRessource.addMethod(
      "PUT",
      new LambdaIntegration(booksLambdaFn, { proxy: true })
    );

    booksRessource.addMethod(
      "DELETE",
      new LambdaIntegration(booksLambdaFn, { proxy: true })
    );

    const coversRessource = api.root.addResource("covers");

    coversRessource.addMethod(
      "POST",
      new LambdaIntegration(bookCoversLambdaFn, { proxy: true })
    );
  }
}
