package sendcloud_sms_go

const (
	SMS          = 0
	MMS          = 1
	INTERNAT_SMS = 2
	VOICE        = 3
	QR_CODE      = 4
	YX           = 5
)

type SendSmsTemplateArgs struct {
	TemplateId  int
	LabelId     int
	MsgType     int
	Phone       string
	Vars        string
	SendRequestId string
	Tag string
}


type SendSmsResult struct {
	Result     bool        `json:"result"`
	StatusCode int         `json:"statusCode"`
	Message    string      `json:"message"`
	Info       interface{} `json:"info"`
}

