package command

import (
	"fmt"
	"log"
	"strconv"

	qqbotapi "github.com/catsworld/qq-bot-api"
)

//GetDetail deal with resp and return detail map, rowCount, msgSender
func GetDetail(upd *qqbotapi.Update, bot *qqbotapi.BotAPI, mapResult map[string]interface{}) ([]interface{}, int) {
	if mapResult["result"] == "failed" {
		bot.NewMessage(upd.Message.Chat.ID, upd.Message.Chat.Type).Text(mapResult["message"].(string)).Send()
		return nil, 0
	}
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
