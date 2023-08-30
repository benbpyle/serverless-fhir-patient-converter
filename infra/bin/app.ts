#!/usr/bin/env node
import "source-map-support/register";
import * as cdk from "aws-cdk-lib";
import { MainStack } from "../lib/main-stack";

const app = new cdk.App();

new MainStack(app, `PatientFhirConverterStack`, {
    tags: {
        billingTag: "patient-fhir-converter",
        service: "patient-fhir-converter",
    },
});
