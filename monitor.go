package main

import (
	"fmt"
	"github.com/gungoren/rss-telegram-bot-go/pkg/database"
	"github.com/mmcdole/gofeed"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Monitor struct{}

func NewMonitor() *Monitor {
	return &Monitor{}
}

var (
	minHourlyPrice = 35.0
	minFixPrice    = 400
)

var visitedLinks = map[string]time.Time{}

func (m *Monitor) rssMonitor() {

	feeds := listFeeds()
	bans := listBans()

	fp := gofeed.NewParser()
	for name, feed := range feeds {
		rss, err := fp.ParseURL(feed.link)
		if err != nil {
			log.Println(err)
			continue
		}
		lastLink := feed.last
		for _, item := range rss.Items {
			if item.Link != lastLink {
				sendMessageToChat(name, item, bans)
			} else {
				break
			}
		}

		lastLink = rss.Items[0].Link

		db := database.GetDB()
		stmt, err := db.Prepare("UPDATE rss SET last = ? WHERE name = ? AND link = ?")
		if err != nil {
			log.Fatal(err)
		}
		if _, err = stmt.Exec(lastLink, name, feed.link); err != nil {
			log.Printf("ERROR: DB update lastlink error %v", err)
			continue
		}
		_ = stmt.Close()
	}
}

func (m *Monitor) checkedExpiredVisitedLinks() {
	v := map[string]time.Time{}
	yesterday := time.Now().Add(-24 * time.Hour)
	for link, t := range visitedLinks {
		if t.After(yesterday) {
			v[link] = t
		}
	}
	visitedLinks = v
}

func sendMessageToChat(name string, item *gofeed.Item, bans []string) {
	detail := item.Content

	sendMessage := true
	budget := ""

	if _, ok := visitedLinks[item.Link]; ok {
		return
	}

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

	if checkEntryContainsBannedWord(bans, strings.ToLower(detail)) {
		return
	}

	visitedLinks[item.Link] = time.Now()
	//log.Printf(fmt.Sprintf("%s %s %s %s", prefix, strings.ReplaceAll(item.Link, "?source=rss", ""), name, budget))
	_, _ = bot.Send(chatId, fmt.Sprintf("%s %s %s %s", prefix, strings.ReplaceAll(item.Link, "?source=rss", ""), name, budget))
}

func checkEntryContainsBannedWord(bans []string, entryDetail string) bool {
	for _, ban := range bans {
		if strings.Contains(entryDetail, ban) {
			return true
		}
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
