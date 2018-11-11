
CREATE EXTENSION IF NOT EXISTS citext;

DROP TABLE IF EXISTS post;
DROP TABLE IF EXISTS vote;
DROP TABLE IF EXISTS thread;
DROP TABLE IF EXISTS forum;
DROP TABLE IF EXISTS "user";

CREATE TABLE "user" (
  nickname citext PRIMARY KEY COLLATE "POSIX",
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
  slug citext unique not null ,
  forum_slug citext references forum,
  user_nick citext references "user",
  created timestamp with time zone default now(),
  title text not null,
  votes integer default 0 not null,
  message text not null
);

CREATE TABLE vote (
  id BIGSERIAL PRIMARY KEY,
  nickname citext references "user",
  voice boolean not null,
  thread_id integer references thread,
  CONSTRAINT unique_vote UNIQUE (nickname, thread_id)
);

CREATE TABLE post (
  id BIGSERIAL PRIMARY KEY,
  path integer[],
  author citext references "user",
  created timestamp with time zone,
  edited boolean,
  message text,
  parent_id integer references post (id),
  thread_id integer references thread NOT NULL
);

CREATE OR REPLACE FUNCTION change_edited_post() RETURNS trigger as $change_edited_post$
BEGIN
  IF NEW.message <> OLD.message THEN
    NEW.edited = true;
  END IF;
  
  return NEW;
END;
$change_edited_post$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS change_edited_post ON post;

CREATE TRIGGER change_edited_post BEFORE UPDATE ON post
  FOR EACH ROW EXECUTE PROCEDURE change_edited_post();

CREATE OR REPLACE FUNCTION create_path() RETURNS trigger as $create_path$
BEGIN
   IF NEW.parent_id IS NULL THEN
     NEW.path := (ARRAY [(SELECT COUNT(*) + 1 from post p where p.parent_id is null and p.thread_id = NEW.thread_id )]);
     return NEW;
   end if;

   NEW.path := (SELECT array_append(p.path, (SELECT COUNT(*)::INTEGER + 1 from post p1 where p1.parent_id = NEW.parent_id and p1.thread_id = NEW.thread_id))
                from post p where p.id = NEW.parent_id);
  RETURN NEW;
END;
$create_path$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS create_path ON post;

CREATE TRIGGER create_path BEFORE INSERT ON post
  FOR EACH ROW EXECUTE PROCEDURE create_path();