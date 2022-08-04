package smtprelay

import (
	"bytes"
	"fmt"
	"runtime/debug"

	"github.com/emersion/go-sasl"
	gosmtp "github.com/emersion/go-smtp"
	"github.com/gopistolet/smtp/smtp"
	log "github.com/sirupsen/logrus"
)

// Handle implements the GoPistolet SMTP Handler interface
func (r *SMTPRelay) Handle(s *smtp.State) error {

	if r.config.HandleAsync {
		log.Debugln("handling incoming mail async")
		state := *s
		go func() {
			err := r.handleMail(&state)
			if err != nil {
				log.Errorf("couldn't handle incoming mail: %v", err)
			}
		}()
	} else {
		log.Debugln("handling incoming mail sync")
		err := r.handleMail(s)
		if err != nil {
			log.Errorf("couldn't handle incoming mail: %v", err)
			return err
		}
	}

	return nil

}

func (r *SMTPRelay) handleMail(s *smtp.State) error {

	// Recover from panics
	if r.config.SentryDSN != "" {
		defer func() {
			if err := recover(); err != nil {
				log.WithField("stacktrace", string(debug.Stack())).Fatalf("panic: %s", err)
			}
		}()
	}

	log.Println("parsing incoming mail...")
	//log.Debugf("%+v", s)

	reader := bytes.NewReader(s.Data)

	from := s.From.Address
	to := func() []string {
		addresses := []string{}
		for _, address := range s.To {
			addresses = append(addresses, address.Address)
		}
		return addresses
	}()
	remoteSMTPAddress := fmt.Sprintf("%s:%d", r.config.RemoteSMTP.Host, r.config.RemoteSMTP.Port)
	auth := sasl.NewPlainClient("", r.config.RemoteSMTP.User, r.config.RemoteSMTP.Password)

	err := gosmtp.SendMail(remoteSMTPAddress, auth, from, to, reader)
	if err != nil {
		log.Errorln("couldn't send email: %v", err)
		return fmt.Errorf("couldn't send email: %v", err)
	}

	log.Println("Relayed email to original recipient.")

	return nil
}
