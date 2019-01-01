
CREATE EXTENSION IF NOT EXISTS citext;

DROP TABLE IF EXISTS post;
DROP TABLE IF EXISTS vote;
DROP TABLE IF EXISTS thread;
DROP TABLE IF EXISTS forum_user;
DROP TABLE IF EXISTS forum;
DROP TABLE IF EXISTS "user";

CREATE TABLE "user" (
  nickname citext PRIMARY KEY COLLATE "POSIX",
  fullname text,
  about text,
  email citext unique not null
);

CREATE INDEX idx_nick_nick ON "user" (nickname);
CREATE INDEX idx_nick_email ON "user" (email);
CREATE INDEX idx_nick_cover ON "user" (nickname, fullname, about, email);

CREATE TABLE forum (
  user_nick   citext references "user" not null,
  slug        citext PRIMARY KEY,
  title       text not null,
  thread_count integer default 0 not null,
  post_count integer default 0 not null
);

CREATE INDEX idx_forum_slug ON forum (slug);

CREATE TABLE forum_user (
  nickname citext references "user",
  forum_slug citext references "forum",
  CONSTRAINT unique_forum_user UNIQUE (nickname, forum_slug)
);

CREATE INDEX idx_forum_user ON forum_user (nickname, forum_slug);

CREATE TABLE thread (
  id BIGSERIAL PRIMARY KEY,
  slug citext unique ,
  forum_slug citext references forum not null,
  user_nick citext references "user" not null,
  created timestamp with time zone default now(),
  title text not null,
  votes integer default 0 not null,
  message text not null
);

CREATE INDEX idx_thread_id ON thread(id);
CREATE INDEX idx_thread_slug ON thread(slug);
CREATE INDEX idx_thread_f_slug ON thread(forum_slug, created);

CREATE TABLE vote (
  nickname citext references "user",
  voice boolean not null,
  thread_id integer references thread,
  CONSTRAINT unique_vote UNIQUE (nickname, thread_id)
);

CREATE INDEX idx_vote ON vote(thread_id, voice);

CREATE TABLE post (
  id BIGSERIAL PRIMARY KEY,
  path integer[],
  author citext references "user",
  created timestamp with time zone,
  edited boolean,
  message text,
  parent_id integer references post (id),
  forum_slug citext,
  thread_id integer references thread NOT NULL
);

CREATE INDEX idx_post_id ON post(id);
CREATE INDEX idx_post_thread_id ON post(thread_id);
CREATE INDEX idx_post_cr_id ON post(created, id, thread_id);
CREATE INDEX idx_post_thread_id_cr_i ON post(thread_id, id);
CREATE INDEX idx_post_thread_id_p_i ON post(thread_id, (path[1]), id);

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
     NEW.path := (ARRAY [NEW.id]);
     return NEW;
   end if;

   NEW.path := (SELECT array_append(p.path, NEW.id::integer)
                from post p where p.id = NEW.parent_id);
  RETURN NEW;
END;
$create_path$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS create_path ON post;

CREATE TRIGGER create_path BEFORE INSERT ON post
  FOR EACH ROW EXECUTE PROCEDURE create_path();