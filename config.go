package main

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"os"
)

func setDefaultGlobalConfig() {
	// set default
	viper.SetDefault("PORT", 3000)

	if viper.GetString("APP_KEY") == "" {
		panic("APP_KEY is required.")
	}

	viper.SetDefault("APP_ENV", "production")

	viper.SetDefault("DB_NAME", "rent-n-go")
	viper.SetDefault("DB_USER", "root")
	viper.SetDefault("DB_PASS", "")
	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_PORT", 3306)
}

/**
* Load configuration
 */
func init() {
	// load from file
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed, please re run:", e.Name)
		os.Exit(1)
	})
	viper.WatchConfig()
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error regrading config file: %w", err))
	}

	setDefaultGlobalConfig()
}
