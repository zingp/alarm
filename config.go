package main
import (
	"path"
	"io/ioutil"
	yaml "gopkg.in/yaml.v2"
	"log"
	"strings"
)

type Rule struct {
	Minute int `yaml:"minute"`
	Freq  int `yaml:"freq"`
	Ring  int `yaml:"ring"`
	Sign string `yaml:"sign"`
	Condition  int `yaml:"condition"`
}

type Yaml struct {
	Domain string  `yaml:"domain"` 
	Hosts []string `yaml:"hosts"`
	Rules map[string]Rule `yaml:"rules"`
	Maillist []string `yaml:"maillist"`
}

func getYamlList(dir string)(fileSlice []string, err error){
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Printf("read dir error:%v", err)
		return
	}
	for _,file := range files {
		if file.IsDir() {
			continue
		} 
		if strings.HasSuffix(file.Name(), ".yaml") {
			fileSlice = append(fileSlice, file.Name())
			continue
		}
	}
	return
}


func loadYaml(file string) (conf *Yaml, err error) {
	conf = new(Yaml)
    yamlFile, err := ioutil.ReadFile(file)
    if err != nil {
		log.Printf("parse yaml file error:%v ", err)
		return
    }
    err = yaml.Unmarshal(yamlFile, &conf)
    if err != nil {
		log.Printf("Unmarshal: %v", err)
		return
    }
	return 
}

func initConfigMap(dir string)(err error){
	fileList, err := getYamlList(dir)
	if err != nil {
		log.Printf("get yaml list error:%v", err)
		return
	}
	log.Printf("get yaml list sucess, files:%v", fileList)

	for _, f := range fileList {
		file := path.Join(dir, f)

		conf, err := loadYaml(file)
		if err != nil {
			log.Printf("init config map load yaml error:%v", err)
			continue
		}

		configMap[conf.Domain] = conf
	}
	return
}

