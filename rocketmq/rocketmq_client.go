package rocketmq

import (
	"encoding/json"
	"errors"
	mapset "github.com/deckarep/golang-set/v2"
	"github.com/oliveagle/jsonpath"
	"libra-go/config"
	"libra-go/util"
)

type DefaultConsumerInstance struct {
	DiffTotal   int64
	ClientInfos []string
}

type RocketMQClient interface {
	GetInstances(topic string, subscriptionGroup string) (*DefaultConsumerInstance, error)
}

type RocketConsoleClient struct {
	host string
}

func (r RocketConsoleClient) GetInstances(topic string, subscriptionGroup string) (*DefaultConsumerInstance, error) {
	respBody, err := util.HttpGet(r.host + "/topic/queryConsumerByTopic.query?topic=" + topic)
	if err != nil {
		return nil, err
	}

	var jsonData interface{}
	json.Unmarshal([]byte(respBody), &jsonData)
	status, err := jsonpath.JsonPathLookup(jsonData, "$.status")
	if err != nil {
		return nil, errors.New("解析jsonpath失败")
	}
	if status.(float64) != 0 {
		errMsg, _ := jsonpath.JsonPathLookup(jsonData, "$.errMsg")
		return nil, errors.New(errMsg.(string))
	}
	clientInfo, _ := jsonpath.JsonPathLookup(jsonData, "$.data."+subscriptionGroup+".queueStatInfoList.clientInfo")
	diffTotal, _ := jsonpath.JsonPathLookup(jsonData, "$.data."+subscriptionGroup+".diffTotal")

	required := mapset.NewSet[string]()
	for _, host := range clientInfo.([]interface{}) {
		required.Add(host.(string))
	}

	return &DefaultConsumerInstance{
		DiffTotal:   int64(diffTotal.(float64)),
		ClientInfos: required.ToSlice(),
	}, nil
}

func NewRocketConsoleClient(host string) (*RocketConsoleClient, error) {
	if len(host) < 3 {
		return nil, errors.New("check host")
	}
	return &RocketConsoleClient{host: host}, nil
}

var RocketMQClients = make(map[string]RocketMQClient)

func Setup(rocketMQConfig []config.RocketMQClientConfig) {
	for name := range RocketMQClients {
		delete(RocketMQClients, name)
	}
	if len(rocketMQConfig) == 0 {
		return
	}

	for _, rcConfig := range rocketMQConfig {
		client, err := NewRocketConsoleClient(rcConfig.Host)
		if err != nil {
			panic(err)
		}
		RocketMQClients[rcConfig.Name] = client
	}
}
