# SMTP Relay

SMTP Relay is a very simple SMTP server that will relay all incoming emails to a remote mail service.  
I use as a single entrypoint to relay all mails from my Docker containers to Mailgun.

[![Build Status](https://travis-ci.com/DenBeke/smtprelay.svg?branch=master)](https://travis-ci.com/DenBeke/smtprelay)
[![Go Report Card](https://goreportcard.com/badge/github.com/DenBeke/smtprelay)](https://goreportcard.com/report/github.com/DenBeke/smtprelay)
[![Docker Image Size (latest by date)](https://img.shields.io/docker/image-size/denbeke/smtprelay?sort=date)](https://hub.docker.com/r/denbeke/smtprelay)


## Usage (Docker)

### Docker-compose

The easiest way to run SMTP Relay is with docker-compose.
Edit the `.env` file with your settings,  download the [docker-compose.yml](./docker-compose.yml) file and run it with:

```bash
docker-compose up -d
```


### Docker run

If you don't want to use Docker compose, you can always run the command manually:

```bash
docker run -it\
    -e REMOTE_SMTP_HOST=${REMOTE_SMTP_HOST} \
    -e REMOTE_SMTP_PORT=${REMOTE_SMTP_PORT} \
    -e REMOTE_SMTP_DISABLE_TLS=${REMOTE_SMTP_DISABLE_TLS} \
    -e REMOTE_SMTP_USER=${REMOTE_SMTP_USER} \
    -e REMOTE_SMTP_PASSWORD=${REMOTE_SMTP_PASSWORD} \
    -p 25:25 \
    denbeke/smtprelay
```



## Usage (binary)

Download the latest SMTP Relay from the [releases page](https://github.com/DenBeke/smtprelay/releases).

Configure your settings in the `.env` and run the SMTP Relay with:

```bash
./smtprelay
```


## Development

Run it manually with Go (requires Go 1.15 or newer):

```bash
go run cmd/smtprelay/*.go
```

To test the email functionality, you can send the `test.txt` SMTP mail with a tool like netcat:

```bash
nc localhost 25 -i 1 < mail.txt
```


## Acknowledgments

- [gopistolet/smtp](https://github.com/gopistolet/smtp)
- [sirupsen/logrus](https://github.com/sirupsen/logrus)
- [evalphobia/logrus_sentry](https://github.com/evalphobia/logrus_sentry)
- [emersion/go-smtp](https://github.com/emersion/go-smtp)



## Author

[Mathias Beke](https://denbeke.be)

