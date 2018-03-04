# S3 read/write

This is a silly/simple AWS/Golang/Lambda example to demonstrate how Go Lambdas can be deployed with the Serverless framework.

When .txt file is uploaded to a S3 bucket the deployed Lambda will trigger and read the file and write a copy of it into same S3 bucket.

## Prerequisites

* AWS account
* basic knowledge how to work with S3 buckets
* [Go](https://golang.org/)
* [dep](https://github.com/golang/dep) (go package manager)
* [Serverless](https://github.com/serverless/serverless)

## How to try it yourself

* Edit `serverless.yml` to match your AWS profile definition and your bucket name
  * Note: the Serverless will create bucket with the given name, so don't create it yourself
* Deploy with `serverless deploy`
* Throw some file (with .txt suffix) into your bucket and see how it gets duplicated by the triggered Lambda
