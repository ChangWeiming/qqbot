package command

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"qqbot/src/constant"

	qqbotapi "github.com/catsworld/qq-bot-api"
)

//ConnectWithServer connects with api server
func ConnectWithServer(upd *qqbotapi.Update, bot *qqbotapi.BotAPI, mod constant.MOD, params url.Values) (map[string]interface{}, error) {
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
	case constant.AddDDL:
		php = "add_ddl.php"
	}

	var resp *http.Response
	var err error
	if params == nil {
		resp, err = http.PostForm(constant.ServerName+php,
			url.Values{"auth_key": {"513106c051f94528f1d386926aa65e1a"}})
	} else {
		params["auth_key"] = []string{"513106c051f94528f1d386926aa65e1a"}
		resp, err = http.PostForm(constant.ServerName+php,
			params)
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
