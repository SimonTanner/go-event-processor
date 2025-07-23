# go-event-processor

The Challenge to produce an event driven solution which allows receiving data from multiple sources with different types which require validating and persisting so that these can be processed by multiple clients at some later point.

## Idea

Using a kinesis stream which can deliver high throughput low-latency ingestion of events and useful when the time ordering of events is important. This is connected to a lambda which processes the events, validates the data and then persists these. Due to the nature of these I chose to use dynamodb for the flexibility due to being schemaless as well as allowing high throughput and the ability to partition data, based on a chosen attribute.

I tried to think of financial applications and the simple kind of events that could come from different sources, such as a transaction, an alert for fraudulent activity leading to tighter monitoring of account or freezing certain actions, etc.

## Solution

I chose to use localstack to allow simulating AWS's services locally and have included certain scripts to facilitate testing this.

## Running locally

The go code was written using go1.24 and gvm for managing multiple versions. My local machine is Ubuntu 24.04, which posed numerous issues with running docker so you might need to tweak some of the scripts. This is also why some commands have to be run as root.

Before creating the go binary for the lamba you must have the go/aws lambda zip binary installed:
```go install github.com/aws/aws-lambda-go/cmd/build-lambda-zip@latest```

This should install it to `$GOPATH/bin/build-lambda-zip` however this could be different depending on the os.

In order to start the localstack container enter:

```docker compose up --build```

Once the localstack container is in a ready state you can run the bash script as follows:

```bash create-resources.sh```

__N.B__ this requires `jq` in order to work

this will create the different infrastructure:
1. Kinesis Stream
2. Lambda
3. Lambda Kinesis Event Source Mapping
4. DynamoDB Table

In order to send data to the stream there's some test data and a script which will send this, by running:

```bash put-records.sh```

Unfortunately in order to see the logs for the lambda you need to call docker log with the container ID that is spawned by localstack in order to run lambdas, not the main localstack container. I've added a script `get-lambda-logs.sh` for simplicity. This requires `jq` in order to get the correct container.




