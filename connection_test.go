package main

import (
	"errors"
	"os"
	"testing"
	"time"
)

// TODO: FIX TEST
func TestMessageSending(t *testing.T) {
	log := createLogFile()
	//start messageListener
	go func() {
		err := messageListener(log)
		if err != nil {
			os.Exit(1)
		}
	}()

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
