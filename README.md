[![Build Status](https://travis-ci.com/elegant-bro/amqp-recipient.svg?branch=master)](https://travis-ci.com/elegant-bro/amqp-recipient)

# Useful go components for amqp consumers
This module based on popular https://github.com/streadway/amqp library.

## Preface
Working on several microservices that communicate through rabbitmq I found there are common tasks not related
to the domain I need to solve. Most often I need consumers to be idempotent. Logs are important too. Probably
if message handling fails you want retry it later or even republish to another exchange/queue. So I made handlers
to solve this tasks. 

## Install
`go get https://github.com/elegant-bro/amqp-recipient@v{version}`

## Usage
Suppose you need consumer that able to listen events from queue `user.registered` and send Welcome email using some
mail service's api.

Let's look at [basic example](./_examples/basic.go).

