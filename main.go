package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go/aws"
)

var (
	ctx          = context.Background()
	functionName string
)

func main() {

	flag.StringVar(&functionName, "functionName", "", "Name of the lambda function to proxy requests to")
	flag.Parse()

	if functionName == "" {
		log.Fatal("missing flag functionName")
	}

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatal(err)
	}

	lambdaClient := lambda.NewFromConfig(cfg)

	r := http.DefaultServeMux

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		output, err := lambdaClient.Invoke(ctx, &lambda.InvokeInput{
			FunctionName: aws.String(functionName),
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write(output.Payload)
	})

	fmt.Println("running. listening on port :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
