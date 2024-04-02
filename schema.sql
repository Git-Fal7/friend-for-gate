CREATE TABLE user_friend (
    id SERIAL PRIMARY KEY,
    uid1 uuid NOT NULL,
    uid2 uuid NOT NULL,
    friend_status enum('REQ_UID1', 'REQ_UID2', 'FRIEND') NOT NULL
);