// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: queries.sql

package database

import (
	"context"

	"github.com/google/uuid"
)

const acceptFriendRequest = `-- name: AcceptFriendRequest :exec
UPDATE user_friend
SET friend_status = 'FRIEND'
WHERE (uid1 = $1 AND uid2 = $2) OR (uid1 = $2 AND uid2 = $1)
`

type AcceptFriendRequestParams struct {
	Uid1 uuid.UUID
	Uid2 uuid.UUID
}

func (q *Queries) AcceptFriendRequest(ctx context.Context, arg AcceptFriendRequestParams) error {
	_, err := q.db.ExecContext(ctx, acceptFriendRequest, arg.Uid1, arg.Uid2)
	return err
}

const createFriendRequest = `-- name: CreateFriendRequest :exec
INSERT INTO user_friend (
    uid1, uid2, friend_status
) VALUES (
    $1, $2, 'PENDING'
)
`

type CreateFriendRequestParams struct {
	Uid1 uuid.UUID
	Uid2 uuid.UUID
}

func (q *Queries) CreateFriendRequest(ctx context.Context, arg CreateFriendRequestParams) error {
	_, err := q.db.ExecContext(ctx, createFriendRequest, arg.Uid1, arg.Uid2)
	return err
}

const createFriendstatusType = `-- name: CreateFriendstatusType :exec
DO $$ BEGIN
    IF to_regtype('friendstatus') IS NULL THEN
        CREATE TYPE friendstatus AS enum (
            'PENDING',
            'FRIEND'
        );
    END IF;
END $$
`

func (q *Queries) CreateFriendstatusType(ctx context.Context) error {
	_, err := q.db.ExecContext(ctx, createFriendstatusType)
	return err
}

const createLookupUserTable = `-- name: CreateLookupUserTable :exec
CREATE TABLE lookup_users (
    user_uuid uuid PRIMARY KEY NOT NULL,
    user_name varchar(16) NOT NULL
)
`

func (q *Queries) CreateLookupUserTable(ctx context.Context) error {
	_, err := q.db.ExecContext(ctx, createLookupUserTable)
	return err
}

const createUsersTable = `-- name: CreateUsersTable :exec
CREATE TABLE IF NOT EXISTS user_friend (
    id SERIAL PRIMARY KEY,
    uid1 uuid NOT NULL,
    uid2 uuid NOT NULL,
    friend_status friendstatus NOT NULL
)
`

func (q *Queries) CreateUsersTable(ctx context.Context) error {
	_, err := q.db.ExecContext(ctx, createUsersTable)
	return err
}

const getFriendStatus = `-- name: GetFriendStatus :one
SELECT friend_status FROM user_friend
WHERE (uid1 = $1 AND uid2 = $2) OR (uid1 = $2 AND uid2 = $1) 
LIMIT 1
`

type GetFriendStatusParams struct {
	Uid1 uuid.UUID
	Uid2 uuid.UUID
}

func (q *Queries) GetFriendStatus(ctx context.Context, arg GetFriendStatusParams) (Friendstatus, error) {
	row := q.db.QueryRowContext(ctx, getFriendStatus, arg.Uid1, arg.Uid2)
	var friend_status Friendstatus
	err := row.Scan(&friend_status)
	return friend_status, err
}

const getUserUUIDFromLookupTable = `-- name: GetUserUUIDFromLookupTable :one
SELECT user_uuid, user_name FROM lookup_users
WHERE user_name = $1 
LIMIT 1
`

func (q *Queries) GetUserUUIDFromLookupTable(ctx context.Context, userName string) (LookupUser, error) {
	row := q.db.QueryRowContext(ctx, getUserUUIDFromLookupTable, userName)
	var i LookupUser
	err := row.Scan(&i.UserUuid, &i.UserName)
	return i, err
}

const getUsernameFromLookupTable = `-- name: GetUsernameFromLookupTable :one
SELECT user_uuid, user_name FROM lookup_users
WHERE user_uuid = $1 
LIMIT 1
`

func (q *Queries) GetUsernameFromLookupTable(ctx context.Context, userUuid uuid.UUID) (LookupUser, error) {
	row := q.db.QueryRowContext(ctx, getUsernameFromLookupTable, userUuid)
	var i LookupUser
	err := row.Scan(&i.UserUuid, &i.UserName)
	return i, err
}

const listFriends = `-- name: ListFriends :many
SELECT id, uid1, uid2, friend_status FROM user_friend
WHERE (uid1 = $1 OR uid2 = $1) AND friend_status = 'FRIEND'
`

func (q *Queries) ListFriends(ctx context.Context, uid1 uuid.UUID) ([]UserFriend, error) {
	rows, err := q.db.QueryContext(ctx, listFriends, uid1)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []UserFriend
	for rows.Next() {
		var i UserFriend
		if err := rows.Scan(
			&i.ID,
			&i.Uid1,
			&i.Uid2,
			&i.FriendStatus,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const logIntoLookupTable = `-- name: LogIntoLookupTable :exec
INSERT INTO lookup_users (
    user_uuid, user_name
) VALUES (
    $1, $2
)
ON CONFLICT(user_uuid)
DO UPDATE SET
user_name = $2
`

type LogIntoLookupTableParams struct {
	UserUuid uuid.UUID
	UserName string
}

func (q *Queries) LogIntoLookupTable(ctx context.Context, arg LogIntoLookupTableParams) error {
	_, err := q.db.ExecContext(ctx, logIntoLookupTable, arg.UserUuid, arg.UserName)
	return err
}

const removeFriendRequest = `-- name: RemoveFriendRequest :exec
DELETE FROM user_friend
WHERE (uid1 = $1 AND uid2 = $2) OR (uid1 = $2 AND uid2 = $1)
`

type RemoveFriendRequestParams struct {
	Uid1 uuid.UUID
	Uid2 uuid.UUID
}

func (q *Queries) RemoveFriendRequest(ctx context.Context, arg RemoveFriendRequestParams) error {
	_, err := q.db.ExecContext(ctx, removeFriendRequest, arg.Uid1, arg.Uid2)
	return err
}
