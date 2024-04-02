-- name: CreateUsersTable :exec
CREATE TABLE IF NOT EXISTS user_friend (
    id SERIAL PRIMARY KEY,
    uid1 uuid NOT NULL,
    uid2 uuid NOT NULL,
    friend_status friendstatus NOT NULL
);

-- name: CreateFriendstatusType :exec
DO $$ BEGIN
    IF to_regtype('friendstatus') IS NULL THEN
        CREATE TYPE friendstatus AS enum (
            'PENDING',
            'FRIEND'
        );
    END IF;
END $$;

-- name: CreateLookupUserTable :exec
CREATE TABLE lookup_users (
    user_uuid uuid PRIMARY KEY NOT NULL,
    user_name varchar(16) NOT NULL
);


-- name: CreateFriendRequest :exec
INSERT INTO user_friend (
    uid1, uid2, friend_status
) VALUES (
    $1, $2, 'PENDING'
);

-- name: GetFriendStatus :one
SELECT friend_status FROM user_friend
WHERE (uid1 = $1 AND uid2 = $2) OR (uid1 = $2 AND uid2 = $1) 
LIMIT 1;

-- name: AcceptFriendRequest :exec
UPDATE user_friend
SET friend_status = 'FRIEND'
WHERE (uid1 = $1 AND uid2 = $2) OR (uid1 = $2 AND uid2 = $1);

-- name: RemoveFriendRequest :exec
DELETE FROM user_friend
WHERE (uid1 = $1 AND uid2 = $2) OR (uid1 = $2 AND uid2 = $1);

-- name: ListFriends :many
SELECT * FROM user_friend
WHERE (uid1 = $1 OR uid2 = $1) AND friend_status = 'FRIEND';

-- name: LogIntoLookupTable :exec
INSERT INTO lookup_users (
    user_uuid, user_name
) VALUES (
    $1, $2
)
ON CONFLICT(user_uuid)
DO UPDATE SET
user_name = $2;

-- name: GetUsernameFromLookupTable :one
SELECT user_name FROM lookup_users
WHERE user_uuid = $1 
LIMIT 1;