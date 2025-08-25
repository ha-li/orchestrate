package config

import (
	"encoding/json"
	"os"
)

type DockerConfig struct {
	Name  string `json:"docker_name"`
	Image string `json:"docker_image"`
}

type ServerConfig struct {
	Port int `json:"port"`
}

type DatabaseConfig struct {
	User     string `json:"user"`
	Password string `json:"psswd"`
}

type AppConfig struct {
	Docker   DockerConfig   `json:"docker"`
	Server   ServerConfig   `json:"server"`
	Database DatabaseConfig `json:"database"`
}

func New() *AppConfig {
	return &AppConfig{}
}

func (c *AppConfig) LoadConfig(path string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	err = json.Unmarshal(content, c)
	if err != nil {
		return err
	}

	return nil
}

func (c *AppConfig) GetDockerConfig() *DockerConfig {
	return &c.Docker
}

func (c *AppConfig) GetServerConfig() ServerConfig {
	return c.Server
}

func (c *AppConfig) GetDatabaseConfig() *DatabaseConfig {
	return &c.Database
}

// database configurations
func (db *DatabaseConfig) GetUser() string {
	return db.User
}

func (db *DatabaseConfig) GetPassword() string {
	return db.Password
}

func (dc *DockerConfig) GetName() string {
	return dc.Name
}

func (dc *DockerConfig) GetImage() string {
	return dc.Image
}
