// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package database

import (
	"github.com/google/uuid"
)

type UserFriend struct {
	ID           int32
	Uid1         uuid.UUID
	Uid2         uuid.UUID
	FriendStatus interface{}
}
