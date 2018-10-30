package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"
	"os"
	tb "gopkg.in/tucnak/telebot.v2"
)

const functionalitiesText = "Use /currentreport for the latest topic." +
	"\nUse /reportcount for the number of reports." +
	"\n\nUse /subscribe will notify you when there is a new report." +
	"\nUse /unsubscribe will stop notifying you when there is a new report."

var bot *tb.Bot

func main() {
	TELEGRAM_BOT_API_KEY := os.Getenv("TELEGRAM_BOT_API_KEY")
	LISTENER_FILE := os.Getenv("LISTENER_FILE")

	if LISTENER_FILE == "" && TELEGRAM_BOT_API_KEY == ""{
		fmt.Println("Environment variables not set!\nTELEGRAM_BOT_API_KEY & LISTENER_FILE")
		os.Exit(0)
	}

	bot = initBot()
	Init()
	checkForNewData()
	go bot.Start()

	ticker := time.NewTicker(time.Hour * 24)

	func() {
		for t := range ticker.C {
			fmt.Println("Checking data at : " + t.String())
			checkForNewData()
		}
	}()
}

func checkForNewData() {
	fmt.Println("Starting new data check!")
	stats := GetReportStats()

	fmt.Println("Data received : ")
	prettified, _ := prettyify(stats)
	fmt.Println(prettified)

	isDataNew := UpdateStats(stats)

	if isDataNew {
		NotifySubscribers(bot)
	}
}

func prettyify(data interface{}) (string, error) {
	bytes, e := json.MarshalIndent(data, "", "  ")

	if e != nil {
		return "", e
	}

	return string(bytes), nil
}

func initBot() *tb.Bot {
	bot, err := tb.NewBot(tb.Settings{
		Token:  os.Getenv("TELEGRAM_BOT_API_KEY"),
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		log.Fatal(err)
		return nil
	}

	bot.Handle("/currentreport", func(m *tb.Message) {
		bot.Send(m.Sender, GetLatestTopic())
	})

	bot.Handle("/reportcount", func(m *tb.Message) {
		bot.Send(m.Sender, strconv.Itoa(GetTotalTopicCount()))
	})

	bot.Handle("/subscribe", func(m *tb.Message) {
		bot.Send(m.Sender, "🎈🎈When there is a new report, I will notify you!🎉🎉")
		AddUser(New(m.Sender.ID, m.Sender.FirstName))
	})

	bot.Handle("/unsubscribe", func(m *tb.Message) {
		deleted := RemoveUser(New(m.Sender.ID, m.Sender.FirstName))
		if deleted {
			bot.Send(m.Sender, "😟😟😟\n Thanks for listening!\n\n I'm unsubscribing you now... 😐")
		} else {
			bot.Send(m.Sender, "🤨🤨🤨\n\n...😔\n\nYou were never even subscribed... 😞😞😞")
		}

	})

	inlineBtn := tb.InlineButton{
		Unique: "subscribe",
		Text:   "🎉 Subscribe! 🎉",
	}

	inlineKeys := [][]tb.InlineButton{
		{inlineBtn},
	}

	bot.Handle(&inlineBtn, func(c *tb.Callback) {
		// on inline button pressed (callback!)
		AddUser(New(c.Sender.ID, c.Sender.FirstName))
		bot.Respond(c, &tb.CallbackResponse{Text: "🎉🎉 You are now subscribed! 🎉🎉"})
	})

	bot.Handle("/start", func(m *tb.Message) {
		bot.Send(m.Sender, functionalitiesText, &tb.ReplyMarkup{
			InlineKeyboard: inlineKeys,
		})
	})

	return bot
}
