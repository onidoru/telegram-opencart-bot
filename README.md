# telegram-opencart-bot
[![Build Status](https://travis-ci.com/onidoru/telegram-opencart-bot.svg?branch=master)](https://travis-ci.com/onidoru/telegram-opencart-bot)

Bot for Telegram OpenCart Backend

Project uses standart $GOPATH dep management. Config file ``config.yml`` must be edited accordingly before running the bot.

Compile & Run:

  ```
  $ go get github.com/onidoru/telegram-opencart-bot
  $ cd $GOPATH/src/github.com/onidoru/telegram-opencart-bot
  $ go get ./...
  $ go build main.go
  $ ./main
  ```
  
  Build with Docker:
  ```
  $ sh createImage.sh
  $ docker run telegram-opencart-bot
  ```
  
