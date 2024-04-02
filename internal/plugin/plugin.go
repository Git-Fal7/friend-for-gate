package plugin

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/git-fal7/friend-for-gate/internal/database"
	"go.minekube.com/gate/pkg/edition/java/proxy"
)

var DB *database.Queries

func InitPlugin(ctx context.Context, proxy *proxy.Proxy) error {
	db, err := sql.Open("pgx", fmt.Sprintf("user=%s password=%s host=%s port=%s database=%s sslmode=disable", "admin", "adminpassword", "localhost", "5432", "friends"))
	if err != nil {
		return err
	}
	DB = database.New(db)
	// Init tables
	err = DB.CreateUsersTable(context.Background())
	if err != nil {
		return err
	}

	return nil
}
