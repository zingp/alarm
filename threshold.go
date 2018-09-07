package main

import (
	"sync"
	"time"
)

type AlarmMap struct {
	lock       sync.Mutex
	dataStatus map[string]*DataStatus
}

type AlarmRule struct {
	ipSilice  []string
	alarmNum  int
	alarmTime time.Duration
}

type DataStatus struct {
	domain   string
	hostsNum int
	rules    []map[string]*AlarmRule
}

var alarmMap = &AlarmMap{}

func NewAlarmMap(c *Yaml) error {
	alarmRule := &AlarmRule{
		alarmNum: 0,
	}

	var slice []map[string]*AlarmRule
	for k := range c.Rules {
		m := make(map[string]*AlarmRule)
		m[k] = alarmRule
		slice = append(slice, m)
	}
	dataStat := &DataStatus{
		domain:   c.Domain,
		hostsNum: len(c.Hosts),
		rules:    slice,
	}
	m := make(map[string]*DataStatus)
	m[c.Domain] = dataStat
	alarmMap = &AlarmMap{
		dataStatus: m,
	}
	return nil
}
