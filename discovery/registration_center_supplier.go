package discovery

import (
	"libra-go/config"
)

var RegistrationCenters = make(map[string]DiscoveryClient)

func Setup(registrationCenters []config.RegistrationCenterConfig) {
	for name := range RegistrationCenters {
		delete(RegistrationCenters, name)
	}
	if len(registrationCenters) == 0 {
		return
	}

	for _, rcConfig := range registrationCenters {
		if rcConfig.Type == "nacos" {
			client, err := NewNacosDiscoveryClient(rcConfig.Nacos)
			if err != nil {
				panic(err)
			}
			RegistrationCenters[rcConfig.Name] = client
		}
	}
}
