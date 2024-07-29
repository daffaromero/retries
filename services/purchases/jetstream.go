package main

import (
	"log"

	config "github.com/daffaromero/retries/services/purchases/config"
	"github.com/nats-io/nats.go"
)

func JetStreamInit() (nats.JetStreamContext, error) {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	js, err := nc.JetStream(nats.PublishAsyncMaxPending(256))
	if err != nil {
		return nil, err
	}

	err = CreateStream(js)
	if err != nil {
		return nil, err
	}

	return js, nil
}

func CreateStream(jsc nats.JetStreamContext) error {
	stream, err := jsc.StreamInfo(config.StreamName)
	if err != nil {
		log.Printf("error retrieving stream info %v", err)
	}

	if stream == nil {
		log.Printf("creating stream: %s\n", config.StreamName)

		_, err = jsc.AddStream(&nats.StreamConfig{
			Name:     config.StreamName,
			Subjects: []string{config.StreamSubjects},
		})

		if err != nil {
			log.Printf("error creating stream: %s\n", err)
			return err
		}
	}
	return nil
}
