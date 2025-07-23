export AWS_REGION="eu-west-2"
export AWS_ACCESS_KEY_ID="test"
export AWS_SECRET_ACCESS_KEY="test"

ENDPOINT_URL="http://localhost:4566"
STREAM_NAME="eventStream"

test_files=`(ls test-data)`
test_files=`(echo "${test_files//$'\n'/ }")`

for file in $test_files; do
    TEST_DATA=`(cat test-data/$file | base64)`

    echo $TEST_DATA

    aws --endpoint-url="$ENDPOINT_URL" kinesis put-record --stream-name "$STREAM_NAME" \
        --partition-key 1 --data "${TEST_DATA}"
done