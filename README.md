# Uptime monitor
Environment variables:

- `TIMEOUT` - Timeout in seconds
- `DYNAMO_TABLE_EXECUTIONS` - DynamoDB table name in which uptime's executions are stored
- `DYNAMO_TABLE_STATUS` - DynamoDB table name in which uptime's status is stored
- `SNS_TOPIC` - ARN of SNS topic to which are published changes of uptime's status

## Build
Make sure you have installed [build-lambda-zip](https://github.com/aws/aws-lambda-go/tree/master/cmd/build-lambda-zip) tool.\
In order to install it, run:
```
$ go get -u github.com/aws/aws-lambda-go/cmd/build-lambda-zip
```

For building lambda ZIP file use provided `Makefile`, run:
```
$ make
```

after successful build, AWS lambda zip and executable are present in the `build` directory.

Note, that subsequent `make` call, will clean (delete) the `build` directory.
