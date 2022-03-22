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

	loadDbConfiguration()

	loadAppConfiguration()
	loadRabbitConfiguration()
	loadJWTConfiguration()
}

// Load de db congifurations
func loadDbConfiguration() {
	dbConfig := os.Getenv("DB_CONFIG")

	if err := json.Unmarshal([]byte(dbConfig), &DbConfig); err != nil {
		log.Fatalf("%s: ", err)
	}

}

func loadAppConfiguration() {
	appConfig := os.Getenv("APP_SERVER_PROPERTIES")

	if err := json.Unmarshal([]byte(appConfig), &AppConfig); err != nil {
		log.Fatalf("%s: ", err)
	}
}

func loadRabbitConfiguration() {
	rabbitConfig := os.Getenv("RABBIT_CONFIG")

	if err := json.Unmarshal([]byte(rabbitConfig), &RabbitConfig); err != nil {
		log.Fatalf("%s: ", err)
	}
}

func loadJWTConfiguration() {
	jwtConfig := os.Getenv("JWT_CONFIG")

	if err := json.Unmarshal([]byte(jwtConfig), &JwtConfig); err != nil {
		log.Fatalf("%s: ", err)
	}
}
