package config

type LibraConfig struct {
	RegistrationCenters []RegistrationCenterConfig `yaml:"registration-centers"`
	RocketMQClients     []RocketMQClientConfig     `yaml:"rocket-mq"`
	Services            []ServiceConfig            `yaml:"services"`
	Consumers           []ConsumerConfig           `yaml:"consumers"`
}

type ServiceConfig struct {
	Name               string      `yaml:"name"`
	ServiceId          string      `yaml:"service-id"`
	RegistrationCenter string      `yaml:"registration-center"`
	Alarm              AlarmConfig `yaml:"alarm"`
}

type ConsumerConfig struct {
	Name              string      `yaml:"name"`
	Topic             string      `yaml:"topic"`
	SubscriptionGroup string      `yaml:"subscription-group"`
	MQ                string      `yaml:"mq"`
	Alarm             AlarmConfig `yaml:"alarm"`
}

type AlarmConfig struct {
	Check     bool     `yaml:"check"`
	MinSize   int      `yaml:"min-size"`
	MustHosts []string `yaml:"must-hosts"`
}

type RocketMQClientConfig struct {
	Name string `yaml:"name"`
	Host string `yaml:"host"`
}

type RegistrationCenterConfig struct {
	Name  string      `yaml:"name"`
	Type  string      `yaml:"type"`
	Nacos NacosConfig `yaml:"nacos"`
}

type NacosConfig struct {
	Namespace  string `yaml:"namespace"`
	ServerAddr string `yaml:"server-addr"`
	Group      string `yaml:"group"`
	Username   string `yaml:"username"`
	Password   string `yaml:"password"`
}
