package main

import (
	"fmt"
	"log"
	"net/http"
	"qqbot/src/command"

	qqbotapi "github.com/catsworld/qq-bot-api"
)

func main() {
	bot, err := qqbotapi.NewBotAPI("", "http://0.0.0.0:5700", "")
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = false

	u := qqbotapi.NewWebhook("/")
	u.PreloadUserInfo = true
	updates := bot.ListenForWebhook(u)
	go http.ListenAndServe("localhost:2333", nil)
	go command.TimerSender(bot)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		fmt.Printf("[%s] %s", update.Message.From.String(), update.Message.Text)

		retStr := ""

		switch update.Message.Text {
		case "!recentlists", "！recentlists":
			if err := command.RecentDDL(&update, bot); err != nil {
				log.Println(err)
			}
		case "!todaylists", "！todaylists":
			if err := command.TodayDDL(&update, bot); err != nil {
				log.Println(err)
			}
		case "!ping", "！ping":
			retStr = "pong"
			bot.NewMessage(update.Message.Chat.ID, update.Message.Chat.Type).Text(retStr).Send()
		case "!failed", "！failed":
			if err := command.OverdueDDL(&update, bot); err != nil {
				log.Println(err)
			}
		case "!myddl", "！myddl":
			if err := command.UserRecentDDL(&update, bot, fmt.Sprint(update.Message.From.ID)); err != nil {
				log.Println(err)
			}
		case "!mysb", "！mysb":
			if err := command.UserOverdueDDL(&update, bot, fmt.Sprint(update.Message.From.ID)); err != nil {
				log.Println(err)
			}
		case "!help", "！help":
			command.UsageHelp(&update, bot)
		}
		/*
			if len(update.Message.Text) >= 8 && update.Message.Text[0:4] == "\\add" {
				var opt, orin, tar string
				if len(mp) > 10 {
					bot.NewMessage(update.Message.Chat.ID, update.Message.Chat.Type).At(fmt.Sprint(update.Message.From.ID)).Text("别瞎JB添加了").Send()
					continue
				}
				fmt.Sscanf(update.Message.Text, "%s %s %s", &opt, &orin, &tar)
				mp[orin] = tar
				bot.NewMessage(update.Message.Chat.ID, update.Message.Chat.Type).At(mp[orin]).NewLine().Text("binded with name " + orin).Send()
			} else
		*/
		if len(update.Message.Text) >= 8 && len(update.Message.Text) <= 100 {
			if update.Message.Text[0:4] == "!add" || update.Message.Text[0:6] == "！add" {
				command.AddUserDDL(&update, bot, fmt.Sprint(update.Message.From.ID))
			}
		}

	}
}
