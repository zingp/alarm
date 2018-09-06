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
	if len(procMap) != 0 {
		title := "Tcloud proc relaod"
		maillistStr := strings.Join(c.Maillist, ";")
		// 这里需要补充，并完善报警接口
		alarmMailObj := &alarmMail{
			api:      appConf.mailApi,
			frName:   appConf.frName,
			frAddr:   appConf.frAddr,
			maillist: maillistStr,
			title:    title,
		}

		alarmChan <- alarmMailObj
	}

	itemMap := analysis(data, "items")
	if len(itemMap) == 0 {
		return
	}
	for k, v := range itemMap {
		intV, err := strconv.Atoi(v)
		if err != nil {
			log.Printf("strconv string to int error:%v", err)
			continue
		}
		ok := judge(c, k, intV)
		if ok {
			// 告警结构体；
			// 添加ip等操作
		}
	}
}

func judge(c *Yaml, k string, v int) bool {
	value, ok := c.Rules[k]
	if !ok {
		return false
	}
	switch {
	case value.Sign == ">":
		if v >= value.Condition {
			return true
		}
		return false
	case value.Sign == "<":
		if v <= value.Condition {
			return true
		}
		return false
	}
	return false
}
