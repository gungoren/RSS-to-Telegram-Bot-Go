![RSSTT](https://github.com/gungoren/RSS-to-Telegram-Bot-Go/raw/master/img/rsstt.png)

# RSS to Telegram bot with golang

> Note: This repo is migrated from [Python version](https://github.com/gungoren/RSS-to-Telegram-Bot/)

A self-hosted telegram golang bot that dumps posts from RSS feeds to a telegram chat. This script was created because all the third party services were unreliable.

![Image of help menu](https://github.com/gungoren/RSS-to-Telegram-Bot-Go/raw/master/img/telegram.png)

### Docker

For the docker image go to: https://hub.docker.com/r/mehgungoren/rss.to.telegram/

### Installation

Golang 1.14

```sh
go mod download
```

A telegram bot is needed that the script will connect to. https://botsfortelegram.com/project/the-bot-father/
Running the script and typing in /help will reveal the current chatId, this needs to be set also in the script

1. Clone the script
2. Replace your chatID and Token on the top of the script.
3. Edit the delay. (seconds)
4. Save and run
5. Use the telegram commands to manage feeds

Warning! Without chatID the bot wont be able to send automated messages and will only be able to respond to messages.

# Usage

send /help to the bot to get this message:

> RSS to Telegram bot
>
> After successfully adding a RSS link, the bot starts fetching the feed every 60 seconds. (This can be set)
> Titles are used to easily manage RSS feeds and need to contain only one word
>
> commands:
>
> **/help** Shows this text
>
> **/add** title http://www(.)URL(.)com
>
> **/remove** !Title! removes the RSS link
>
> **/list** Lists all the titles and the RSS links from the DB
>
> **/add_ban** word add word to ban list
>
> **/remove_ban** word removes word from ban word list
>
> **/list_ban** Lists all the banned words from the DB
>
> **/test** Inbuilt command that fetches a post from Reddits RSS.
>
> The current chatId is: 20416xxxx

# Known issues

If the bot is set to for example 5 minutes and one feed manages to get 2 new posts before the bot can check. Only the newest post will show up on telegram.

# Docker

```
docker create \
  --name=rss.to.telegram \
  -e DELAY=180 \
  -e TOKEN=<TOKEN> \
  -e CHATID=<CHATID> \
  -v /path/to/host/config:/config \
  --restart unless-stopped \
  mehgungoren/rss.to.telegram
```
