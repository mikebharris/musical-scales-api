package main

import (
	"musical-scales-lambda/handler"

	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(handler.Handler{}.HandleRequest)
}
