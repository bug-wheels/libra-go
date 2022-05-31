package monitor

import (
	"fmt"
	mapset "github.com/deckarep/golang-set/v2"
	"github.com/fatih/color"
	"github.com/gosuri/uitable"
	"libra-go/config"
	"libra-go/discovery"
	"libra-go/util"
	"log"
)

func MonitoringService(services []config.ServiceConfig) {
	util.CallClear()
	log.Println("开始检测服务健康状态")
	table := uitable.New()
	table.MaxColWidth = 100
	table.AddRow("Name", "Status", "Reason")
	for _, service := range services {
		status, reason := monitorService(service)
		table.AddRow(service.Name, util.Any(status, color.GreenString("正常"), color.RedString("警告")), reason)
	}
	fmt.Println(table)
	log.Println("结束检测服务健康状态: 最终结果为")
}

func monitorService(service config.ServiceConfig) (bool, string) {
	if service.Alarm.Check {
		discoveryClient := discovery.RegistrationCenters[service.RegistrationCenter]
		instances, _ := discoveryClient.GetInstances(service.ServiceId)
		if len(instances) < service.Alarm.MinSize {
			return false, fmt.Sprintf("服务数量不对，告警最小数量应为 %d, 现在为 %d", service.Alarm.MinSize, len(instances))
		}
		if len(service.Alarm.MustHosts) > 0 {
			required := mapset.NewSet[string]()
			for _, host := range service.Alarm.MustHosts {
				required.Add(host)
			}
			originHostSet := required.String()
			for _, instance := range instances {
				required.Remove(instance.GetHost())
			}
			if len(required.ToSlice()) > 0 {
				return false, fmt.Sprintf("MustHosts error，%s 不存在，%s 必须存在", required.String(), originHostSet)
			}
		}
	}
	return true, ""
}
