#!/usr/bin/env node
import "source-map-support/register";
import * as cdk from "aws-cdk-lib";
import { InfrastructureStack } from "../lib/infrastructure-stack";
import { HttpApiStack } from "../lib/lambda-rest-api-stack";
import { FrontendSpaStack } from "../lib/frontend-spa-stack";
import { CertificateStack } from "../lib/certificate-stack";
import { DnsStack } from "../lib/dns-stack";

// Purchased domain name and corresponding HostedZone from AWS R53
const domainName = "DOMAIN-NAME.COM HERE";
const hostedZoneId = "HOSTED-ZONE-HERE";

const app = new cdk.App();
const env = { account: "ACCOUNT_ID", region: "REGION" };

const { lambdaFunction } = new HttpApiStack(app, "RestApiStack", { env });
const { hostedZone, certificate } = new CertificateStack(
  app,
  "CertificateStack",
  { env, domainName, hostedZoneId }
);
const { distribution } = new FrontendSpaStack(app, "FrontendSpaStack", {
  env,
  domainName,
  certificate,
});
new DnsStack(app, "DnsStack", { env, hostedZone, distribution, domainName });
new InfrastructureStack(app, "InfrastructureStack", {
  env,
  lambdaFunction,
});
