package discovery

import (
	"errors"
	"github.com/hashicorp/consul/api"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"libra-go/config"
	"strconv"
	"strings"
)

type DiscoveryClient interface {
	GetInstances(serviceId string) ([]ServiceInstance, error)
}

type ConsulDiscoveryClient struct {
	client api.Client
}

func NewConsulServiceRegistry(host string, port int, token string) (*ConsulDiscoveryClient, error) {
	if len(host) < 3 {
		return nil, errors.New("check host")
	}

	if port <= 0 || port > 65535 {
		return nil, errors.New("check port, port should between 1 and 65535")
	}

	config := api.DefaultConfig()
	config.Address = host + ":" + strconv.Itoa(port)
	config.Token = token
	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}

	return &ConsulDiscoveryClient{client: *client}, nil
}

func (c ConsulDiscoveryClient) GetInstances(serviceId string) ([]ServiceInstance, error) {
	catalogService, _, _ := c.client.Catalog().Service(serviceId, "", nil)
	if len(catalogService) > 0 {
		result := make([]ServiceInstance, len(catalogService))
		for index, sever := range catalogService {
			s := DefaultServiceInstance{
				ServiceId: sever.ServiceName,
				Host:      sever.ServiceAddress,
				Port:      sever.ServicePort,
				Metadata:  sever.ServiceMeta,
			}
			result[index] = s
		}
		return result, nil
	}
	return nil, nil
}

type NacosDiscoveryClient struct {
	nacosConfig config.NacosConfig
	iClient     naming_client.INamingClient
}

func (nacos NacosDiscoveryClient) GetInstances(serviceId string) ([]ServiceInstance, error) {
	instances, err := nacos.iClient.SelectInstances(vo.SelectInstancesParam{
		ServiceName: serviceId,
		GroupName:   nacos.nacosConfig.Group,
		HealthyOnly: true,
	})
	if err != nil {
		return nil, err
	}

	if len(instances) > 0 {
		result := make([]ServiceInstance, len(instances))
		for index, sever := range instances {
			s := DefaultServiceInstance{
				ServiceId: serviceId,
				Host:      sever.Ip,
				Port:      int(sever.Port),
				Metadata:  sever.Metadata,
			}
			result[index] = s
		}
		return result, nil
	}
	return nil, nil
}

func NewNacosDiscoveryClient(nacosConfig config.NacosConfig) (*NacosDiscoveryClient, error) {
	clientConfig := *constant.NewClientConfig(
		constant.WithNamespaceId(nacosConfig.Namespace),
		constant.WithTimeoutMs(5000),
		constant.WithNotLoadCacheAtStart(true),
		constant.WithUsername(nacosConfig.Username),
		constant.WithPassword(nacosConfig.Password),
		constant.WithUpdateCacheWhenEmpty(true),
		constant.WithLogLevel("warn"),
	)

	serverAddrSplit := strings.Split(nacosConfig.ServerAddr, ":")
	serverHost := serverAddrSplit[0]
	serverPort, err := strconv.ParseUint(serverAddrSplit[1], 10, 64)
	if err != nil {
		return nil, err
	}
	serverConfigs := []constant.ServerConfig{
		{
			IpAddr:      serverHost,
			ContextPath: "/nacos",
			Port:        serverPort,
			Scheme:      "http",
		},
	}

	namingClient, err := clients.CreateNamingClient(map[string]interface{}{
		"serverConfigs": serverConfigs,
		"clientConfig":  clientConfig,
	})
	if err != nil {
		return nil, err
	}
	return &NacosDiscoveryClient{
		nacosConfig: nacosConfig, iClient: namingClient,
	}, nil
}
