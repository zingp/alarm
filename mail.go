package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"io/ioutil"
)
type alarmMail struct {
	api string
	frName string
	frAddr string
	maillist string
	title string
	body string
	mode string
	attname string
	attbody string
}

func (a *alarmMail)sendMailGet()(resp *http.Response, err error){
	url := fmt.Sprintf(a.api, a.frName, a.frAddr, a.maillist, a.title, a.body, a.mode, a.attname, a.attbody)
	resp, err = http.Get(url)
	if err != nil {
		log.Printf("http error:%v", err)
	}
	return
}

func sendMail() {
	/*
	uid - 申请权限的user_id，请使用sogou-inc邮箱账号
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
	// api := "http://mail.e.sogou/sendMail?uid=cp_notice@sogou-inc.com&fr_name=%s&fr_addr=%s&maillist=%s&title=%s&body=%s&mode=%s&attname=&attbody="
	api := "http://mail.portal.sogou/portal/tools/send_mail.php?uid=pangbowen@sogou-inc.com&fr_name=%s&fr_addr=%s&title=%s&body=%s&mode=%s&maillist=%s&attname=&attbody="
	frName := "Liuyouyuan"
	frAddr := "dt_op@sogou-inc.com"
	maillist := "liuyouyuan@sogou-inc.com"
	title := "title"
	body := fmt.Sprintf(htm, "Tcloud proc relaod", "abc.com","10.134.239.239","cont")
	mode := "html"
	
	url := fmt.Sprintf(api, frName, frAddr,maillist,title, body,mode)
	log.Println(url)
	ret, err := http.Get(url)
	if err != nil {
		log.Printf("http error:%v", err)
	}
	fmt.Println("ret=", ret)
}

func sendMailPost() {
	api := "http://mail.portal.sogou/portal/tools/send_mail.php"
	// api := "http://mail.e.sogou/sendMail"
	m := make(map[string]string)
	// m["uid"] = "cp_notice@sogou-inc.com"
	m["uid"] = "pangbowen@sogou-inc.com"
	m["fr_name"] = "Liuyouyuan"
	m["fr_addr"] = "dt_op@sogou-inc.com"
	m["maillist"] = "liuyouyuan@sogou-inc.com"
	m["title"] = "标题"
	m["body"] = "content"
	m["mode"] = "html"
	m["attname"] = ""
	m["attbody"] = ""

	b, err := json.Marshal(m)
	if err != nil {
		log.Printf("json marshal error:%v", err)
		return
	}
	body := bytes.NewBuffer([]byte(b))
	res,err := http.Post(api, "application/json;charset=utf-8", body)
        if err != nil {
                log.Fatal(err)
                return
        }
        result, err := ioutil.ReadAll(res.Body)
        res.Body.Close()
        if err != nil {
                log.Fatal(err)
                return
        }
        fmt.Printf("prost result::%s", result)
}
