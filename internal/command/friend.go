package command

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/git-fal7/friend-for-gate/internal/config"
	"github.com/git-fal7/friend-for-gate/internal/database"
	"github.com/google/uuid"
	"go.minekube.com/brigodier"
	"go.minekube.com/common/minecraft/component"
	"go.minekube.com/gate/pkg/command"
	"go.minekube.com/gate/pkg/edition/java/proxy"
)

func replaceAll(str string, replaceMap map[string]string) string {
	for key, value := range replaceMap {
		str = strings.ReplaceAll(str, key, value)
	}
	return str
}

func friendCommand(p *proxy.Proxy) brigodier.LiteralNodeBuilder {
	return brigodier.Literal("friend").
		Executes(command.Command(func(c *command.Context) error {
			player, ok := c.Source.(proxy.Player)
			if !ok {
				return nil
			}
			player.SendMessage(&component.Text{
				Content: config.ViperConfig.GetString("messages.friendHelpMessage"),
			})
			return nil
		})).Then(brigodier.Argument("arg-1", brigodier.String).
		Executes(command.Command(func(c *command.Context) error {
			player, ok := c.Source.(proxy.Player)
			if !ok {
				return nil
			}
			arg1 := c.String("arg-1")
			if strings.ToLower(arg1) == "list" {
				friends, err := database.DB.ListFriendsLookup(context.Background(), uuid.UUID(player.ID()))
				if err != nil {
					player.SendMessage(&component.Text{
						Content: "An error occured, please try again",
					})
					log.Println(err)
					return nil
				}
				for _, friend := range friends {
					player.SendMessage(&component.Text{
						Content: replaceAll(config.ViperConfig.GetString("messages.friendListFriend"), map[string]string{
							"%friend%": friend.UserName,
						}),
					})
				}
			}
			return nil
		})).Then(brigodier.Argument("player", brigodier.String).
		Executes(command.Command(func(c *command.Context) error {
			player, ok := c.Source.(proxy.Player)
			if !ok {
				return nil
			}
			arg1 := c.String("arg-1")
			targetStr := c.String("player")
			target := p.PlayerByName(targetStr)
			var targetUUID uuid.UUID
			var targetUsername string
			if target == nil {
				// Lookup player
				lookupUserResult, err := database.DB.GetUserUUIDFromLookupTable(context.Background(), targetStr)
				if err != nil {
					player.SendMessage(&component.Text{
						Content: config.ViperConfig.GetString("messages.errorPlayerNotFound"),
					})
					return nil
				}
				targetUUID = lookupUserResult.UserUuid
				targetUsername = lookupUserResult.UserName
			} else {
				targetUUID = uuid.UUID(target.ID())
				targetUsername = target.Username()
			}
			if strings.ToLower(arg1) == "add" {
				// check if have relations
				getFriendStatusParams := database.GetFriendStatusParams{
					Uid1: uuid.UUID(player.ID()),
					Uid2: targetUUID,
				}
				friendStatus, err := database.DB.GetFriendStatus(context.Background(), getFriendStatusParams)
				if err == nil {
					if friendStatus == database.FriendstatusPENDING {
						player.SendMessage(&component.Text{
							Content: replaceAll(config.ViperConfig.GetString("messages.errorAlreadyRequest"), map[string]string{
								"%sender%":   player.Username(),
								"%receiver%": targetUsername,
							}),
						})
						return nil
					} else if friendStatus == database.FriendstatusFRIEND {
						player.SendMessage(&component.Text{
							Content: replaceAll(config.ViperConfig.GetString("messages.errorAlreadyFriends"), map[string]string{
								"%sender%":   player.Username(),
								"%receiver%": targetUsername,
							}),
						})
						return nil
					}
				}

				createFriendRequestParams := database.CreateFriendRequestParams{
					Uid1: uuid.UUID(player.ID()),
					Uid2: targetUUID,
				}
				err = database.DB.CreateFriendRequest(context.Background(), createFriendRequestParams)
				if err != nil {
					player.SendMessage(&component.Text{
						Content: "An error occured, please try again",
					})
					log.Println(err)
					return nil
				}
				player.SendMessage(&component.Text{
					Content: replaceAll(config.ViperConfig.GetString("messages.friendSenderSentRequest"), map[string]string{
						"%sender%":   player.Username(),
						"%receiver%": targetUsername,
					}),
				})
				if target != nil {
					target.SendMessage(&component.Text{
						Content: replaceAll(config.ViperConfig.GetString("messages.friendReceiverSentRequest"), map[string]string{
							"%sender%":   player.Username(),
							"%receiver%": targetUsername,
						}),
					})
				}
			} else if strings.ToLower(arg1) == "remove" {
				// remove
				removeFriendRequestParam := database.RemoveFriendRequestParams{
					Uid1: uuid.UUID(player.ID()),
					Uid2: targetUUID,
				}
				err := database.DB.RemoveFriendRequest(context.Background(), removeFriendRequestParam)
				if err != nil {
					player.SendMessage(&component.Text{
						Content: "An error occured, please try again",
					})
					log.Println(err)
					return nil
				}
				// just send to the player
				player.SendMessage(&component.Text{
					Content: fmt.Sprintf("Removed %s", targetUsername),
				})
			} else if strings.ToLower(arg1) == "accept" {
				// check if have relations
				getFriendStatusParams := database.GetFriendStatusParams{
					Uid1: uuid.UUID(player.ID()),
					Uid2: targetUUID,
				}
				friendStatus, err := database.DB.GetFriendStatus(context.Background(), getFriendStatusParams)
				if err != nil {
					player.SendMessage(&component.Text{
						Content: replaceAll(config.ViperConfig.GetString("messages.errorNoRequest"), map[string]string{
							"%sender%":   player.Username(),
							"%receiver%": targetUsername,
						}),
					})
					log.Println(err)
					return nil
				}
				if friendStatus == database.FriendstatusFRIEND {
					player.SendMessage(&component.Text{
						Content: replaceAll(config.ViperConfig.GetString("messages.errorAlreadyFriends"), map[string]string{
							"%sender%":   player.Username(),
							"%receiver%": targetUsername,
						}),
					})
					return nil
				}

				acceptFriendRequetsParam := database.AcceptFriendRequestParams{
					Uid1: targetUUID,
					Uid2: uuid.UUID(player.ID()),
				}
				err = database.DB.AcceptFriendRequest(context.Background(), acceptFriendRequetsParam)
				if err != nil {
					player.SendMessage(&component.Text{
						Content: "An error occured, please try again",
					})
					log.Println(err)
					return nil
				}
				friendStatus, _ = database.DB.GetFriendStatus(context.Background(), getFriendStatusParams)
				if friendStatus == database.FriendstatusFRIEND {
					player.SendMessage(&component.Text{
						Content: replaceAll(config.ViperConfig.GetString("messages.friendNowFriends"), map[string]string{
							"%player%": targetUsername,
						}),
					})
					if target != nil {
						target.SendMessage(&component.Text{
							Content: replaceAll(config.ViperConfig.GetString("messages.friendNowFriends"), map[string]string{
								"%player%": player.Username(),
							}),
						})
					}
				} else {
					player.SendMessage(&component.Text{
						Content: replaceAll(config.ViperConfig.GetString("messages.errorAlreadyRequest"), map[string]string{
							"%sender%":   player.Username(),
							"%receiver%": targetUsername,
						}),
					})
				}
			}
			return nil
		})),
	))

}
