package main

import (
	"fmt"
	"github.com/gungoren/rss-telegram-bot-go/pkg/database"
	"github.com/mmcdole/gofeed"
	tb "gopkg.in/tucnak/telebot.v2"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	token             = os.Getenv("TOKEN")
	chat_id           = os.Getenv("CHATID")
	delay             = os.Getenv("DELAY")
	chatId  tb.ChatID = 0
	bot     *tb.Bot
)

func init() {
	database.Setup()
}

func main() {
	b, err := tb.NewBot(tb.Settings{
		Token:  token,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		log.Fatalf("Telegram bot error : %v", err)
		return
	}

	bot = b

	c, err := strconv.Atoi(chat_id)
	if err != nil {
		log.Fatalf("Invalid chat id, %s", chat_id)
	}
	chatId = tb.ChatID(c)

	bot.Handle("/test", cmdTest())
	bot.Handle("/help", cmdRssHelp())

	bot.Handle("/add", cmdRssAdd())
	bot.Handle("/list", cmdRssList())
	bot.Handle("/remove", cmdRssRemove())

	bot.Handle("/add_ban", cmdRssAddBan())
	bot.Handle("/list_ban", cmdRssListBan())
	bot.Handle("/remove_ban", cmdRssRemoveBan())

	monitor := NewMonitor()
	go func(m *Monitor) {
		d, _ := strconv.Atoi(delay)
		for range time.Tick(time.Duration(d) * time.Second) {
			m.rssMonitor()
		}
	}(monitor)

	bot.Start()
}

type Feed struct {
	link string
	last string
}

func listFeeds() map[string]Feed {
	db := database.GetDB()
	rows, err := db.Query(`SELECT * FROM rss`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	feeds := map[string]Feed{}
	for rows.Next() {
		var name string
		var link string
		var last string
		err = rows.Scan(&name, &link, &last)
		if err != nil {
			log.Fatal(err)
		}
		feeds[name] = Feed{link: link, last: last}
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	return feeds
}

func listBans() []string {
	db := database.GetDB()
	rows, err := db.Query(`SELECT * FROM banned_word`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var result []string
	for rows.Next() {
		var word string
		err = rows.Scan(&word)
		if err != nil {
			log.Fatal(err)
		}
		result = append(result, word)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	return result
}

func cmdRssRemoveBan() func(m *tb.Message) {
	return func(m *tb.Message) {
		word := m.Payload

		db := database.GetDB()
		stmt, err := db.Prepare("DELETE FROM banned_word WHERE value = ?")
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()
		if _, err = stmt.Exec(word); err != nil {
			_, _ = bot.Send(m.Chat, "ERROR: DB delete ban error %v", err)
			return
		}

		_, _ = bot.Send(m.Chat, fmt.Sprintf("Removed: %s\n", word))
	}
}

func cmdRssListBan() func(m *tb.Message) {
	return func(m *tb.Message) {
		bans := listBans()
		if len(bans) == 0 {
			_, _ = bot.Send(m.Chat, "Database empty")
		} else {

			for _, ban := range bans {
				msg := fmt.Sprintf("Word: %s \n", ban)
				_, _ = bot.Send(m.Chat, msg)
			}
		}
	}
}

func cmdRssAddBan() func(m *tb.Message) {
	return func(m *tb.Message) {
		if !m.Private() {
			return
		}

		ban := m.Payload

		if len(ban) == 0 {
			_, _ = bot.Send(m.Chat, "ERROR: ban(word) is empty")
			return
		}

		db := database.GetDB()
		stmt, err := db.Prepare("INSERT INTO banned_word('value') VALUES(?)")
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()
		if _, err = stmt.Exec(ban); err != nil {
			_, _ = bot.Send(m.Chat, "ERROR: DB save ban error %v", err)
			return
		}

		_, _ = bot.Send(m.Chat, fmt.Sprintf("added \nBanned word: %s", ban))
	}
}

func cmdRssRemove() func(m *tb.Message) {
	return func(m *tb.Message) {

		params := strings.Split(m.Payload, " ")

		if len(params) != 1 {
			_, _ = bot.Send(m.Chat, "ERROR: The format needs to be: /remove title")
			return
		}

		title := params[0]

		db := database.GetDB()
		stmt, err := db.Prepare("DELETE FROM rss WHERE name = ?")
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()
		if _, err = stmt.Exec(title); err != nil {
			_, _ = bot.Send(m.Chat, "ERROR: DB delete error %v", err)
			return
		}

		_, _ = bot.Send(m.Chat, fmt.Sprintf("Removed: %s\n", title))
	}
}

func cmdRssList() func(m *tb.Message) {
	return func(m *tb.Message) {
		feeds := listFeeds()
		if len(feeds) == 0 {
			if _, err := bot.Send(m.Chat, "Database empty"); err != nil {
				fmt.Printf("Error %v", err)
			}
		} else {
			for name, feed := range feeds {
				msg := fmt.Sprintf("Title: %s \nrss url: %s \nlast checked article: %s", name, feed.link, feed.last)
				_, _ = bot.Send(m.Chat, msg)
			}
		}
	}
}

func cmdTest() func(m *tb.Message) {
	return func(m *tb.Message) {
		url := "https://www.reddit.com/r/funny/new/.rss"
		if f, err := gofeed.NewParser().ParseURL(url); err != nil {
			_, _ = bot.Send(m.Chat, fmt.Sprintf("ERROR : %v", err))
		} else {
			_, _ = bot.Send(m.Chat, f.Items[0].Link)
		}
	}
}

func cmdRssHelp() func(m *tb.Message) {
	return func(m *tb.Message) {
		help := fmt.Sprintf(
			"RSS to Telegram bot"+
				"\n\nAfter successfully adding a RSS link, the bot starts fetching the feed every "+
				delay+" seconds. (This can be set)"+
				"\n\nTitles are used to easily manage RSS feeds and need to contain only one word"+
				"\n\ncommands:"+
				"\n/help Posts this help message"+
				"\n/add title http://www(.)RSS-URL(.)com"+
				"\n/remove !Title! removes the RSS link"+
				"\n/list Lists all the titles and the RSS links from the DB"+
				"\n/add_ban word"+
				"\n/list_ban Lists all the banned words"+
				"\n/remove_ban word Delete word from the banned words"+
				"\n/test Inbuilt command that fetches a post from Reddits RSS."+
				"\n\nThe current chatId is: %d ", m.Chat.ID)

		_, _ = bot.Send(m.Chat, help)
	}
}

func cmdRssAdd() func(m *tb.Message) {
	return func(m *tb.Message) {
		fp := gofeed.NewParser()

		if !m.Private() {
			return
		}

		params := strings.Split(m.Payload, " ")

		if len(params) != 2 {
			_, _ = bot.Send(m.Chat, "ERROR: The format needs to be: /add title http://www.URL.com")
			return
		}

		title := params[0]
		url := params[1]

		feed, err := fp.ParseURL(url)
		if err != nil {
			_, _ = bot.Send(m.Chat, "ERROR: The link does not seem to be a RSS feed or is not supported")
			return
		}

		db := database.GetDB()
		stmt, err := db.Prepare("INSERT INTO rss('name','link','last') VALUES(?,?,?)")
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()
		if _, err = stmt.Exec(title, url, feed.Items[0].Link); err != nil {
			_, _ = bot.Send(m.Chat, "ERROR: DB save error %v", err)
			return
		}

		_, _ = bot.Send(m.Chat, fmt.Sprintf("added \nTITLE: %s\nRSS: %s", title, url))
	}
}
