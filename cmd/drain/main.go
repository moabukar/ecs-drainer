package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	drain "github.com/moabukar/ecs-drainer"
)

func main() {
	// main entry point for the lambda function
	lambda.Start(drain.HandleRequest)
}
