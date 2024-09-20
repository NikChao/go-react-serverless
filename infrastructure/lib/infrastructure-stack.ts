import * as cdk from "aws-cdk-lib";
import { AttributeType, Table } from "aws-cdk-lib/aws-dynamodb";
import {
  Effect,
  IGrantable,
  PolicyDocument,
  PolicyStatement,
  Role,
  ServicePrincipal,
} from "aws-cdk-lib/aws-iam";
import { Bucket, HttpMethods } from "aws-cdk-lib/aws-s3";
import { Construct } from "constructs";

interface InfrastructureStackProps {
  lambdaFunction: IGrantable;
}

export class InfrastructureStack extends cdk.Stack {
  public readonly householdsTable: Table;
  public readonly tasksTable: Table;
  public readonly groceriesTable: Table;
  public readonly usersTable: Table;
  public readonly expensesTable: Table;
  public readonly catalogBucket: Bucket;
  public readonly unprocessedReceiptsBucket: Bucket;
  public readonly processedReceiptsBucket: Bucket;

  constructor(
    scope: Construct,
    id: string,
    props?: cdk.StackProps & InfrastructureStackProps
  ) {
    super(scope, id, props);

    this.groceriesTable = new Table(this, "Groceries", {
      tableName: "Groceries",
      partitionKey: {
        type: AttributeType.STRING,
        name: "householdId",
      },
      sortKey: {
        type: AttributeType.STRING,
        name: "id",
      },
    });
    this.groceriesTable.grantFullAccess(props!.lambdaFunction);

    this.householdsTable = new Table(this, "Households", {
      tableName: "Households",
      partitionKey: {
        type: AttributeType.STRING,
        name: "id",
      },
    });
    this.householdsTable.grantFullAccess(props!.lambdaFunction);

    this.usersTable = new Table(this, "Users", {
      tableName: "Users",
      partitionKey: {
        type: AttributeType.STRING,
        name: "id",
      },
    });
    this.usersTable.grantFullAccess(props!.lambdaFunction);

    this.catalogBucket = new Bucket(this, "CatalogBucket", {
      bucketName: "store-comparison-bucket-001",
    });
    this.catalogBucket.grantReadWrite(props!.lambdaFunction);

    this.unprocessedReceiptsBucket = new Bucket(this, "UnprocessedReceipts", {
      bucketName: "unprocessed-receipts-001",
      cors: [
        {
          allowedMethods: [
            HttpMethods.GET,
            HttpMethods.PUT,
            HttpMethods.POST,
            HttpMethods.DELETE,
            HttpMethods.HEAD,
          ],
          allowedOrigins: ["*"],
          allowedHeaders: ["*"],
          exposedHeaders: [
            "x-amz-server-side-encryption",
            "x-amz-request-id",
            "x-amz-id-2",
            "ETag",
          ],
          maxAge: 3000,
        },
      ],
    });
    this.unprocessedReceiptsBucket.grantReadWrite(props!.lambdaFunction);
    this.createPresignedUrlRole(this.unprocessedReceiptsBucket);

    this.processedReceiptsBucket = new Bucket(this, "ProcessedReceipts", {
      bucketName: "processed-receipts-001",
    });
  }

  private createPresignedUrlRole(bucket: Bucket) {
    return new Role(this, "PresignedUrlRole", {
      assumedBy: new ServicePrincipal("lambda.amazonaws.com"),
      inlinePolicies: {
        AllowS3BucketObjectAccess: new PolicyDocument({
          statements: [
            new PolicyStatement({
              effect: Effect.ALLOW,
              actions: ["s3:GetObject", "s3:PutObject"],
              resources: [bucket.bucketArn + "/*"],
            }),
          ],
        }),
      },
    });
  }
}
