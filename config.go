package main


type Rule struct {
	Minute int `yaml:"minute"`
	Freq  int `yaml:"freq"`
	Condition  string `yaml:"condition"`
}

type Yaml struct {
	Domain string  `yaml:"domain"` 
	Hosts []string `yaml:"hosts"`
	Rules map[string]Rule `yaml:"rules"`
}

