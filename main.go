package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go/aws"
	flag "github.com/spf13/pflag"
)

var (
	ctx          = context.Background()
	help         bool
	functionName string
	listenAddr   string
	listenPort   int
)

func main() {

	flag.StringVarP(&functionName, "functionName", "f", "", "Name or ARN of the lambda function to proxy requests to")
	flag.StringVarP(&listenAddr, "listenAddr", "a", "0.0.0.0", "Address to listen on")
	flag.IntVarP(&listenPort, "listenPort", "p", 8080, "Port to listen on")
	flag.BoolVarP(&help, "help", "h", false, "Print help")
	flag.Parse()

	if help {
		flag.PrintDefaults()
		os.Exit(0)
	}

	if functionName == "" {
		flag.PrintDefaults()
		log.Fatal("missing flag functionName")
	}

	functionName = strings.TrimSpace(functionName)
	listenAddr = strings.TrimSpace(listenAddr)

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

	listen := fmt.Sprintf("%s:%d", listenAddr, listenPort)

	log.Printf("running. listening on %s", listen)
	if err := http.ListenAndServe(listen, r); err != nil {
		log.Fatal(err)
	}
}
