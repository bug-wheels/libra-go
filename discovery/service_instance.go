package discovery

type ServiceInstance interface {
	GetServiceId() string
	GetHost() string
	GetPort() int
	GetMetadata() map[string]string
}

type DefaultServiceInstance struct {
	ServiceId string
	Host      string
	Port      int
	Metadata  map[string]string
}

func NewDefaultServiceInstance(serviceId string, host string, port int,
	metadata map[string]string) (*DefaultServiceInstance, error) {
	return &DefaultServiceInstance{ServiceId: serviceId, Host: host, Port: port, Metadata: metadata}, nil
}

func (serviceInstance DefaultServiceInstance) GetServiceId() string {
	return serviceInstance.ServiceId
}

func (serviceInstance DefaultServiceInstance) GetHost() string {
	return serviceInstance.Host
}

func (serviceInstance DefaultServiceInstance) GetPort() int {
	return serviceInstance.Port
}

func (serviceInstance DefaultServiceInstance) GetMetadata() map[string]string {
	return serviceInstance.Metadata
}
