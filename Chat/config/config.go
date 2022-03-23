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

	loadConfigurations("DB_CONFIG", &DbConfig)
	loadConfigurations("APP_SERVER_PROPERTIES", &AppConfig)
	loadConfigurations("RABBIT_CONFIG", &RabbitConfig)
	loadConfigurations("JWT_CONFIG", &JwtConfig)
}

// Load the congifuration requested
func loadConfigurations(envVariable string, i interface{}) {
	dbConfig := os.Getenv(envVariable)

	if err := json.Unmarshal([]byte(dbConfig), i); err != nil {
		log.Fatalf("%s: ", err)
	}

}
