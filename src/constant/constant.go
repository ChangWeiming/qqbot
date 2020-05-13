package constant

//MOD means the module of the bot
type MOD int

const (
	_ MOD = iota
	//TODAY means MOD ddl finishes today
	TODAY
	//RECENT means MOD ddl updates recently
	RECENT
	//FAILED means MOD ddl that fails
	FAILED
	//UserRECENT means MOD certain user's ddl recently
	UserRECENT
	//UserFAILED means MOD certain user's failed ddl today
	UserFAILED
)

//ServerName JZDK server
const ServerName string = "https://test.example.com/api/"

//NameToQQ mapping JZDK ID to QQ number
var NameToQQ map[string]string = map[string]string{
	"test": "123456789",
}

//QQToName mapping QQ number to JZDK ID
var QQToName map[string]string = map[string]string{
	"123456789":  "test",
}
