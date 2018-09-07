package main

import (
	"errors"
	"fmt"
	"github.com/axgle/mahonia"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

var emailURL = "http://mail.portal.sogou/portal/tools/send_mail.php"
var htmlBody = `<!DOCTYPE html>
				<html lang="en">
				<head>
					<meta charset="UTF-8">
				</head>
				<body>
					<p>Title:%s</p>
					<p>Domain:%s</p>
					<p>IP:%s</p>
					<p>Detail:%s</p>
				</body>
				</html>`

type alarmMail struct {
	frName   string
	frAddr   string
	maillist string
	title    string
	body     string
	mode     string
	attname  string
	attbody  string
}

func (a alarmMail) MailSend() (bool, error) {
	if _, err := checkEmail(a.maillist); err != nil {
		return false, err
	}

	enc := mahonia.NewEncoder("GB18030")
	req, _ := url.Parse(emailURL)
	q := req.Query()
	q.Add("fr_name", a.frName)
	q.Add("fr_addr", a.frAddr)
	q.Add("title", enc.ConvertString(a.title))
	q.Add("body", enc.ConvertString(a.body))
	q.Add("mode", a.mode)
	q.Add("maillist", a.maillist)
	q.Add("attname", a.attname)
	q.Add("attbody", a.attbody)
	req.RawQuery = q.Encode()
	resp, err := http.Get(req.String())
	if err != nil {
		return false, err
	}
	defer resp.Body.Close() // email body为空

	if resp.StatusCode != 200 {
		return false, errors.New("Email Send failed")
	}
	return true, nil
}

func sendMail() {
	/*
		fr_name - 发信人姓名
		fr_addr - 发信人email
		title - 邮件标题
		body - 邮件内容
		mode - 邮件类型，html或txt
		maillist - 收信人邮箱，多个邮箱用";"分隔
		attname - 附件文件名
		attbody - 附件正文*/
	htm := `<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
	</head>
	<body>
		<p>Title:%s</p>
		<p>Domain:%s</p>
		<p>IP:%s</p>
		<p>Detail:%s</p>
	</body>
	</html>`
	frName := "Liuyouyuan"
	frAddr := "dt_op@sogou-inc.com"
	maillist := "liuyouyuan@sogou-inc.com"
	title := "title"
	body := fmt.Sprintf(htm, "Tcloud proc relaod", "abc.com", "10.134.239.239", "cont")

	r, _ := EmailSend(maillist, body, frName, frAddr, title)
	fmt.Println("sendmail:", r)
}

func checkEmail(addr string) (bool, error) {
	addrs := strings.Split(addr, ";")
	for _, v := range addrs {
		if len(v) > 0 {
			matched, err := regexp.MatchString("^[a-zA-Z0-9_-]+@sogou-inc.com$", v)
			if !matched || err != nil {
				return matched, errors.New("邮箱格式错误，请确认！")
			}
		}
	}
	return true, nil
}

// EmailSend send email to users
// maillist: 用逗号分割的邮件名
// name: From邮件名
// emailFrom：From邮件地址
func EmailSend(maillist, content, name, emailFrom, title string) (bool, error) {
	if _, err := checkEmail(maillist); err != nil {
		return false, err
	}
	enc := mahonia.NewEncoder("GB18030")
	req, _ := url.Parse(emailURL)
	q := req.Query()
	q.Add("fr_name", name)
	q.Add("fr_addr", emailFrom)
	q.Add("title", enc.ConvertString(title))
	q.Add("body", enc.ConvertString(content))
	q.Add("mode", "html")
	q.Add("maillist", maillist)
	q.Add("attname", "")
	q.Add("attbody", "")
	req.RawQuery = q.Encode()
	resp, err := http.Get(req.String())
	if err != nil {
		return false, err
	}
	defer resp.Body.Close() // email body为空
	if resp.StatusCode != 200 {
		return false, errors.New("Email Send failed")
	}
	return true, nil
}
