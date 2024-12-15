import {GoFunction} from "@aws-cdk/aws-lambda-go-alpha";
import {CfnOutput, RemovalPolicy, Stack, StackProps} from "aws-cdk-lib";
import {Key, KeySpec, KeyUsage} from "aws-cdk-lib/aws-kms";
import {Construct} from "constructs";
import {HttpApi, HttpMethod} from "aws-cdk-lib/aws-apigatewayv2";
import {HttpLambdaIntegration} from "aws-cdk-lib/aws-apigatewayv2-integrations";
import {AttributeType, BillingMode, Table} from "aws-cdk-lib/aws-dynamodb";

export class AppStack extends Stack {
    constructor(scope: Construct, id: string, props?: StackProps) {
        super(scope, id, props);
       const table = new Table(this, "Table", {
            billingMode: BillingMode.PAY_PER_REQUEST,
            tableName: 'clients',
            partitionKey: {
                type: AttributeType.STRING,
                name: 'clientId'
            },
           removalPolicy: RemovalPolicy.DESTROY,
        })
        new Key(this, "Key", {
            keySpec: KeySpec.HMAC_256,
            keyUsage: KeyUsage.GENERATE_VERIFY_MAC,
        });
        const fn = new GoFunction(this, "Function", {
            entry: __dirname + "/../../cmd/idp-server/main.go",
            moduleDir: __dirname + "/../../go.mod",
            functionName: "idp-server",
        });
        table.grantFullAccess(fn);
        const httpApi = new HttpApi(this, 'HttpApi', {
            apiName: 'idp-server',
            description: 'This is the API for the IDP server',
        });
        httpApi.addRoutes({
            path: '/token',
            methods: [HttpMethod.POST],
            integration: new HttpLambdaIntegration('Integration', fn)
        });
        httpApi.addRoutes({
            path: '/introspect',
            methods: [HttpMethod.POST],
            integration: new HttpLambdaIntegration('Integration', fn)
        });
        new CfnOutput(this, 'ApiUrl', {
            value: httpApi.apiEndpoint,
            key: 'ApiUrl',
            description: 'The URL of the API',
        });
    }
}
