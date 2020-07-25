package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Azure/azure-storage-queue-go/azqueue"
)

func main() {
	fmt.Println("Azure Strage Queue Sender")
	connectionString := os.Getenv("ConnectionString")
	queueName := os.Getenv("QueueName")
	if len(os.Args) != 2 {
		log.Fatalf("Specify the counter parameter. e.g. send 100 Parameter length: %d\n", len(os.Args))
	}
	count, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatalf("count should be integer : %s", os.Args[1])
	}
	credential, endpoint, err := parseAzureStorageConnectionString(connectionString)
	if err != nil {
		log.Fatalf("%s: %s", "Invalid ConnectionString", err)
	}
	p := azqueue.NewPipeline(credential, azqueue.PipelineOptions{})
	serviceURL := azqueue.NewServiceURL(*endpoint, p)
	queueURL := serviceURL.NewQueueURL(queueName)
	ctx := context.Background()
	_, err = queueURL.Create(ctx, azqueue.Metadata{})
	if err != nil {
		log.Fatalf("%s: %s", "Can not create a queue", err)
	}

	messageURL := queueURL.NewMessagesURL()
	_, err = messageURL.Clear(ctx)
	if err != nil {
		log.Fatalf("Cannot clear the message: %v", err)
	}
	for i := 0; i < count; i++ {
		props, err := queueURL.GetProperties(ctx)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("1: Approximate number of messages in the queue=%d\n", props.ApproximateMessagesCount())
		visibleCount, err := getVisibleCount(&messageURL, 32)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("2: Visible count of messages in the queue=%d\n", visibleCount)

		response, err := messageURL.Enqueue(ctx, "Hello!", 30*time.Second, 0) // TimeToLive = 0 means 7 days.
		fmt.Printf("3: Enqueue Hello %d: %v \n", i, response.MessageID)
		props, err = queueURL.GetProperties(ctx)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("4: Approximate number of messages in the queue=%d\n", props.ApproximateMessagesCount())
		visibleCount, err = getVisibleCount(&messageURL, 32)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("5: Visible count of messages in the queue=%d\n", visibleCount)
	}
}

func getVisibleCount(messagesURL *azqueue.MessagesURL, maxCount int32) (int32, error) {
	ctx := context.Background()
	queue, err := messagesURL.Peek(ctx, maxCount)
	if err != nil {
		return 0, err
	}
	num := queue.NumMessages()
	return num, nil
}

func parseAzureStorageConnectionString(connectionString string) (azqueue.Credential, *url.URL, error) {
	parts := strings.Split(connectionString, ";")

	getValue := func(pair string) string {
		parts := strings.SplitN(pair, "=", 2)
		if len(parts) == 2 {
			return parts[1]
		}
		return ""
	}

	var name, key string
	for _, v := range parts {
		if strings.HasPrefix(v, "AccountName") {
			name = getValue(v)
		} else if strings.HasPrefix(v, "AccountKey") {
			key = getValue(v)
		}
	}

	credential, err := azqueue.NewSharedKeyCredential(name, key)
	if err != nil {
		return nil, nil, err
	}

	if name == "" || key == "" {
		return nil, nil, errors.New("can't parse storage connection string. Missing key or name")
	}

	endpoint, _ := url.Parse(fmt.Sprintf("https://%s.queue.core.windows.net", name))

	return credential, endpoint, nil
}
