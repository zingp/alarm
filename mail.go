package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"io/ioutil"
)

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
	api := "http://mail.e.sogou/sendMail?uid=cp_notice@sogou-inc.com&fr_name=%s&fr_addr=%s&maillist=%s&title=%s&body=%s&mode=%s&attname=&attbody="
	frName := "Liuyouyuan"
	frAddr := "dt_op@sogou-inc.com"
	maillist := "liuyouyuan@sogou-inc.com"
	title := "title"
	body := "content"
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
	api := "http://mail.e.sogou/sendMail"
	m := make(map[string]string)
	m["uid"] = "cp_notice@sogou-inc.com"
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
