CREATE TYPE friendstatus AS enum (
  'PENDING',
  'FRIEND'
);

CREATE TABLE user_friend (
    id SERIAL PRIMARY KEY,
    uid1 uuid NOT NULL,
    uid2 uuid NOT NULL,
    friend_status friendstatus NOT NULL
);