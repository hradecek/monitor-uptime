# Uptime monitor
Environment variables:

 - TIMEOUT timeout in seconds
 - DYNAMO_TABLE_EXECUTIONS DynamoDB table name in which uptime's executions are stored
 - DYNAMO_TABLE_STATUS DynamoDB table name in which uptime's status is stored
 - SNS_TOPIC ARN of SNS topic to which are published changes of uptime's status

## Build
Using provided `Makefile`, run:
```
$ make
```

after successful build, AWS lambda zip and executable are present in the `build` directory.

Note, that subsequent `make` call, will clean (delete) the `build` directory.

