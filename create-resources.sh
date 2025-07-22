#!/bin/sh

# sleep 5

export AWS_REGION="eu-west-2"
export AWS_ACCESS_KEY_ID="test"
export AWS_SECRET_ACCESS_KEY="test"

ENDPOINT_URL="http://host.docker.internal:4566"
ENDPOINT_URL="http://172.17.0.2:4566"
ENDPOINT_URL="http://localhost:4566"
# ENDPOINT_URL="http://127.0.0.1:4566"
# ENDPOINT_URL="http://simon-HP-ENVY-Laptop:4566"

STREAM_NAME="eventStream"


# aws configure

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

echo "Kinesis ARN: $KINESIS_ARN"

# aws --endpoint-url="$ENDPOINT_URL" kinesis describe-stream --stream-name "$STREAM_NAME" \
#     --region "$AWS_REGION" --profile localstack > kinesis

# # echo "Stream details: \n $KINESIS_INFO"

GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bootstrap -tags lambda.norpc lambda/main.go
build-lambda-zip -o myFunction.zip bootstrap

aws --endpoint-url="$ENDPOINT_URL" lambda create-function --function-name myFunction \
    --runtime provided.al2023 --handler bootstrap \
    --architectures x86_64 \
    --role arn:aws:iam::000000000000:role/lambda-ex \
    --zip-file fileb://myFunction.zip \
    --profile localstack > /dev/null 2>&1

sleep 5

aws --endpoint-url="$ENDPOINT_URL" lambda create-event-source-mapping --function-name myFunction \
    --event-source $KINESIS_ARN \
    --batch-size 1 --starting-position LATEST \
    --profile localstack > /dev/null 2>&1

# ls -al ~/

# cat ~/.aws/config 
# cat /etc/hosts
# echo "\n"

# sleep 5


# cat ~/.aws/credentials 