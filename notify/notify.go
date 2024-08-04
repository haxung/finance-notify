package notify

import (
	"encoding/json"
	"strings"
	"time"

	"finance-notify/client"
	"finance-notify/common"

	"gopkg.in/gomail.v2"
)

func Notify(body string) error {
	c := common.EnvConf.Email

	// 准备邮件内容
	m := gomail.NewMessage()
	m.SetHeader("From", c.From)
	m.SetHeader("To", c.To...)
	m.SetHeader("Subject", "Finance Notify")
	m.SetBody("text/plain", body)
	m.SetHeader("Date", time.Now().Local().UTC().Format(time.RFC1123Z))

	// 发送邮件
	dialer := gomail.NewDialer(c.Server, c.Port, c.Username, c.Password)
	return dialer.DialAndSend(m)
}

func CoinResp(crs []*client.CoinResp) string {
	b := strings.Builder{}
	b.WriteString("=================================\n")
	for i, resp := range crs {
		bytes, _ := json.Marshal(resp)
		b.Write(bytes)
		if i != len(crs)-1 {
			b.WriteString("\n-------------------------------\n")
		}
	}

	b.WriteString("=================================\n")
	return b.String()
}
