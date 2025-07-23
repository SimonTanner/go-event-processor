#!/bin/sh

export AWS_REGION="eu-west-2"
export AWS_ACCESS_KEY_ID="test"
export AWS_SECRET_ACCESS_KEY="test"

ENDPOINT_URL="http://localhost:4566"
STREAM_NAME="eventStream"
FUNCTION_NAME="goEventProcessor"
DYNAMO_TABLE_NAME="GoEvents"

# echo "home: $HOME, user: $USER, localstack hostname: $HOSTNAME, $LOCALSTACK_HOST"

echo "[+] Creating Kinesis stream: $STREAM_NAME"

aws --endpoint-url="$ENDPOINT_URL" kinesis create-stream --stream-name "$STREAM_NAME" \
    --shard-count 1 \
    --region "$AWS_REGION" \
    --profile localstack

sleep 5

KINESIS_ARN=`(aws --endpoint-url="$ENDPOINT_URL" kinesis describe-stream --stream-name "$STREAM_NAME" \
    --region "$AWS_REGION" --profile localstack | jq .StreamDescription.StreamARN)`

KINESIS_ARN=`(echo ${KINESIS_ARN//\"})`

# echo "Kinesis ARN: $KINESIS_ARN"

echo "[+] building go binary"

GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bootstrap -tags lambda.norpc lambda/main.go
$GOPATH/bin/build-lambda-zip -o "$FUNCTION_NAME".zip bootstrap

echo "[+] Deploying Lambda"

aws --endpoint-url="$ENDPOINT_URL" lambda create-function --function-name "$FUNCTION_NAME" \
    --runtime provided.al2023 --handler bootstrap \
    --architectures x86_64 \
    --role arn:aws:iam::000000000000:role/lambda-ex \
    --zip-file fileb://"$FUNCTION_NAME".zip \
    --profile localstack >> /dev/null 2>&1

sleep 5

echo "[+] Adding event source mapping between Lambda and Kinesis Stream, ARN: $KINESIS_ARN"

aws --endpoint-url="$ENDPOINT_URL" lambda create-event-source-mapping --function-name "$FUNCTION_NAME" \
    --event-source $KINESIS_ARN \
    --batch-size 1 --starting-position LATEST \
    --profile localstack >> /dev/null 2>&1

echo "[+] Creating DynamoDB table"

aws --endpoint-url="$ENDPOINT_URL" dynamodb create-table \
    --table-name "$DYNAMO_TABLE_NAME" \
    --attribute-definitions \
        AttributeName=Client,AttributeType=S \
        AttributeName=CustomerID,AttributeType=S \
    --key-schema AttributeName=Client,KeyType=HASH AttributeName=CustomerID,KeyType=RANGE \
    --billing-mode PAY_PER_REQUEST \
    --table-class STANDARD >> /dev/null 2>&1

# aws --endpoint-url="$ENDPOINT_URL" dynamodb create-table \
#     --table-name "$DYNAMO_TABLE_NAME" \
#     --attribute-definitions \
#         AttributeName=Client,AttributeType=S \
#         AttributeName=CustomerID,AttributeType=S \
#         AttributeName=AccountID,AttributeType=S \
#         AttributeName=EventID,AttributeType=S \
#         AttributeName=Source,AttributeType=S \
#         AttributeName=Timestamp,AttributeType=N \
#         AttributeName=Type,AttributeType=S \
#         AttributeName=Payload,AttributeType=B \
#     --key-schema AttributeName=Client,KeyType=HASH AttributeName=AccountID,KeyType=RANGE \
#     --billing-mode PAY_PER_REQUEST \
#     --table-class STANDARD
    #     AttributeName=CustomerID,KeyType=RANGE AttributeName=EventID,KeyType=RANGE \
    #     AttributeName=Source,KeyType=RANGE AttributeName=Timestamp,KeyType=RANGE \
    #     AttributeName=CustomerID,KeyType=RANGE AttributeName=CustomerID,KeyType=RANGE \
    #     AttributeName=CustomerID,KeyType=RANGE AttributeName=CustomerID,KeyType=RANGE \

    # --billing-mode PAY_PER_REQUEST \
    # --table-class STANDARD
