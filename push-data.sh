export AWS_REGION="eu-west-2"
export AWS_ACCESS_KEY_ID="test"
export AWS_SECRET_ACCESS_KEY="test"

ENDPOINT_URL="http://localhost:4566"
STREAM_NAME="eventStream"

TEST_DATA=`(echo testdata | base64)`
TEST_DATA=`(cat test-data.json | base64)`

aws --endpoint-url="$ENDPOINT_URL" kinesis put-record --stream-name "$STREAM_NAME" \
    --partition-key 1 --data "${TEST_DATA}" \
    --profile localstack

CONTAINER="$(sudo docker container ls -a --format json | grep lambda | jq .Names)"
CONTAINER=`(echo ${CONTAINER//\"})`

sudo docker logs $CONTAINER