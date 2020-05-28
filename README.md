## yf-bot-go

Slack-bot returns [yahoo! finance](https://finance.yahoo.com/) url for symbol.

## How to install

1. `$ git clone https://github.com/msawady/yf-bot-go.git`
2. create new bot and get token from `https://<your-team>.slack.com/apps/A0F7YS25R-bots`
3. `$ cp config.toml.sample config.toml` and paste your token
4. `$ make run`

## Create docker image

* `$ make docker`
  * create docker image `msawady/yf-bot-go:latest` 
* `$ make run-docker` 
  * run created docker image

