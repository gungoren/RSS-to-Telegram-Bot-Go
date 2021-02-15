package main

import (
	"fmt"
	"github.com/gungoren/rss-telegram-bot-go/pkg/database"
	"github.com/mmcdole/gofeed"
	"log"
	"regexp"
	"strconv"
	"strings"
)

type Monitor struct{}

func NewMonitor() *Monitor {
	return &Monitor{}
}

var (
	minHourlyPrice = 35.0
	minFixPrice    = 400
)

func (m *Monitor) rssMonitor() {

	feeds := listFeeds()
	bans := listBans()

	for name, feed := range feeds {
		fp := gofeed.NewParser()
		rss, err := fp.ParseURL(feed.link)
		if err != nil {
			log.Println(err)
			return
		}
		lastLink := feed.last
		for _, item := range rss.Items {
			if item.Link != feed.last {
				sendMessageToChat(name, item, bans)
			} else {
				break
			}
		}

		if len(rss.Items) > 0 {
			lastLink = rss.Items[0].Link
		}

		db := database.GetDB()
		stmt, err := db.Prepare("UPDATE rss SET last = ? WHERE name = ? AND link = ?")
		if err != nil {
			log.Fatal(err)
		}
		if _, err = stmt.Exec(lastLink, name, feed.link); err != nil {
			log.Printf("ERROR: DB save ban error %v", err)
			return
		}
		_ = stmt.Close()
	}
}

func sendMessageToChat(name string, item *gofeed.Item, bans []string) {
	detail := item.Content

	sendMessage := true
	budget := ""

	if strings.Contains(detail, "Hourly Range") {
		sendMessage, budget = getHourlyPrice(detail)
	} else if strings.Contains(detail, "Budget</b>") {
		sendMessage, budget = checkEntryBudget(detail)
	}

	prefix := ""
	if sendMessage && strings.Contains(detail, "Country") && !checkBlockedCountry(detail) {
		prefix = "⚠️⚠️⚠️⚠️"
	}

	if !sendMessage {
		return
	}

	if isMessageAlreadySend(item.Link) {
		return
	}

	if checkEntryContainsBannedWord(bans, strings.ToLower(detail)) {
		return
	}

	saveMessageSend(item.Link)
	_, _ = bot.Send(chatId, fmt.Sprintf("%s %s %s %s", prefix, strings.ReplaceAll(item.Link, "?source=rss", ""), name, budget))
}

func saveMessageSend(link string) {
	db := database.GetDB()
	stmt, err := db.Prepare("INSERT INTO messages_send('link') VALUES(?)")
	if err != nil {
		log.Fatal(err)
	}
	if _, err = stmt.Exec(link); err != nil {
		log.Printf("Link save to history failed")
		return
	}
	_ = stmt.Close()
}

func checkEntryContainsBannedWord(bans []string, entryDetail string) bool {
	for _, ban := range bans {
		if strings.Contains(entryDetail, ban) {
			return true
		}
	}
	return false
}

func isMessageAlreadySend(link string) bool {
	db := database.GetDB()
	stmt, err := db.Prepare(`SELECT * FROM messages_send WHERE link = ?`)
	if err != nil {
		log.Fatal(err)
	}

	if results, err := stmt.Query(link); err != nil {
		return false
	} else {
		defer results.Close()
		return results.Next()
	}

	return false
}

func checkBlockedCountry(detail string) bool {
	reg := regexp.MustCompile(`Country.*?: (.*)\n`)
	country := reg.FindStringSubmatch(detail)[1]
	country = strings.ToLower(country)
	return !strings.Contains(country, "india")
}

func checkEntryBudget(detail string) (bool, string) {
	reg := regexp.MustCompile(`Budget.*?: .?([0-9,]+)`)
	matches := reg.FindStringSubmatch(detail)
	if len(matches) > 0 {
		budget := strings.ReplaceAll(matches[1], ",", "")
		b, _ := strconv.Atoi(budget)
		if b >= minFixPrice {
			return true, budget
		}
		return false, ""
	}
	return true, ""
}

func getHourlyPrice(detail string) (bool, string) {
	reg := regexp.MustCompile(`Hourly Range.*?: (.*)\n`)
	price := reg.FindStringSubmatch(detail)[1]
	price = strings.ReplaceAll(price, "$", "")
	prices := strings.Split(price, "-")
	if len(prices) == 1 {
		p, _ := strconv.ParseFloat(prices[0], 64)
		return p >= minHourlyPrice, price
	}
	if len(prices) == 2 {
		l, _ := strconv.ParseFloat(prices[0], 64)
		h, _ := strconv.ParseFloat(prices[1], 64)
		return l >= minHourlyPrice || h >= minHourlyPrice, price
	}
	return true, price
}
