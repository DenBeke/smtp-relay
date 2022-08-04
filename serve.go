package smtprelay

import (
	"fmt"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"

	"github.com/evalphobia/logrus_sentry"
	"github.com/gopistolet/smtp/mta"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

const tmpDir = "./tmp/"

// Serve runs the actual SMTP server and handles all the mails
func (r *SMTPRelay) Serve() {

	// Validate config
	err := r.config.Validate()
	if err != nil {
		log.Fatalf("Config file is not valid: %v", err)
	}

	// Error logging with Sentry
	if r.config.SentryDSN != "" {
		hook, err := logrus_sentry.NewSentryHook(r.config.SentryDSN, []logrus.Level{
			logrus.PanicLevel,
			logrus.FatalLevel,
			logrus.ErrorLevel,
		})

		hook.StacktraceConfiguration.Enable = true

		if err == nil {
			log.AddHook(hook)
		}

		defer func() {
			if err := recover(); err != nil {
				log.WithField("stacktrace", string(debug.Stack())).Fatalf("panic: %s", err)
			}
		}()
	}

	log.WithField("config", fmt.Sprintf("%+v", r.config)).Println("Starting SMTP Relay ✉️")

	if r.config.HandleAsync {
		log.Warnln("Async mail handling is enabled. Mails might 'disappear' silently on crashes or errors.")
	}

	// Configure and start GoPistolet SMTP server
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt, syscall.SIGTERM)

	// Default config
	smtpConfig := mta.Config{
		Hostname: "localhost",
		Port:     25,
	}

	// create new MTA with SMTP config and SMTP Relay as the email handler
	mta := mta.NewDefault(smtpConfig, r)
	go func() {
		<-sigc
		mta.Stop()
	}()
	err = mta.ListenAndServe()
	if err != nil {
		log.Fatalln(err)
	}
}
