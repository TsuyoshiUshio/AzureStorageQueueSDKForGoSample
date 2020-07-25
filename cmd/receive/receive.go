package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/Azure/azure-storage-queue-go/azqueue"
)

func main() {
	fmt.Println("Azure Strage Queue Receiver")
	connectionString := os.Getenv("ConnectionString")
	queueName := os.Getenv("QueueName")
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
	for {
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
		dequeue, err := messageURL.Dequeue(ctx, 1, 30*time.Second)
		if err != nil {
			log.Fatal(err)
		}
		if dequeue.NumMessages() == 0 {
			time.Sleep(time.Second * 10)
		} else {
			for m := int32(0); m < dequeue.NumMessages(); m++ {

				message := dequeue.Message(m)
				// txt, err := base64.StdEncoding.DecodeString(message.Text)

				fmt.Printf("3: Dequeue [%s] : Message: %s\n", message.InsertionTime.Format("2009-11-10 23:00:00 +0000 UTC"), message.Text)
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
				msgIDURL := messageURL.NewMessageIDURL(message.ID)
				_, err = msgIDURL.Delete(ctx, message.PopReceipt)
				fmt.Printf("6: Delete message %v\n", message.ID)
				if err != nil {
					log.Fatal(err)
				}
				fmt.Printf("7: Approximate number of messages in the queue=%d\n", props.ApproximateMessagesCount())
				visibleCount, err = getVisibleCount(&messageURL, 32)
				if err != nil {
					log.Fatal(err)
				}
				fmt.Printf("8: Visible count of messages in the queue=%d\n", visibleCount)
			}
			// return
		}
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
