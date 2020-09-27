package command

import (
	"fmt"
	"log"
	"qqbot/src/constant"

	qqbotapi "github.com/catsworld/qq-bot-api"
)

//TodayDDL send today ddl to QQ chat
func TodayDDL(upd *qqbotapi.Update, bot *qqbotapi.BotAPI) error {
	mapResult, err := ConnectWithServer(upd, bot, constant.TODAY, nil)

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
			if val, ok := constant.NameToQQ[name]; ok {
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
	mapResult, err := ConnectWithServer(upd, bot, constant.RECENT, nil)

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
			if val, ok := constant.NameToQQ[name]; ok {
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
	mapResult, err := ConnectWithServer(upd, bot, constant.FAILED, nil)
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
			if val, ok := constant.NameToQQ[name]; ok {
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

//UsageHelp sends usage of the bot
func UsageHelp(upd *qqbotapi.Update, bot *qqbotapi.BotAPI) {
	msg := bot.NewMessage(upd.Message.Chat.ID, upd.Message.Chat.Type)
	var msgSender *qqbotapi.FlatSender = msg.FlatSender

	msgSender.Text("!help 帮助").NewLine().
		Text("!myddl 查看自己的ddl").NewLine().
		Text("!mysb 查看自己的失败ddl").NewLine().
		Text("!recentlists 所有人最近ddl").NewLine().
		Text("!todaylists 所有人今天ddl").NewLine().
		Text("!failed 所有人最近失败").NewLine().
		Text("!ping 判断机器人挂了没").NewLine().
		Text("!add dates description [priority] [totoal parts]添加ddl 方框表示可选项（可倒序缺失），空格分开，dates两种格式一个是日期20200101或者1表示1天后即当日凌晨之前").NewLine().
		Text("!check ddlID 打卡编号为ddlID的ddl").NewLine().
		Text("!del ddlID 删除编号为ddlID的ddl").NewLine().
		Text("!mkfail ddlID 使编号为ddlID的ddl失败").
		Send()
}
