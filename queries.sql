-- name: CreateUsersTable :exec
CREATE TABLE user_friend (
    id SERIAL PRIMARY KEY,
    uid1 uuid NOT NULL,
    uid2 uuid NOT NULL,
    friend_status enum('REQ_UID1', 'REQ_UID2', 'FRIEND') NOT NULL
);

-- name: CreateFriendRequest :exec
INSERT INTO user_friend (
    uid1, uid2, friend_status
) VALUES (
    $1, $2, 'REQ_UID1'
);