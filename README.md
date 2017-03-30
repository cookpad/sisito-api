# sisito-api

[Sisito](https://github.com/winebarrel/sisito) API server.

[![Build Status](https://travis-ci.org/winebarrel/sisito-api.svg?branch=master)](https://travis-ci.org/winebarrel/sisito-api)

## Getting Started

```sh
cd docker
docker-compose build
docker-compose up
```

```sh
$ curl -u foo:bar localhost:8080/blacklist
{"recipients":["foo@example.com"]}
```

## API

```sh
$ curl -u foo:bar localhost:8080/blacklist
{"recipients":["foo@example.com"]}
```

```sh
$ curl -u foo:bar localhost:8080/listed?recipient=foo@example.com
{"listed":true}
```

```sh
$ curl -s -u foo:bar localhost:8080/recent?recipient=foo@example.com | jq .
{
  "addresser": "no-reply@sender.example.com",
  "alias": "foo@example.com",
  "created_at": "2017-03-01T00:00:00Z",
  "deliverystatus": "5.0.0",
  "destination": "example.com",
  "diagnosticcode": "550 Unknown user foo@example.com",
  "digest": "767e74eab7081c41e0b83630511139d130249666",
  "lhost": "mail.sender.example.com",
  "messageid": "foo_example_com_message_id",
  "reason": "filtered",
  "recipient": "foo@example.com",
  "rhost": "mail.example.com",
  "senderdomain": "sender.example.com",
  "smtpagent": "MTA::Postfix",
  "smtpcommand": "DATA",
  "softbounce": true,
  "subject": "how are you?",
  "timestamp": "2017-03-01T00:00:00Z",
  "timezoneoffset": "+0900",
  "updated_at": "2017-03-01T00:00:00Z",
  "whitelisted": false
}
```
