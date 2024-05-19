package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge/types"
)

type Request struct {
	ID    float64 `json:"id"`
	Value string  `json:"value"`
}

type Payload struct {
	ID      float64 `json:"id"`
	Value   string  `json:"value"`
	BusName string  `json:"busName"`
}

// type Resources struct {
// 	Resources []string `json:"resources"`
// }

// type Response struct {
// 	Detail  string `json:"detail"`
// 	Message string `json:"message"`
// 	Ok      bool   `json:"ok"`
// }

func Haddler(ctx context.Context, request Request) (string, error) {
	fmt.Println("test...")
	// var setResources Resources

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatal(err)
	}

	svc := eventbridge.NewFromConfig(cfg)

	setPayload := Payload{
		ID:      request.ID,
		Value:   request.Value,
		BusName: "send-notice-failure",
	}

	detail, err := json.Marshal(setPayload)
	if err != nil {
		return "", fmt.Errorf("Marshal fail to write in json %w", err)
	}

	setInput := &eventbridge.PutEventsInput{
		Entries: []types.PutEventsRequestEntry{
			{
				Time:         aws.Time(time.Now()),
				Detail:       aws.String(string(detail)),
				DetailType:   aws.String("Message"),
				EventBusName: aws.String("arn:aws:events:ap-southeast-1:058264531773:event-bus/event-bus-send-notice-fail"),
				Source:       aws.String("lambda publish"),
				// Resources:    setResources,
			},
		},
	}

	result, err := svc.PutEvents(ctx, setInput)
	fmt.Println("result => ", result)
	if err != nil {
		return "", fmt.Errorf("error PutEvents %w", err)
	}

	return "sucess send eventbus", nil
}

func main() {
	lambda.Start(Haddler)
}
