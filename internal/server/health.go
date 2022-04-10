package server

import "encoding/json"

type HealthCheck struct {
	Status string `json:"status"`
}

var (
	healthCheckUpJson string
)

func HealthCheckUpJson() string {
	if healthCheckUpJson == "" {
		check, _ := json.Marshal(HealthCheck{"UP"})
		healthCheckUpJson = string(check)
	}

	return healthCheckUpJson
}
