package pubnub

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

// These are primarily integration tests with PubNub. Make sure you have network or elese
// these will fail

func TestSDKVersion(t *testing.T) {
	version := GetPubNubVersion()
	if version == "" {
		t.Error("GetPubNubVersion returned empty string")
	}

	t.Log(version)
}

func TestNewPubNub(t *testing.T) {
	_, err := NewPubNub()
	if err != nil {
		t.Error(err.Error())
	}

}

func TestExtracAndPublishMessage(t *testing.T) {
	// Create channels where Subscribe method will push messages onto
	successChannel := make(chan []byte)
	errorChannel := make(chan []byte)

	// Get new PubNub client
	pn, err := NewPubNub()
	if err != nil {
		t.Fatalf("could not create PubNub client: %s", err)
	}

	// Subscribe to a channel
	go pn.Subscribe("test", "", successChannel, false, errorChannel)

	// create a message
	message := struct{ Text string }{
		Text: "success",
	}

	// Test PublishMessage function
	err = PublishMessage("test", message, pn)
	if err != nil {
		t.Fatalf("PublishMessage function returned error: %s", err)
	}

	// Select on channels to determine if round trip was successful. Use timer so test
	// does not hang forever
	var timeoutSeconds time.Duration = 4
	timer := time.NewTimer(time.Second * timeoutSeconds)
	for {
		select {
		case response := <-successChannel:
			json, err := ExtractMessage(response)
			if err != nil {
				if err.Error() != "diag message" {
					t.Fatal("encountered non diag error: %s", err)
				}
				continue
			}
			// found json
			if json["Text"] == "success" {
				return
			}
			if json["Text"] != "success" {
				t.Fatalf("received message text other then what was published: %s", json["Text"])
			}
		case <-timer.C:
			t.Fatalf("did not receive published message in %s second(s)", timeoutSeconds)
		}
	}
}

func TestRoundTrip(t *testing.T) {

	// Create channels where Subscribe method will push messages onto
	successChannel := make(chan []byte)
	errorChannel := make(chan []byte)

	// Get new PubNub client
	pn, err := NewPubNub()
	if err != nil {
		t.Fatalf("could not create PubNub client: %s", err)
	}

	// Subscribe to a channel
	go pn.Subscribe("test", "", successChannel, false, errorChannel)

	// publish a message
	psuccessChannel := make(chan []byte)
	perrorChannel := make(chan []byte)
	message := struct{ Text string }{
		Text: "success",
	}

	// Publish message
	go pn.Publish("test", message, psuccessChannel, perrorChannel)

	// Select on channels to determine if round trip was successful. Use timer so test
	// does not hang forever
	var timeoutSeconds time.Duration = 4
	timer := time.NewTimer(time.Second * timeoutSeconds)
	for {
		select {
		case response := <-successChannel:
			var msg []interface{}

			type Resp struct{ Text string }

			err := json.Unmarshal(response, &msg)
			if err != nil {
				fmt.Println(err)
				return
			}

			switch m := msg[0].(type) {
			case float64:
				t.Logf("found diagnostic message of len: %d", len(msg))
				t.Log(msg[0])
				t.Log(msg[1].(string))
				t.Log(msg[2])
			case []interface{}:
				t.Logf("found message of len: %d", len(msg))
				t.Log(msg[0])
				t.Log(msg[1])
				t.Log(msg[2])
				mapp := m[0].(map[string]interface{})
				// Make sure message value is what we sent
				if mapp["Text"] != "success" {
					t.Fatalf("received message with different text field: %s", mapp["Test"])
				}

				return
			default:
				t.Fatalf(fmt.Sprintf("Unknown type: %T", m))
			}

		case err := <-errorChannel:
			t.Fatalf("received error on subscription: %s", err)
		case err := <-perrorChannel:
			t.Fatalf("received error on publish: %s", err)
		case <-timer.C:
			t.Fatalf("did not receive published message in %s second(s)", timeoutSeconds)
		}

	}
}
