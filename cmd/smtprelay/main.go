package main

import (
	log "github.com/sirupsen/logrus"

	smtprelay "github.com/DenBeke/smtp-relay"
)

func main() {

	log.SetLevel(log.DebugLevel)

	config := smtprelay.BuildConfigFromEnv()
	SMTPRelay, err := smtprelay.New(config)
	if err != nil {
		log.Fatalf("couldn't create SMTP Relay instance: %v", err)
	}

	SMTPRelay.Serve()

}
