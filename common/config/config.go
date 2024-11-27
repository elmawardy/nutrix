package config

import (
	"context"
	"fmt"
	"time"

	"github.com/elmawardy/nutrix/common/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ConfigFactory creates a Config object based on the provided type and path
func ConfigFactory(t string, path string, logger logger.ILogger) Config {
	switch t {
	case "viper":
		viper_config := NewViperConfig(logger)
		viper_config.ReadFile(path)
		config, err := viper_config.GetConfig()
		if err != nil {
			logger.Error("can't reat config")
		}
		return config
	}

	return Config{}
}

// IConfig interface defines methods for config management
type IConfig interface {
	ReadFile(path string)
	GetConfig() (config Config)
}

// Settings represents the configuration settings structure
type Settings struct {
	Id        string `bson:"id,omitempty" json:"id"`
	Inventory struct {
		DefaultInventoryQuantityWarn float64 `json:"default_inventory_quantity_warn" bson:"default_inventory_quantity_warn"`
	} `bson:"inventory" json:"inventory"`
	Orders struct {
		Queues []struct {
			Prefix string `json:"prefix" bson:"prefix"`
			Next   uint32 `json:"next" bson:"next"`
		} `json:"queues" bson:"queues"`
	} `bson:"orders" json:"orders"`
}

// LoadFromDB loads settings from the database using the provided Config
func (s *Settings) LoadFromDB(config Config) error {
	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%v", config.Databases[0].Host, config.Databases[0].Port))

	db_connection_deadline := 5 * time.Second
	if config.Env == "dev" {
		db_connection_deadline = 1000 * time.Second
	}

	// Create a context with a timeout (optional)
	ctx, cancel := context.WithTimeout(context.Background(), db_connection_deadline)

	defer cancel()

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return err
	}

	// Ping the database to check connectivity
	err = client.Ping(ctx, nil)
	if err != nil {
		return err
	}

	// Get the "settings" collection from the database
	collection := client.Database("waha").Collection("settings")
	err = collection.FindOne(ctx, bson.D{}).Decode(s)
	if err != nil {
		return err
	}

	return nil
}

// ZitadelConfig holds the configuration for Zitadel
type ZitadelConfig struct {
	Domain  string `mapstructure:"domain"`
	Port    uint8  `mapstructure:"port"`
	KeyPath string `mapstructure:"key_path"`
}

// Config represents the overall configuration structure
type Config struct {
	Databases    []Database
	Zitadel      ZitadelConfig `mapstructure:"zitadel"`
	Env          string        `mapstructure:"env"`
	JwtSecretKey string        `mapstructure:"jwt_secret_key"`
	TimeZone     string        `mapstructure:"timezone"`
}

// Database holds the configuration for database connections
type Database struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Type     string `mapstructure:"type"`
	Name     string `mapstructure:"name"`
	Database string `mapstructure:"database"`
}