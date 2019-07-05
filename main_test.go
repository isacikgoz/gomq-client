package main

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/isacikgoz/gomq/api"
)

type mockReader struct {
	bytes []byte
}

func (mr *mockReader) Read(bytes []byte) (int, error) {
	message := &BareMessage{
		Message: "testing",
	}
	payload, _ := json.Marshal(message)
	mockMsg := &api.AnnotatedMessage{
		Command: "PUB",
		Target:  "topic_1",
		Payload: payload,
	}
	src, _ := json.Marshal(mockMsg)
	for i, b := range src {
		bytes[i] = b
	}
	return len(src), nil
}

func TestListenFromUser(t *testing.T) {
	messages := make(chan api.AnnotatedMessage)

	r := strings.NewReader("hak ben")

	go listenFromUser(r, messages)

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	select {
	case <-ctx.Done():
		t.Fatal("test failed")
	case msg := <-messages:
		t.Logf("test passed: %v", msg)
		break
	}
}

func TestListenFromBroker(t *testing.T) {
	messages := make(chan api.AnnotatedMessage)

	reader := &mockReader{}
	go listenFromBroker(reader, messages)

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	select {
	case <-ctx.Done():
		t.Fatal("test failed")
	case msg := <-messages:
		t.Logf("test passed: %v", msg)
		break
	}
}
