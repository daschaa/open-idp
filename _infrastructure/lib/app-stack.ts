import { GoFunction } from "@aws-cdk/aws-lambda-go-alpha";
import { Stack, StackProps } from "aws-cdk-lib";
import { Key, KeySpec, KeyUsage } from "aws-cdk-lib/aws-kms";
import { Construct } from "constructs";
export class AppStack extends Stack {
  constructor(scope: Construct, id: string, props?: StackProps) {
    super(scope, id, props);
    new Key(this, "Key", {
      keySpec: KeySpec.HMAC_256,
      keyUsage: KeyUsage.GENERATE_VERIFY_MAC,
    });
    new GoFunction(this, "Function", {
      entry: __dirname + "/../../cmd/idp-server/main.go",
      moduleDir: __dirname + "/../../go.mod",
      functionName: "idp-server",
    });
  }
}
