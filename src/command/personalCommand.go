package command

import (
	"fmt"
	"log"
	"net/url"
	"qqbot/src/constant"
	"strconv"
	"time"

	qqbotapi "github.com/catsworld/qq-bot-api"
)

//UserOverdueDDL send failed ddl of certain user (mysb)
func UserOverdueDDL(upd *qqbotapi.Update, bot *qqbotapi.BotAPI, QQ string) error {
	params := url.Values{"username": {constant.QQToName[QQ]}}
	mapResult, err := ConnectWithServer(upd, bot, constant.UserFAILED, params)
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
	params := url.Values{"username": {constant.QQToName[QQ]}}
	mapResult, err := ConnectWithServer(upd, bot, constant.UserRECENT, params)
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

	var upd *qqbotapi.Update = &qqbotapi.Update{
		Message: &qqbotapi.Message{
			Chat: &qqbotapi.Chat{
				ID:   constant.AlertGroup,
				Type: "group",
			}}}
	for {
		hr := time.Now().Hour()
		if hr == constant.Timepoint[p] && hr != hvSend {
			hvSend = hr
			for _, v := range constant.NameToQQ {
				UserRecentDDL(upd, bot, v)
			}
			p = (p + 1) % 4
		}
	}
}

//AddUserDDL user can add ddl using add command
func AddUserDDL(upd *qqbotapi.Update, bot *qqbotapi.BotAPI, QQ string) {
	const negVal = -23333
	var opt, date, desc string
	var priority, totalParts int = negVal, negVal

	fmt.Sscanf(upd.Message.Text, "%s %s %s %d %d", &opt, &date, &desc, &priority, &totalParts)

	msg := bot.NewMessage(upd.Message.Chat.ID, upd.Message.Chat.Type)
	var msgSender *qqbotapi.FlatSender = msg.FlatSender

	if opt != "!add" && opt != "！add" {
		msgSender.At(QQ).Text("命令后加空格！").Send()
		return
	}
	if len(date) <= 3 {
		if d, err := strconv.Atoi(date); err != nil {
			msgSender.At(QQ).Text("date格式不对！！").Send()
			return
		} else {
			d--
			d *= 24
			date = time.Now().Add(time.Hour * time.Duration(d)).Format("20060102")
		}
	}
	params := url.Values{"name": {constant.QQToName[QQ]}, "deadline": {date}, "description": {desc}}

	if priority != negVal {
		params["priority"] = []string{fmt.Sprint(priority)}
	} else {
		params["priority"] = []string{"0"}
	}

	if totalParts != negVal {
		params["totalparts"] = []string{fmt.Sprint(totalParts)}
	} else {
		params["totalparts"] = []string{"1"}
	}

	if resp, err := ConnectWithServer(upd, bot, constant.AddDDL, params); err != nil {
		msgSender.At(QQ).Text(fmt.Sprint(err)).Send()
	} else {
		if resp["result"].(string) == "success" {
			msgSender.At(QQ).Text("成了！！").Send()
			return
		}
		msgSender.At(QQ).Text(resp["message"].(string)).Send()
		return
	}
}
