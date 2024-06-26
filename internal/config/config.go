package config

import (
	"log"

	"github.com/spf13/viper"
)

var ViperConfig = viper.New()

func InitConfig() {
	ViperConfig.AddConfigPath(".")
	ViperConfig.SetConfigName("friendforgate")
	ViperConfig.SetConfigType("yaml")

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
	ViperConfig.SetDefault("messages.friendSenderSentRequest", "Sent friend request to %receiver%")
	ViperConfig.SetDefault("messages.friendReceiverSentRequest", "%sender% sent you a friend request, accept by /friend accept %sender%")
	ViperConfig.SetDefault("messages.friendNowFriends", "You are now friends with %player%")
	ViperConfig.SetDefault("messages.friendListFriend", "Friendlist: %friend%")
	ViperConfig.SetDefault("messages.friendSucessfullyRemoved", "Removed %friend%")
	ViperConfig.SetDefault("messages.msgHelpMessage", "/msg [player] [message]")
	ViperConfig.SetDefault("messages.msgOutgoingMessage", "To %target%: %message%")
	ViperConfig.SetDefault("messages.msgIncomingMessage", "From %receiver%: %message%")
	ViperConfig.SetDefault("messages.errorPlayerNotFound", "Couldnt find that player")
	ViperConfig.SetDefault("messages.errorAlreadyRequest", "You have already a request pending with %receiver%")
	ViperConfig.SetDefault("messages.errorAlreadyFriends", "You have already friends with %receiver%")
	ViperConfig.SetDefault("messages.errorNoRequest", "You dont have a request from this player")
	ViperConfig.SetDefault("messages.errorNotInFriendList", "That player is not in your friend list")
	ViperConfig.SetDefault("messages.errorOccured", "An error occured, please try again")
	ViperConfig.SetDefault("messages.errorSelfAdd", "You cant add yourself")

	err := ViperConfig.ReadInConfig()
	if err != nil {
		// Create config file
		log.Println("Couldn't find friendforgate.yml, creating a new config file")
		ViperConfig.WriteConfigAs("./friendforgate.yml")
	}

}
