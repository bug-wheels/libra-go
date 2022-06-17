package monitor

import (
	"fmt"
	mapset "github.com/deckarep/golang-set/v2"
	"github.com/fatih/color"
	"github.com/gosuri/uitable"
	"libra-go/config"
	"libra-go/discovery"
	"libra-go/rocketmq"
	"libra-go/util"
	"log"
	"strconv"
)

func Monitoring(libra config.LibraConfig) {
	util.CallClear()
	log.Println("开始检测服务健康状态")
	monitoringService(libra.Services)
	monitoringRocketMQ(libra.Consumers)
	log.Println("结束检测服务健康状态")
}

func monitoringService(services []config.ServiceConfig) {
	table := uitable.New()
	table.MaxColWidth = 100
	table.AddRow("Name", "Status", "Reason")
	for _, service := range services {
		status, reason, instancesCount := monitorService(service)
		table.AddRow(service.Name, util.Any(status, color.GreenString("正常("+strconv.Itoa(instancesCount)+")"), color.RedString("警告("+strconv.Itoa(instancesCount)+")")), reason)
	}
	fmt.Println(table)
}

func monitoringRocketMQ(consumers []config.ConsumerConfig) {
	table := uitable.New()
	table.MaxColWidth = 100
	table.AddRow("Name", "Status", "Reason")
	for _, consumer := range consumers {
		status, reason, instancesCount := monitorConsumer(consumer)
		table.AddRow(consumer.Name, util.Any(status, color.GreenString("正常("+strconv.Itoa(instancesCount)+")"), color.RedString("警告("+strconv.Itoa(instancesCount)+")")), reason)
	}
	fmt.Println(table)
}

func monitorConsumer(consumer config.ConsumerConfig) (bool, string, int) {
	instancesCount := 0
	if consumer.Alarm.Check {
		rocketMQClient := rocketmq.RocketMQClients[consumer.MQ]
		instances, _ := rocketMQClient.GetInstances(consumer.Topic, consumer.SubscriptionGroup)
		instancesCount = len(instances.ClientInfos)
		if instancesCount < consumer.Alarm.MinSize {
			return false, fmt.Sprintf("服务数量不对，告警最小数量应为 %d, 现在为 %d", consumer.Alarm.MinSize, instancesCount), instancesCount
		}
		if len(consumer.Alarm.MustHosts) > 0 {
			required := mapset.NewSet[string]()
			for _, host := range consumer.Alarm.MustHosts {
				required.Add(host)
			}
			originHostSet := required.String()
			for _, instance := range instances.ClientInfos {
				required.Remove(instance)
			}
			if len(required.ToSlice()) > 0 {
				return false, fmt.Sprintf("MustHosts error，%s 不存在，%s 必须存在", required.String(), originHostSet), instancesCount
			}
		}
	}
	return true, "", instancesCount
}

func monitorService(service config.ServiceConfig) (bool, string, int) {
	instancesCount := 0
	if service.Alarm.Check {
		discoveryClient := discovery.RegistrationCenters[service.RegistrationCenter]
		instances, _ := discoveryClient.GetInstances(service.ServiceId)
		instancesCount = len(instances)
		if instancesCount < service.Alarm.MinSize {
			return false, fmt.Sprintf("服务数量不对，告警最小数量应为 %d, 现在为 %d", service.Alarm.MinSize, instancesCount), instancesCount
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
				return false, fmt.Sprintf("MustHosts error，%s 不存在，%s 必须存在", required.String(), originHostSet), instancesCount
			}
		}
	}
	return true, "", instancesCount
}
