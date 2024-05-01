package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	drain "github.com/moabukar/ecs-drainer"
)

func main() {
	lambda.Start(drain.HandleRequest)
}
