package command

import (
	"context"
	"fmt"

	"github.com/git-fal7/friend-for-gate/internal/config"
	"github.com/git-fal7/friend-for-gate/internal/database"
	"github.com/google/uuid"
	"go.minekube.com/brigodier"
	"go.minekube.com/common/minecraft/component"
	"go.minekube.com/gate/pkg/command"
	"go.minekube.com/gate/pkg/edition/java/proxy"
)

func msgCommand(p *proxy.Proxy) brigodier.LiteralNodeBuilder {
	return brigodier.Literal("msg").
		Executes(command.Command(func(c *command.Context) error {
			player, ok := c.Source.(proxy.Player)
			if !ok {
				return nil
			}
			player.SendMessage(&component.Text{
				Content: "/msg [player] [message]",
			})
			return nil
		})).Then(brigodier.Argument("player", brigodier.String).
		Executes(command.Command(func(c *command.Context) error {
			player, ok := c.Source.(proxy.Player)
			if !ok {
				return nil
			}
			player.SendMessage(&component.Text{
				Content: "/msg [player] [message]",
			})
			return nil
		})).Then(brigodier.Argument("message", brigodier.StringPhrase).
		Executes(command.Command(func(c *command.Context) error {
			player, ok := c.Source.(proxy.Player)
			if !ok {
				return nil
			}
			message := c.String("message")
			target := p.PlayerByName(c.String("player"))
			if target == nil {
				player.SendMessage(&component.Text{
					Content: config.ViperConfig.GetString("messages.errorPlayerNotFound"),
				})
				return nil
			}
			getFriendStatusParams := database.GetFriendStatusParams{
				Uid1: uuid.UUID(player.ID()),
				Uid2: uuid.UUID(target.ID()),
			}
			friendStatus, err := database.DB.GetFriendStatus(context.Background(), getFriendStatusParams)
			if err != nil || friendStatus != database.FriendstatusFRIEND {
				player.SendMessage(&component.Text{
					Content: "That player isn't in your friend list",
				})
				return nil
			}
			player.SendMessage(&component.Text{
				Content: fmt.Sprintf("%s >> %s: %s", player.Username(), target.Username(), message),
			})
			target.SendMessage(&component.Text{
				Content: fmt.Sprintf("%s >> %s: %s", player.Username(), target.Username(), message),
			})
			return nil
		})),
	))
}
