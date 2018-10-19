
CREATE EXTENSION IF NOT EXISTS citext;


DROP TABLE IF EXISTS thread;
DROP TABLE IF EXISTS forum;
DROP TABLE IF EXISTS "user";

CREATE TABLE "user" (
  nickname citext PRIMARY KEY,
  fullname text,
  about text,
  email citext unique not null
);

CREATE TABLE forum (
  user_nick   citext references "user",
  slug        citext PRIMARY KEY,
  title       text not null
);

CREATE TABLE thread (
  id BIGSERIAL PRIMARY KEY,
  forum_slug citext references forum,
  user_nick citext references "user",
  created timestamp,
  slug text not null,
  title text not null,
  message text not null
);