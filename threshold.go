package main

import (
	"time"
	"sync"
)

type AlarmMap struct {
	lock sync.Mutex
	dataStatus map[string]DataStatus
}

type AlarmRule struct {
	ipSilice []string
	alarmNume int
	alarmTime time.Duration
}

type DataStatus struct {
	domain string
	hostsNum int
	rules []map[string]AlarmRule
}

