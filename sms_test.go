package sendcloud_sms_go

import "testing"

func TestSendSmsTemplate(t *testing.T) {
	client, err := NewSmsClient("**", "**")
    if err != nil {
        t.Error(err)
    }
    result, err := client.SendSmsTemplate(&SendSmsTemplateArgs{
        TemplateId: 1,
        LabelId:    1,
        Phone:      "13800138000,13800138001",
        MsgType:    SMS,
    })
    if err != nil {
        t.Error(err)
    }
    t.Log(result)
}


