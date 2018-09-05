package main

import (
	"strings"
	"net/http"
	"io/ioutil"
	"fmt"
	"time"
	"regexp"
)

/*
2018/09/03 15:36:24.145 [I] items [upgrade=2,]
2018/09/03 15:46:28.149 [I] items [upgrade=1,]
2018/09/03 16:01:05.874 [I] items [upgrade=1,]
2018/09/03 16:32:02.983 [I] items [upgrade=1,]
2018/09/03 16:37:37.998 [I] items [upgrade=1,]

proc [name=%s,cont=%s]
*/

func requestGet(url string)(result string, err error){
	res, err := http.Get(url)
	if err != nil {
		fmt.Println(" http get error:", err)
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

func getPart(item string)(part string) {
	partStr := `%s(?P<part1>.+) %s \[(?P<part2>.+)\]`
	timeMinStr := time.Now().Format("2006/01/02 15:04")
	part = fmt.Sprintf(partStr, timeMinStr, item)
	return 
}

func analysis(s string, p string)(m map[string]string){
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

func stringToMap(s string) (m map[string]string){
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