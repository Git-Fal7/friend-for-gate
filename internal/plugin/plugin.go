package plugin

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/git-fal7/friend-for-gate/internal/command"
	"github.com/git-fal7/friend-for-gate/internal/database"
	"github.com/git-fal7/friend-for-gate/internal/event"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.minekube.com/gate/pkg/edition/java/proxy"
)

func InitPlugin(ctx context.Context, proxy *proxy.Proxy) error {
	db, err := sql.Open("pgx", fmt.Sprintf("user=%s password=%s host=%s port=%s database=%s sslmode=disable", "admin", "adminpassword", "localhost", "5432", "friends"))
	if err != nil {
		return err
	}
	database.DB = database.New(db)
	// Init tables
	err = database.DB.CreateFriendstatusType(context.Background())
	if err != nil {
		return err
	}
	err = database.DB.CreateUsersTable(context.Background())
	if err != nil {
		return err
	}
	err = database.DB.CreateLookupUserTable(context.Background())
	if err != nil {
		return err
	}
	command.Init(proxy)
	event.Init(proxy)

	return nil
}
