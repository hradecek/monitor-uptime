# Uptime service
Environment variables:

 - TIMEOUT timeout in seconds
 - DYNAMO_TABLE DynamoDB table name
 - SNS_TOPIC ARN of SNS topic

## Build
On Windows (CMD):
```
> set GOOS=linux
> go build -o main main.go uptime.go dynamodb.go
```

## AWS Lambda zip
On Windows (CMD):
```
> %USERPROFILE%\Go\bin\build-lambda-zip.exe main
```
