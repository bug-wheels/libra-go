package rocketmq

import (
	"log"
	"testing"
)

func TestNewRocketConsoleClient(t *testing.T) {
	client, _ := NewRocketConsoleClient("http://123123:8200")
	instances, err := client.GetInstances("sdfasdf", "asdfasdfa")
	if err != nil {
		log.Panicf("%s", err)
	}
	log.Println(instances)
}
