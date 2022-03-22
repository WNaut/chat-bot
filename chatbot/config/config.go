package config

import (
	"encoding/json"
	"log"
	"os"
)

type AppServerConfiguration struct {
	Port string `json:"port"`
}

type JWTConfiguration struct {
	SecretKey string `json:"secretKey"`
}

type DbConfiguration struct {
	Host string `json:"host"`
	Db   string `json:"db"`
	User string `json:"user"`
	Pwd  string `json:"pwd"`
}

type RabbitConfiguration struct {
	Host        string `json:"host"`
	Port        string `json:"port"`
	User        string `json:"user"`
	Pwd         string `json:"pwd"`
	ClientQueue string `json:"clientQueue"`
	StooqQueue  string `json:"stooqQueue"`
}

var (
	DbConfig     DbConfiguration
	AppConfig    AppServerConfiguration
	RabbitConfig RabbitConfiguration
	JwtConfig    JWTConfiguration
)

// Load all the required configurations
func Load() {
	loadDbConfigurations("RABBIT_CONFIG", &RabbitConfig)
}

// Load de db congifurations
func loadDbConfigurations(envVariable string, i interface{}) {
	dbConfig := os.Getenv(envVariable)

	if err := json.Unmarshal([]byte(dbConfig), i); err != nil {
		log.Fatalf("%s: ", err)
	}

}
