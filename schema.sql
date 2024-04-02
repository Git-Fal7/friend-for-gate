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

CREATE TABLE lookup_users (
  user_uuid uuid PRIMARY KEY NOT NULL,
  user_name varchar(16) NOT NULL
);