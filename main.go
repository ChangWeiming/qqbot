package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	constant "qqbot/src/constant"
	"strconv"
	"time"

	qqbotapi "github.com/catsworld/qq-bot-api"
)

var mp map[string]string

//ConnectWithServer connects with api server
func ConnectWithServer(upd *qqbotapi.Update, bot *qqbotapi.BotAPI, mod constant.MOD, params string) (map[string]interface{}, error) {
	php := ""
	switch mod {
	case constant.FAILED:
		php = "today_overdue_ddl.php"
	case constant.TODAY:
		php = "today_upcoming_ddl.php"
	case constant.RECENT:
		php = "recent_updated_ddl.php"
	case constant.UserFAILED:
		php = "user_overdue_ddl.php"
	case constant.UserRECENT:
		php = "user_upcoming_ddl.php"
	}

	var resp *http.Response
	var err error
	if params == "" {
		resp, err = http.PostForm(constant.ServerName+php,
			url.Values{"auth_key": {"513106c051f94528f1d386926aa65e1a"}})
	} else {
		resp, err = http.PostForm(constant.ServerName+php,
			url.Values{"auth_key": {"513106c051f94528f1d386926aa65e1a"}, "username": {params}})
	}

	if err != nil {
		bot.NewMessage(upd.Message.Chat.ID, upd.Message.Chat.Type).Text("JZDKServer returns " + fmt.Sprint(err)).Send()
		return nil, err
	}

	if resp.Status != "200 OK" {
		tmp, _ := ioutil.ReadAll(resp.Body)
		bot.NewMessage(upd.Message.Chat.ID, upd.Message.Chat.Type).Text("JZDKServer returns " + resp.Status + string(tmp)).Send()
		return nil, nil
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var mapResult map[string]interface{}
	//log.Println(">>>>>>>>>>>>>>>>>")
	//log.Println(string(body))
	if err := json.Unmarshal(body, &mapResult); err != nil {
		return nil, err
	}

	log.Print(mapResult)
	return mapResult, nil
}

//GetDetail deal with resp and return detail map, rowCount, msgSender
func GetDetail(upd *qqbotapi.Update, bot *qqbotapi.BotAPI, mapResult map[string]interface{}) ([]interface{}, int) {
	if mapResult["result"] == "failed" {
		bot.NewMessage(upd.Message.Chat.ID, upd.Message.Chat.Type).Text(mapResult["message"].(string)).Send()
		return nil, 0
	} else {
		//rowCount, err := strconv.Atoi(mapResult["row_count"].(string))
		//log.Print(mapResult)
		rowCount, err := strconv.Atoi(fmt.Sprint(mapResult["row_count"]))
		if err != nil {
			log.Println(err)
			return nil, 0
		}

		detailInterface := mapResult["detail"].([]interface{})
		return detailInterface, rowCount
	}

}

//TodayDDL send today ddl to QQ chat
func TodayDDL(upd *qqbotapi.Update, bot *qqbotapi.BotAPI) error {
	mapResult, err := ConnectWithServer(upd, bot, constant.TODAY, "")

	if err != nil {
		return err
	}

	msg := bot.NewMessage(upd.Message.Chat.ID, upd.Message.Chat.Type)
	var msgSender *qqbotapi.FlatSender = msg.FlatSender
	detailInterface, rowCount := GetDetail(upd, bot, mapResult)
	lstname := ""

	if rowCount != 0 {
		msgSender = msgSender.Text("哈哈，今天ddl截止了")
	} else {
		msgSender = msgSender.Text("今天暂无ddl，没想到吧，快来加点吧")
	}

	for i := 0; i < rowCount; i++ {
		detail := detailInterface[i].(map[string]interface{})
		name := detail["name"].(string)
		if name != lstname {
			if val, ok := mp[name]; ok {
				msgSender = msgSender.NewLine().At(val).NewLine()
			}
			lstname = name
			msgSender = msgSender.Text("姓名：" + name).NewLine()
		}
		//fmt.Print(reflect.TypeOf(mapResult["progress"]))
		msgSender = msgSender.Text(detail["description"].(string)).
			Text(" | ").Text(detail["completed_parts"].(string) + "/" + detail["total_parts"].(string)).
			Text(fmt.Sprintf("[%s]", detail["progress"].(string)))
		if i != rowCount-1 {
			msgSender = msgSender.NewLine()
		}
	}
	msgSender.Send()
	return nil
}

//RecentDDL send today ddl to QQ chat
func RecentDDL(upd *qqbotapi.Update, bot *qqbotapi.BotAPI) error {
	mapResult, err := ConnectWithServer(upd, bot, constant.RECENT, "")

	if err != nil {
		return err
	}

	msg := bot.NewMessage(upd.Message.Chat.ID, upd.Message.Chat.Type)
	var msgSender *qqbotapi.FlatSender = msg.FlatSender
	msgSender = msgSender.Text("最近DDL：")
	detailInterface, rowCount := GetDetail(upd, bot, mapResult)
	lstname := ""

	for i := 0; i < rowCount; i++ {
		detail := detailInterface[i].(map[string]interface{})
		name := detail["name"].(string)
		if name != lstname {
			if val, ok := mp[name]; ok {
				msgSender = msgSender.NewLine().At(val).NewLine()
			}
			lstname = name
			msgSender = msgSender.Text("姓名：" + name).NewLine()
		}
		//fmt.Print(reflect.TypeOf(mapResult["progress"]))
		msgSender = msgSender.Text(detail["description"].(string)).
			Text(" | ").Text(detail["completed_parts"].(string) + "/" + detail["total_parts"].(string)).
			Text(fmt.Sprintf("[%s]", detail["progress"].(string))).
			Text(" | ").
			Text(detail["status"].(string))
		if i != rowCount-1 {
			msgSender = msgSender.NewLine()
		}
	}
	msgSender.Send()
	return nil
}

//OverdueDDL send ddls that fail to chat
func OverdueDDL(upd *qqbotapi.Update, bot *qqbotapi.BotAPI) error {
	mapResult, err := ConnectWithServer(upd, bot, constant.FAILED, "")
	if err != nil {
		log.Print(err)
		return err
	}
	detailInter, rowCount := GetDetail(upd, bot, mapResult)

	msg := bot.NewMessage(upd.Message.Chat.ID, upd.Message.Chat.Type)
	var msgSender *qqbotapi.FlatSender = msg.FlatSender
	if rowCount != 0 {
		msgSender = msgSender.Text("哈哈！昨天失败了吧")
	} else {
		msgSender = msgSender.Text("今天干的还行，通报表扬")
	}

	lstname := ""
	for i := 0; i < rowCount; i++ {
		detail := detailInter[i].(map[string]interface{})
		name := detail["name"].(string)
		if name != lstname {
			if val, ok := mp[name]; ok {
				msgSender = msgSender.NewLine().At(val).NewLine()
			}
			lstname = name
			msgSender = msgSender.Text("姓名：" + name).NewLine()
		}
		//fmt.Print(reflect.TypeOf(mapResult["progress"]))
		msgSender = msgSender.Text(detail["description"].(string)).
			Text(" | ").Text(detail["completed_parts"].(string) + "/" + detail["total_parts"].(string)).
			Text(fmt.Sprintf("[%s]", detail["progress"].(string)))
		if i != rowCount-1 {
			msgSender = msgSender.NewLine()
		}
	}
	msgSender.Send()
	return nil
}

//UserOverdueDDL send failed ddl of certain user (mysb)
func UserOverdueDDL(upd *qqbotapi.Update, bot *qqbotapi.BotAPI, QQ string) error {
	mapResult, err := ConnectWithServer(upd, bot, constant.UserFAILED, constant.QQToName[QQ])
	if err != nil {
		log.Print(err)
		return err
	}
	detailInter, rowCount := GetDetail(upd, bot, mapResult)

	msg := bot.NewMessage(upd.Message.Chat.ID, upd.Message.Chat.Type)
	var msgSender *qqbotapi.FlatSender = msg.FlatSender
	if rowCount != 0 {
		msgSender = msgSender.Text("哈哈！昨天失败了吧")
	} else {
		msgSender = msgSender.Text("今天干的还行，你今天不xx了")
	}

	msgSender = msgSender.NewLine().At(QQ).NewLine()
	for i := 0; i < rowCount; i++ {
		detail := detailInter[i].(map[string]interface{})
		//fmt.Print(reflect.TypeOf(mapResult["progress"]))
		msgSender = msgSender.Text(detail["description"].(string)).
			Text(" | ").Text(detail["completed_parts"].(string) + "/" + detail["total_parts"].(string)).
			Text(fmt.Sprintf("[%s]", detail["progress"].(string)))
		if i != rowCount-1 {
			msgSender = msgSender.NewLine()
		}
	}
	msgSender.Send()
	return nil
}

//UserRecentDDL returns recent ddls of certain users (myddl)
func UserRecentDDL(upd *qqbotapi.Update, bot *qqbotapi.BotAPI, QQ string) error {
	mapResult, err := ConnectWithServer(upd, bot, constant.UserRECENT, constant.QQToName[QQ])
	if err != nil {
		log.Print(err)
		return err
	}

	detailInter, rowCount := GetDetail(upd, bot, mapResult)

	msg := bot.NewMessage(upd.Message.Chat.ID, upd.Message.Chat.Type)
	var msgSender *qqbotapi.FlatSender = msg.FlatSender
	if rowCount != 0 {
		msgSender = msgSender.Text("gkd！别天天颓了！！！")
	} else {
		msgSender = msgSender.Text("快去给自己整活去！")
	}

	msgSender = msgSender.NewLine().At(QQ).NewLine()
	for i := 0; i < rowCount; i++ {
		detail := detailInter[i].(map[string]interface{})
		//fmt.Print(reflect.TypeOf(mapResult["progress"]))
		ddlTime, _ := time.Parse("2006-01-02", detail["deadline"].(string))
		ddlTime = ddlTime.UTC().Local().Add(16 * time.Hour)
		gap := ddlTime.Sub(time.Now()).Hours()

		var dateMsg string
		if gap > 0 {
			dateMsg = fmt.Sprint((int)(gap/24)+1) + "天"
		} else {
			dateMsg = fmt.Sprint((int)(gap/24)-1) + "天"
		}

		msgSender = msgSender.Text(detail["description"].(string)).
			Text(" | ").Text(dateMsg).
			Text(" | ").Text(detail["completed_parts"].(string) + "/" + detail["total_parts"].(string)).
			Text(fmt.Sprintf("[%s]", detail["progress"].(string)))
		if i != rowCount-1 {
			msgSender = msgSender.NewLine()
		}
	}
	msgSender.Send()
	return nil
}

//TimerSender sets clock to send ddls regularly
func TimerSender(bot *qqbotapi.BotAPI) {
	var hvSend = -1
	var p = 2
	var timepoint []int = []int{6, 12, 14, 18, 22}
	var upd *qqbotapi.Update = &qqbotapi.Update{
		Message: &qqbotapi.Message{
			Chat: &qqbotapi.Chat{
				ID:   123456789,
				Type: "group",
			}}}
	for {
		hr := time.Now().Hour()
		if hr == timepoint[p] && hr != hvSend {
			hvSend = hr
			for _, v := range mp {
				UserRecentDDL(upd, bot, v)
			}
			p = (p + 1) % 5
		}
	}
}

func main() {
	bot, err := qqbotapi.NewBotAPI("", "http://0.0.0.0:5700", "")
	if err != nil {
		log.Fatal(err)
	}
	mp = constant.NameToQQ

	bot.Debug = false

	u := qqbotapi.NewWebhook("/")
	u.PreloadUserInfo = true
	updates := bot.ListenForWebhook(u)
	go http.ListenAndServe("localhost:2333", nil)
	go TimerSender(bot)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		fmt.Printf("[%s] %s", update.Message.From.String(), update.Message.Text)

		retStr := ""

		switch update.Message.Text {
		case "!recentlists", "！recentlists":
			if err := RecentDDL(&update, bot); err != nil {
				log.Println(err)
			}
		case "!todaylists", "！todaylists":
			if err := TodayDDL(&update, bot); err != nil {
				log.Println(err)
			}
		case "!ping", "！ping":
			retStr = "pong"
			bot.NewMessage(update.Message.Chat.ID, update.Message.Chat.Type).Text(retStr).Send()
		case "!failed", "！failed":
			if err := OverdueDDL(&update, bot); err != nil {
				log.Println(err)
			}
		case "!myddl", "！myddl":
			if err := UserRecentDDL(&update, bot, fmt.Sprint(update.Message.From.ID)); err != nil {
				log.Println(err)
			}
		case "!mysb", "！mysb":
			if err := UserOverdueDDL(&update, bot, fmt.Sprint(update.Message.From.ID)); err != nil {
				log.Println(err)
			}
		}

		if len(update.Message.Text) >= 8 && update.Message.Text[0:4] == "\\add" {
			var opt, orin, tar string
			if len(mp) > 10 {
				bot.NewMessage(update.Message.Chat.ID, update.Message.Chat.Type).At(fmt.Sprint(update.Message.From.ID)).Text("别瞎JB添加了").Send()
				continue
			}
			fmt.Sscanf(update.Message.Text, "%s %s %s", &opt, &orin, &tar)
			mp[orin] = tar
			bot.NewMessage(update.Message.Chat.ID, update.Message.Chat.Type).At(mp[orin]).NewLine().Text("binded with name " + orin).Send()
		}

	}
}
