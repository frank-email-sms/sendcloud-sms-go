package sendcloud_sms_go

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

type SmsClient struct {
	SmsUser    string
	SmsKey     string
	SmsBasePath string
}

const (
	smsBasePath = "https://api.sendcloud.net"
	sendSmsTemplatePath  =  "/smsapi/send"
	sendSmsCodePath   =  "/smsapi/sendCode"
)

func NewSmsClient(smsUser string, smsKey string) (*SmsClient, error) {
	switch {
	case len(smsUser) == 0:
		return nil,errors.New("smsUser cannot be empty")
	case len(smsKey) == 0:
		return nil,errors.New("smsKey cannot be empty")
	}
	return &SmsClient{
        SmsUser:    smsUser,
        SmsKey:     smsKey,
        SmsBasePath: smsBasePath,
    }, nil
}


func (client *SmsClient) SendSmsTemplate(args *SendSmsTemplateArgs) (*SendSmsResult, error) {
	if err := client.validateConfig(); err != nil {
		return nil,fmt.Errorf("failed to send message: %w", err)
	}
	if err := validateSendSmsTemplate(args); err != nil {
		return nil,fmt.Errorf("failed to send message: %w", err)
	}

	params, err := client.prepareSendSmsTemplateParams(args)
	if err != nil {
		return nil,fmt.Errorf("failed to send message: %w", err)
	}

	signature := client.calculateSignature(params)
	params.Set("signature", signature)

	resp, err := http.PostForm(client.SmsBasePath + sendSmsTemplatePath, params)
	if err != nil {
		return nil,fmt.Errorf("failed to send HTTP POST request: %w", err)
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			fmt.Println(err)
			return
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil,fmt.Errorf("HTTP POST request failed with status code %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil,fmt.Errorf("failed to read response body: %w", err)
	}

	var responseData SendSmsResult
	if err := json.Unmarshal(body, &responseData); err != nil {
		return nil,fmt.Errorf("failed to unmarshal response body: %w", err)
	}
	result := &responseData
	if !result.Result {
		return result,fmt.Errorf("SMS sending failed: %s", result.Message)
	}
	return result,nil

}

func (client *SmsClient)calculateSignature(params url.Values) string {
	// 创建一个新的url.Values，以便我们可以修改它而不影响原始params
	sortedParams := url.Values{}

	// 将除了smsKey和signature之外的所有参数复制到新的map中
	for k, v := range params {
		if k != "smsKey" && k != "signature" {
			sortedParams[k] = v
		}
	}

	// 创建一个键的切片，以便我们可以按字母顺序对它们进行排序
	keys := make([]string, 0, len(sortedParams))
	for k := range sortedParams {
		keys = append(keys, k)
	}

	// 对键进行排序
	sort.Strings(keys)

	// 构建param_str
	var paramStr string
	for _, k := range keys {
		paramStr += k + "=" + sortedParams.Get(k) + "&"
	}

	// 去掉paramStr末尾的'&'
	if len(paramStr) > 0 {
		paramStr = paramStr[:len(paramStr)-1]
	}

	// 构建sign_str
	signStr := client.SmsKey + "&" + paramStr + "&" + client.SmsKey

	// 计算sha256值
	hasher := sha256.New()
	hasher.Write([]byte(signStr))
	sha256Bytes := hasher.Sum(nil)

	// 将MD5字节转换为十六进制字符串
	signature := hex.EncodeToString(sha256Bytes)

	return signature
}

func (client *SmsClient) prepareSendSmsTemplateParams(args *SendSmsTemplateArgs) (url.Values, error) {
	params := url.Values{}
	params.Set("smsUser", client.SmsUser)
	params.Set("msgType", strconv.Itoa(args.MsgType))
	params.Set("phone", args.Phone)
	params.Set("templateId", strconv.Itoa(args.TemplateId))
	params.Set("timestamp", strconv.FormatInt(time.Now().UnixNano()/int64(time.Millisecond), 10))
	if len(args.Vars) > 0 {
		varsJSON, err := json.Marshal(args.Vars)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal vars: %v", err)
		}
		params.Set("vars", string(varsJSON))
	}

	return params, nil
}


func (client *SmsClient)validateConfig() error {
	if len(client.SmsBasePath) == 0 {
		client.SmsBasePath = smsBasePath
    }
	switch {
	case len(client.SmsUser) == 0:
		return errors.New("smsUser cannot be empty")
	case len(client.SmsKey) == 0:
		return errors.New("smsKey cannot be empty")
	}
	return nil
}

func isValidMsgType(msgType int) bool {
	return msgType == SMS ||
		msgType == MMS ||
		msgType == INTERNAT_SMS ||
		msgType == VOICE ||
		msgType == QR_CODE ||
		msgType == YX
}

// ValidatePhoneNumbers 校验phone参数中的手机号
func ValidatePhoneNumbers(phone string) error {
	// 分割手机号
	phoneNumbers := strings.Split(phone, ",")

	// 校验手机号数量是否超过2000个
	if len(phoneNumbers) > 2000 {
		return errors.New("the number of mobile phone numbers exceeds the maximum limit of 2,000")
	}

	// 遍历每个手机号并校验格式
	for _, number := range phoneNumbers {
		// 去除可能存在的空白字符
		trimmedNumber := strings.TrimSpace(number)
		if trimmedNumber == "" {
			return errors.New("phone number can not be blank")
		}
	}
	// 所有手机号都校验通过
	return nil
}

func validateSendSmsTemplate(args *SendSmsTemplateArgs) error {
	switch {
	case args.TemplateId == 0:
		return errors.New("templateId cannot be zero")
	case !isValidMsgType(args.MsgType):
		return errors.New("msgType cannot be negative")
	case len(args.Phone) == 0:
		return errors.New("phone cannot be empty")
	}
	if err := ValidatePhoneNumbers(args.Phone); err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}
	return nil
}