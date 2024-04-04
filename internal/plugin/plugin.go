package plugin

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/git-fal7/friend-for-gate/internal/command"
	"github.com/git-fal7/friend-for-gate/internal/config"
	"github.com/git-fal7/friend-for-gate/internal/database"
	"github.com/git-fal7/friend-for-gate/internal/event"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/mattn/go-sqlite3"
	"go.minekube.com/gate/pkg/edition/java/proxy"
)

func InitPlugin(ctx context.Context, proxy *proxy.Proxy) error {
	config.InitConfig()
	// postgres/mysql -> sqlite
	var db *sql.DB
	var err error
	if strings.ToLower(config.ViperConfig.GetString("database.type")) == "postgresql" {
		db, err = sql.Open("pgx", fmt.Sprintf("user=%s password=%s host=%s port=%s database=%s",
			config.ViperConfig.GetString("database.username"), config.ViperConfig.GetString("database.password"),
			config.ViperConfig.GetString("database.hostname"), config.ViperConfig.GetString("database.port"),
			config.ViperConfig.GetString("database.database")))
	} else if strings.ToLower(config.ViperConfig.GetString("database.type")) == "mysql" {
		db, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
			config.ViperConfig.GetString("database.username"), config.ViperConfig.GetString("database.password"),
			config.ViperConfig.GetString("database.hostname"), config.ViperConfig.GetString("database.port"),
			config.ViperConfig.GetString("database.database")))
	} else {
		db, err = sql.Open("sqlite3", "./friendforgate.db")
	}
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
