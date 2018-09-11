package main

import (
	"log"
)

type AppConf struct {
	reqURL  string
	mailAPI string
	frName  string
	frAddr  string
	mode    string
}

var configMap map[string]*Yaml
var appConf AppConf

func main() {
	// url := "http://127.0.0.1:9090/items"
	// HttpGet(url)
	dir := "./conf"
	fileList, err := getYamlList(dir)
	if err != nil {
		log.Printf("get yaml list error:%v", err)
		return
	}
	log.Printf("get yaml list sucess, files:%v", fileList)
	// fmt.Println("conf", conf)
	// fmt.Println("domain", conf.Domain)
	// fmt.Println("hosts", conf.Hosts)
	// for k, v :=  range conf.Hosts {
	// 	fmt.Println(k ,v)
	// }
	// fmt.Println("rules", conf.Rules)
	// for k, v := range conf.Rules {
	// 	fmt.Printf("k=%s  v=%v\n", k, v)
	// }
	// s := `2018/09/04 17:59:24.145 [I] items [upgrade=2,]
	// 2018/09/04 17:59:28.149 [I] items [upgrade=1,]
	// 2018/09/04 17:59:05.874 [I] items [upgrade=1,]
	// 2018/09/04 17:59:02.983 [I] items [upgrade=1,]
	// 2018/09/04 18:06:37.998 [I] items [upgrade=1,]`
	// analysis(s, "proc")
	sendMail()
}