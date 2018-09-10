package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

/*
2018/09/03 15:36:24.145 [I] items [upgrade=2,]
2018/09/03 15:46:28.149 [I] items [upgrade=1,]
2018/09/03 16:01:05.874 [I] items [upgrade=1,]
2018/09/03 16:32:02.983 [I] items [upgrade=1,]
2018/09/03 16:37:37.998 [I] items [upgrade=1,]

proc [name=%s,cont=%s]
*/

var alarmChan = make(chan *alarmMail, 10)

func requestGet(url string) (result string, err error) {
	res, err := http.Get(url)
	if err != nil {
		fmt.Println(" http get error:", err)
		return
	}
	if res.Status != string(200) {
		log.Printf("http get not 200.code=%s, detail=%v", res.Status, res.Body)
		return
	}
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("read data error:", err)
		return
	}
	result = string(data)
	return
}

func getPart(item string) (part string) {
	partStr := `%s(?P<part1>.+) %s \[(?P<part2>.+)\]`
	timeMinStr := time.Now().Format("2006/01/02 15:04")
	part = fmt.Sprintf(partStr, timeMinStr, item)
	return
}

func analysis(s string, p string) (m map[string]string) {
	m = make(map[string]string, 5)
	lineSlice := strings.Split(s, "\n")
	part := getPart(p)
	for _, line := range lineSlice {
		reg := regexp.MustCompile(part)
		param := reg.FindStringSubmatch(line)
		if len(param) > 0 {
			m = stringToMap(param[2])
			return
		}
	}
	return
}

func stringToMap(s string) (m map[string]string) {
	m = make(map[string]string, 5)
	newStr := strings.TrimSpace(s)
	kvSlice := strings.Split(newStr, ",")
	for _, v := range kvSlice {
		trimV := strings.TrimSpace(v)
		if len(trimV) == 0 {
			continue
		}
		subKvSlice := strings.Split(trimV, "=")

		if len(subKvSlice) == 2 {
			mapK := strings.TrimSpace(subKvSlice[0])
			mapV := strings.TrimSpace(subKvSlice[1])
			m[mapK] = mapV
		}
	}
	return
}

func getAgentDate() {
	for _, v := range configMap {
		if len(v.Hosts) == 0 {
			continue
		}

		for _, ip := range v.Hosts {
			// 应该使用协程
			handleData(ip, v)
			// go ------

		}
	}
}

func handleData(ip string, c *Yaml) {
	maillistStr := strings.Join(c.Maillist, ";")

	url := fmt.Sprintf(appConf.reqUrl, ip)
	data, err := requestGet(url)
	if err != nil {
		log.Printf("http get error:%v", err)
		return
	}
	if len(data) == 0 {
		return
	}

	procMap := analysis(data, "proc")
	itemMap := analysis(data, "items")

	if len(procMap) != 0 {
		title := "Tcloud proc relaod"
		proc := procMap["name"]
		cont := procMap["cont"]
		rule := fmt.Sprintf("过去1分钟有进程[%s]被重启", proc)
		body := fmt.Sprintf(htmlBody, rule, c.Domain, ip, cont)
		// 这里需要补充，并完善报警接口
		alarmMailObj := &alarmMail{
			frName:   appConf.frName,
			frAddr:   appConf.frAddr,
			maillist: maillistStr,
			title:    title,
			body:     body,
			mode:     "html",
		}
		alarmChan <- alarmMailObj
	}

	if len(itemMap) == 0 {
		return
	}
	for k, v := range itemMap {
		intV, err := strconv.Atoi(v)
		if err != nil {
			log.Printf("strconv string to int error:%v", err)
			continue
		}
		r, ok := judge(c, k, intV)
		// 如果该ip的某个规则满足告警条件
		v1, ok1 := ipStatus.ipMap[ip][r]
		if ok {
			if !ok1 {
				continue
			}

			if v1 >= (c.Rules[k].Freq -1) {
				// alarm
				title := fmt.Sprintf("日志监控告警%s", r)
				body := fmt.Sprintf(htmlBody, r, c.Domain, ip, v)
				alarmMailObj := &alarmMail{
					frName:   appConf.frName,
					frAddr:   appConf.frAddr,
					maillist: maillistStr,
					title:    title,
					body:     body,
					mode:     "html",
				}
				alarmChan <- alarmMailObj
			}
			ipStatus.Add(ip, r)
			continue
		}
		if v1 >= (c.Rules[k].Freq -1) {
			title := fmt.Sprintf("日志监控告警恢复%s", r)
				body := fmt.Sprintf(htmlBody, r, c.Domain, ip, v)
				alarmMailObj := &alarmMail{
					frName:   appConf.frName,
					frAddr:   appConf.frAddr,
					maillist: maillistStr,
					title:    title,
					body:     body,
					mode:     "html",
				}
				alarmChan <- alarmMailObj
				ipStatus.Add(ip, r)
			// 告警恢复
		}
		err = ipStatus.Zero(ip, r)
		if err != nil {
			continue
		}
	}
}

func judge(c *Yaml, k string, v int) (r string, b bool) {
	value, ok := c.Rules[k]
	r = fmt.Sprintf("%s%s%d", k, value.Sign, value.Condition)
	b =  false
	if !ok {
		return 
	}
	switch {
	case value.Sign == ">":
		if v >= value.Condition {
			b = true
		}
	case value.Sign == "<":
		if v <= value.Condition {
			b = true
		}
	}
	return 
}
