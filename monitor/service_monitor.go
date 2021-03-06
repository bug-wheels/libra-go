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

func Monitoring(libra config.LibraConfig) []MonitoringInfo {
	var monitoringInfo []MonitoringInfo
	util.CallClear()
	log.Println("开始检测服务健康状态")
	log.Println("===================< 服务状态查询 >========================")
	service := monitoringService(libra.Services)
	monitoringInfo = append(monitoringInfo, service)
	log.Println("===================< MQ 状态查询 >========================")
	mq := monitoringRocketMQ(libra.Consumers)
	monitoringInfo = append(monitoringInfo, mq)
	log.Println("结束检测服务健康状态")
	return monitoringInfo
}

type MonitoringInfo struct {
	Name  string
	Title []string
	Data  [][]string
}

func monitoringService(services []config.ServiceConfig) MonitoringInfo {
	monitoringInfo := MonitoringInfo{
		Name:  "Service",
		Title: []string{"Name", "Status", "Reason"},
		Data:  [][]string{},
	}
	table := uitable.New()
	table.MaxColWidth = 100
	table.AddRow("Name", "Status", "Reason")
	for _, service := range services {
		status, reason, instancesCount := monitorService(service)
		table.AddRow(service.Name, util.Any(status, color.GreenString("正常("+strconv.Itoa(instancesCount)+")"), color.RedString("警告("+strconv.Itoa(instancesCount)+")")), reason)
		monitoringInfo.Data = append(monitoringInfo.Data, []string{service.Name, util.Any(status, "正常("+strconv.Itoa(instancesCount)+")", "警告("+strconv.Itoa(instancesCount)+")"), reason})

	}
	fmt.Println(table)
	return monitoringInfo
}

func monitoringRocketMQ(consumers []config.ConsumerConfig) MonitoringInfo {
	monitoringInfo := MonitoringInfo{
		Name:  "MQ",
		Title: []string{"Name", "Status", "Delay", "Reason"},
		Data:  [][]string{},
	}
	table := uitable.New()
	table.MaxColWidth = 100
	table.AddRow("Name", "Status", "Delay", "Reason")

	for _, consumer := range consumers {
		status, reason, instancesCount, delay := monitorConsumer(consumer)
		table.AddRow(consumer.Name, util.Any(status, color.GreenString("正常("+strconv.Itoa(instancesCount)+")"), color.RedString("警告("+strconv.Itoa(instancesCount)+")")), delay, reason)
		monitoringInfo.Data = append(monitoringInfo.Data, []string{consumer.Name, util.Any(status, "正常("+strconv.Itoa(instancesCount)+")", "警告("+strconv.Itoa(instancesCount)+")"), strconv.FormatInt(delay, 10), reason})
	}
	fmt.Println(table)
	return monitoringInfo
}

func monitorConsumer(consumer config.ConsumerConfig) (bool, string, int, int64) {
	instancesCount := 0
	delay := int64(-1)
	if consumer.Alarm.Check {
		rocketMQClient := rocketmq.RocketMQClients[consumer.MQ]
		instances, err := rocketMQClient.GetInstances(consumer.Topic, consumer.SubscriptionGroup)
		if err != nil {
			return false, err.Error(), 0, -1
		}
		instancesCount = len(instances.ClientInfos)
		delay = instances.DiffTotal
		if instancesCount < consumer.Alarm.MinSize {
			return false, fmt.Sprintf("服务数量不对，告警最小数量应为 %d, 现在为 %d", consumer.Alarm.MinSize, instancesCount), instancesCount, instances.DiffTotal
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
				return false, fmt.Sprintf("MustHosts error，%s 不存在，%s 必须存在", required.String(), originHostSet), instancesCount, instances.DiffTotal
			}
		}
	}
	return true, "", instancesCount, delay

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
