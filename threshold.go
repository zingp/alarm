package main

import (
	"fmt"
	"sync"
	"time"
	"errors"
)

type AlarmMap struct {
	lock       sync.Mutex
	dataStatus map[string]*DataStatus
}

type AlarmRule struct {
	ipSilice  [][]string
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

type IPStatus struct {
	lock sync.Mutex
	ipMap map[string]map[string]int
	/*
	ip : {
		redis>2 : 1,
		update>3:2,
	}
	*/
}

var ipStatus = &IPStatus{}

// NewIPStatus初始化ipStatus
func NewIPStatus(c *Yaml){
	for _, ip := range c.Hosts {
		for k, v := range c.Rules {
			key := fmt.Sprintf("%s%s%d", k, v.Sign, v.Condition)
			m := make(map[string]int,5)
			m[key] = 0
			ipStatus.ipMap[ip] = m
		}
	}
}

// Add 给定ip， 规则，对应值+1
func (ips *IPStatus)Add(ip string, rule string) error {
	ips.lock.Lock()
	ips.lock.Unlock()
	if v, ok := ips.ipMap[ip]; ok {
		v[rule]++
		return nil
	}
	return errors.New("ip 不存在")
}

// Sub 给定ip,规则 对应值-1
func (ips *IPStatus)Sub(ip string, rule string) error {
	ips.lock.Lock()
	ips.lock.Unlock()
	if v, ok := ips.ipMap[ip]; ok {
		v[rule]++
		return nil
	}
	return errors.New("ip 不存在")
}
