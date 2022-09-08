package main

import (
	"errors"
	"testing"
	"time"
)

func TestMessageSending(t *testing.T) {
	log := createLogFile()
	//start listener
	err := go listener(log)

	time.Sleep(time.Duration(time.Second * 2))

	// expect message to be sent successfully
	if err := sendRequest("test", log); err != nil {
		if err != nil {
			t.Errorf("Expected %v but got %s", nil, err.Error())
		}
	}

	expectation := errors.New("Message was empty, exiting")
	if err := sendRequest("", log); err != nil {
		if err == expectation {
			return
		} else {
			t.Fatalf("Expected %s but got %s", expectation, err)
		}
	}
}
