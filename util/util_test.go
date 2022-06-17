package util

import (
	"encoding/json"
	"fmt"
	mapset "github.com/deckarep/golang-set/v2"
	"github.com/oliveagle/jsonpath"
	"log"
	"strings"
	"testing"
)

func TestHttpGet(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "TestHttpGet",
			args: args{
				url: "http://1231231:8200/topic/queryConsumerByTopic.query?topic=ab",
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := HttpGet(tt.args.url); strings.EqualFold(got, tt.want) {
				t.Errorf("HttpGet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJSONPath(t *testing.T) {
	var jsonData interface{}
	jsonStr := "{\"status\":1,\"errMsg\":\"哈哈哈哈哈哈\"}"
	json.Unmarshal([]byte(jsonStr), &jsonData)
	status, err := jsonpath.JsonPathLookup(jsonData, "$.status")
	if err != nil {
		log.Panic("解析jsonpath失败")
	}
	if status.(float64) != 0 {
		errMsg, _ := jsonpath.JsonPathLookup(jsonData, "$.errMsg")
		log.Panic(errMsg)
	}

	lookup, err := jsonpath.JsonPathLookup(jsonData, "$.data.asd.queueStatInfoList.clientInfo")
	if err != nil {
		log.Panic("解析jsonpath失败")
	}
	fmt.Println(lookup)

	required := mapset.NewSet[string]()
	for _, host := range lookup.([]interface{}) {
		required.Add(host.(string))
	}
	originHostSet := required.String()
	fmt.Println(originHostSet)
}
