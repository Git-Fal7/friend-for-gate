-- name: CreateUsersTable :exec
CREATE TABLE user_friend (
    id SERIAL PRIMARY KEY,
    uid1 uuid NOT NULL,
    uid2 uuid NOT NULL,
    friend_status friendstatus NOT NULL
);

-- name: CreateFriendstatusType :exec
CREATE TYPE friendstatus AS enum (
  'PENDING',
  'FRIEND'
);

-- name: CreateFriendRequest :exec
INSERT INTO user_friend (
    uid1, uid2, friend_status
) VALUES (
    $1, $2, 'PENDING'
);

-- name: GetFriendStatus :one
SELECT friend_status FROM user_friend
WHERE (uid1 = $1 AND uid2 = $2) OR (uid1 = $2 AND uid = $1) 
LIMIT 1;

-- name: AcceptFriendRequest :exec
UPDATE user_friend
SET friend_status = 'FRIEND'
WHERE (uid1 = $1 AND uid2 = $2) OR (uid1 = $2 AND uid = $1);

-- name: RemoveFriendRequest :exec
DELETE FROM user_friend
WHERE (uid1 = $1 AND uid2 = $2) OR (uid1 = $2 AND uid = $1);