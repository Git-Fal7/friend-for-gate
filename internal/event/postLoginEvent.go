package event

import (
	"context"
	"log"

	"github.com/git-fal7/friend-for-gate/internal/database"
	"github.com/google/uuid"
	"go.minekube.com/gate/pkg/edition/java/proxy"
)

func loginEvent(p *proxy.Proxy) func(*proxy.PostLoginEvent) {
	return func(e *proxy.PostLoginEvent) {
		// log into lookup table
		logIntoLookupTableParam := database.LogIntoLookupTableParams{
			UserUuid: uuid.UUID(e.Player().ID()),
			UserName: e.Player().Username(),
		}
		err := database.DB.LogIntoLookupTable(context.Background(), logIntoLookupTableParam)
		if err != nil {
			log.Println(err)
		}
	}
}
