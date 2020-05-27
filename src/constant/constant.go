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
	//AddDDL means MOD add ones ddl
	AddDDL
	//FinishDDL means MOD finish ddl
	FinishDDL
	//DeleteDDL means MOD finish ddl
	DeleteDDL
)

//ServerName JZDK server
const ServerName string = "https://test.example.com/api/"

//AlertGroup group that is to be sent ddls regularly
const AlertGroup int64 = 123456789

//Timepoint time that sends dlls
var Timepoint = [...]int{10, 14, 18, 22}

//NameToQQ mapping JZDK ID to QQ number
var NameToQQ map[string]string = map[string]string{
	"test": "123456789",
}

//QQToName mapping QQ number to JZDK ID
var QQToName map[string]string = map[string]string{
	"123456789": "test",
}
