package pubnub

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/pubnub/go/messaging"
)

const (
	publishKey   = "pub-c-2a0450f8-865e-44c8-80ad-00a150b666ae"
	subscribeKey = "sub-c-2f63eada-2a02-11e8-9288-daa582a09445"
	secretKey    = "sec-c-YjRlMjFmMGYtMjFhZC00ZGQ3LTg0MjItNDdlYmVkOTE5YjJm"
)

var prohibitedChannelNameChars = []string{",", "/", "\\", ".", "*", ":"}

func GetPubNubVersion() string {
	return fmt.Sprintf("PubNub SDK for go: %s", messaging.VersionInfo())
}

func NewPubNub() (*messaging.Pubnub, error) {
	pn := messaging.NewPubnub(publishKey, subscribeKey, secretKey, "", true, "", nil)
	if pn == nil {
		return nil, fmt.Errorf("could not create new PubNub client")
	}
	return pn, nil
}

// ExtractMessage takes a byte array returned from PubNub success channel and
// returns the json in map[string]interface{} form.
func ExtractMessage(envelope []byte) (map[string]interface{}, error) {
	// Unmarshal byte array into []interface
	var ienvelope []interface{}

	err := json.Unmarshal(envelope, &ienvelope)
	if err != nil {
		return nil, fmt.Errorf("unmarshaling of message has failed: %s", err)
	}

	// if this is a published message first item in ienvelope will be another []interface{}
	switch m := ienvelope[0].(type) {
	case []interface{}:
		// first value of nested []interface{} will be the json
		jmap, ok := m[0].(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("failed to type assert message into json map[string]interface{}")
		}
		log.Printf("extracted message from PubNub: %v", jmap)
		return jmap, nil
	case float64:
		return nil, fmt.Errorf("diag message")
	default:
		return nil, fmt.Errorf("byte array did not container diag or published message")
	}
}

func PublishMessage(pnchan string, msg interface{}, pn *messaging.Pubnub) error {
	// Create channels to judge results
	successChannel := make(chan []byte)
	errorChannel := make(chan []byte)

	// Call publish
	go pn.Publish(pnchan, msg, successChannel, errorChannel)

	// Enter select to wait on response
	select {
	case _ = <-successChannel:
		return nil
	case err := <-errorChannel:
		return fmt.Errorf(string(err))
	case <-messaging.Timeout():
		return fmt.Errorf("publish message timeout")
	}
}
