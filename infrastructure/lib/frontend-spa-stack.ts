import * as cdk from "aws-cdk-lib";
import { Certificate } from "aws-cdk-lib/aws-certificatemanager";
import {
  CloudFrontWebDistribution,
  IDistribution,
  OriginAccessIdentity,
  ViewerCertificate,
} from "aws-cdk-lib/aws-cloudfront";
import { PolicyStatement } from "aws-cdk-lib/aws-iam";
import { Bucket } from "aws-cdk-lib/aws-s3";
import { BucketDeployment, Source } from "aws-cdk-lib/aws-s3-deployment";
import { Construct } from "constructs";

interface FrontendSpaStackProps extends cdk.StackProps {
  domainName: string;
  certificate: Certificate;
}

export class FrontendSpaStack extends cdk.Stack {
  public distribution: IDistribution;

  constructor(scope: Construct, id: string, props: FrontendSpaStackProps) {
    super(scope, id, props);

    const websiteBucket = new Bucket(this, "WebsiteBucket", {
      websiteIndexDocument: "index.html",
      websiteErrorDocument: "index.html",
      publicReadAccess: false,
    });

    const accessIdentity = new OriginAccessIdentity(
      this,
      "OriginAccessIdentity",
      { comment: `${websiteBucket.bucketName}-access-identity` }
    );

    websiteBucket.addToResourcePolicy(
      new PolicyStatement({
        actions: ["s3:GetObject"],
        resources: [websiteBucket.arnForObjects("*")],
        principals: [accessIdentity.grantPrincipal],
      })
    );

    this.distribution = new CloudFrontWebDistribution(
      this,
      "cloudfrontDistribution",
      {
        originConfigs: [
          {
            s3OriginSource: {
              s3BucketSource: websiteBucket,
              originAccessIdentity: accessIdentity,
            },
            behaviors: [{ isDefaultBehavior: true }],
          },
        ],
        errorConfigurations: [
          {
            errorCode: 403,
            responseCode: 200,
            responsePagePath: "/index.html",
          },
          {
            errorCode: 404,
            responseCode: 200,
            responsePagePath: "/index.html",
          },
        ],
        viewerCertificate: ViewerCertificate.fromAcmCertificate(
          props.certificate,
          {
            aliases: [props.domainName],
          }
        ),
      }
    );

    new BucketDeployment(this, "BucketDeployment", {
      sources: [Source.asset("../spa/build")],
      destinationBucket: websiteBucket,
      // Invalidate the cache for / and index.html when we deploy so that cloudfront serves latest site
      distribution: this.distribution,
      distributionPaths: ["/", `/index.html`],
    });

    new cdk.CfnOutput(this, "cloudfront domain", {
      description: "The domain of the website",
      value: this.distribution.distributionDomainName,
    });
  }
}
