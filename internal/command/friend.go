package command

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/git-fal7/friend-for-gate/internal/database"
	"github.com/google/uuid"
	"go.minekube.com/brigodier"
	"go.minekube.com/common/minecraft/component"
	"go.minekube.com/gate/pkg/command"
	"go.minekube.com/gate/pkg/edition/java/proxy"
)

func friendCommand(p *proxy.Proxy) brigodier.LiteralNodeBuilder {
	return brigodier.Literal("friend").
		Executes(command.Command(func(c *command.Context) error {
			player, ok := c.Source.(proxy.Player)
			if !ok {
				return nil
			}
			player.SendMessage(&component.Text{
				Content: "/friend add [player]",
			})
			player.SendMessage(&component.Text{
				Content: "/friend remove [player]",
			})
			player.SendMessage(&component.Text{
				Content: "/friend accept [player]",
			})
			player.SendMessage(&component.Text{
				Content: "/friend list",
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
				friends, err := database.DB.ListFriends(context.Background(), uuid.UUID(player.ID()))
				if err != nil {
					player.SendMessage(&component.Text{
						Content: "An error occured, please try again",
					})
					log.Println(err)
					return nil
				}
				for _, friend := range friends {
					friendUUID := friend.Uid1
					if friendUUID == uuid.UUID(player.ID()) {
						friendUUID = friend.Uid2
					}
					// This is stupid, should be cached in memory or we do a single query instead of multiple.
					lookupResult, err := database.DB.GetUsernameFromLookupTable(context.Background(), friendUUID)
					if err != nil {
						log.Println(err)
						continue
					}
					player.SendMessage(&component.Text{
						Content: fmt.Sprintf("Friend: %s", lookupResult.UserName),
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
			if target == nil {
				// Lookup player
				lookupUserResult, err := database.DB.GetUserUUIDFromLookupTable(context.Background(), strings.ToLower(targetStr))
				if err != nil {
					player.SendMessage(&component.Text{
						Content: "Invalid player",
					})
					return nil
				}
				targetUUID = lookupUserResult.UserUuid
			} else {
				targetUUID = uuid.UUID(target.ID())
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
							Content: fmt.Sprintf("You have already a request pending with %s", target.Username()),
						})
						return nil
					} else if friendStatus == database.FriendstatusFRIEND {
						player.SendMessage(&component.Text{
							Content: fmt.Sprintf("You are already friends with %s", target.Username()),
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
					Content: fmt.Sprintf("Sent friend request to %s", target.Username()),
				})
				target.SendMessage(&component.Text{
					Content: fmt.Sprintf("%s sent you a friend request, accept by /friend accept %s", player.Username(), player.Username()),
				})
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
					Content: fmt.Sprintf("Removed %s", target.Username()),
				})
			} else if strings.ToLower(arg1) == "accept" {
				// check if have relations
				getFriendStatusParams := database.GetFriendStatusParams{
					Uid1: uuid.UUID(player.ID()),
					Uid2: targetUUID,
				}
				friendStatus, err := database.DB.GetFriendStatus(context.Background(), getFriendStatusParams)
				if err == nil {
					if friendStatus == database.FriendstatusFRIEND {
						player.SendMessage(&component.Text{
							Content: fmt.Sprintf("You are already friends with %s", target.Username()),
						})
						return nil
					}
				}

				acceptFriendRequetsParam := database.AcceptFriendRequestParams{
					Uid1: uuid.UUID(player.ID()),
					Uid2: targetUUID,
				}
				err = database.DB.AcceptFriendRequest(context.Background(), acceptFriendRequetsParam)
				if err != nil {
					player.SendMessage(&component.Text{
						Content: "An error occured, please try again",
					})
					log.Println(err)
					return nil
				}
				player.SendMessage(&component.Text{
					Content: fmt.Sprintf("You are now friends with %s", target.Username()),
				})
				target.SendMessage(&component.Text{
					Content: fmt.Sprintf("You are now friends with %s", player.Username()),
				})
			}
			return nil
		})),
	))

}
