package config

import "github.com/spf13/viper"

var ViperConfig = viper.New()

func InitConfig() {
	ViperConfig.SetConfigName("friendforgate")
	ViperConfig.SetConfigType("yaml")
	ViperConfig.ReadInConfig()

	ViperConfig.SetDefault("database.hostname", "localhost")
	ViperConfig.SetDefault("database.port", 5432)
	ViperConfig.SetDefault("database.username", "admin")
	ViperConfig.SetDefault("database.password", "adminpassword")
	ViperConfig.SetDefault("database.database", "friends")

	ViperConfig.SetDefault("messages.friendHelpMessage",
		`/friend add [player]
		/friend remove [player]
		/friend accept [player]
		/friend list`)

}
