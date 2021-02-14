
.PHONY: build stop run

build :
	docker build -f Dockerfile.orig -t my_telegram_rss_bot .

stop :
	docker stop $(docker ps -aqf "name=my_telegram_rss_bot")

rm : stop
	docker rm $(docker ps -aqf "name=my_telegram_rss_bot")

run :
	docker run -d -v `pwd`/config:/app/config --cpus=".7" --memory="750m" --restart always --name my_telegram_rss_bot my_telegram_rss_bot
