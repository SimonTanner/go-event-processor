
export AWS_REGION="eu-west-2"

ENDPOINT_URL="http://localhost:4566"
STREAM_NAME="eventStream"
FUNCTION_NAME="goEventProcessor"
DYNAMO_TABLE_NAME="GoEvents"

aws --endpoint-url="$ENDPOINT_URL" dynamodb get-item \
    --table-name "$DYNAMO_TABLE_NAME" \
    --region "$AWS_REGION" \
    --key='{"Client":{"S":"barclays"}, "CustomerID":{"S":"0124e053-3580-7000-a762-0502e4a1022e"}}' \
    --profile localstack 