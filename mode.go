package zkregister

import (
	"time"

	"github.com/google/uuid"
)

type ServiceInfo struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Port    int    `json:"port"`
}

// Instance zk /services/{name}/{id} 路径下json信息
type Instance struct {
	Name    string `json:"name"`
	ID      string `json:"id"`
	Address string `json:"address"`
	Port    int    `json:"port"`
	SSLPort int    `json:"sslPort,omitempty"`
	Payload struct {
		Class    string `json:"@class"`
		ID       string `json:"id"`
		Name     string `json:"name"`
		Metadata struct {
			InstanceStatus string `json:"instance_status"`
		} `json:"metadata"`
	} `json:"payload"`
	RegistrationTimeUTC int64  `json:"registrationTimeUTC"`
	ServiceType         string `json:"serviceType"`
	URISpec             struct {
		Parts []struct {
			Value    string `json:"value"`
			Variable bool   `json:"variable"`
		} `json:"parts"`
	} `json:"uriSpec"`
}

func GetInstance(service ServiceInfo) Instance {
	uuid := uuid.New().String()
	instance := Instance{
		Name:    service.Name,
		ID:      uuid,
		Address: service.Address,
		Port:    service.Port,
		Payload: struct {
			Class    string `json:"@class"`
			ID       string `json:"id"`
			Name     string `json:"name"`
			Metadata struct {
				InstanceStatus string `json:"instance_status"`
			} `json:"metadata"`
		}{
			Class: "org.springframework.cloud.zookeeper.discovery.ZookeeperInstance",
			ID:    service.Name,
			Name:  service.Name,
			Metadata: struct {
				InstanceStatus string `json:"instance_status"`
			}{
				InstanceStatus: "UP",
			},
		},
		RegistrationTimeUTC: time.Now().UnixMilli(),
		ServiceType:         "DYNAMIC",
		URISpec: struct {
			Parts []struct {
				Value    string `json:"value"`
				Variable bool   `json:"variable"`
			} `json:"parts"`
		}{
			Parts: []struct {
				Value    string `json:"value"`
				Variable bool   `json:"variable"`
			}{
				{
					Value:    "scheme",
					Variable: true,
				},
				{
					Value:    "://",
					Variable: false,
				},
				{
					Value:    "address",
					Variable: true,
				},
				{
					Value:    ":",
					Variable: false,
				},
				{
					Value:    "port",
					Variable: true,
				},
			},
		},
	}
	return instance
}
