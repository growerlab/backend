-- auto-generated definition
CREATE TABLE "user" (
  id bigserial NOT NULL CONSTRAINT user_pkey PRIMARY KEY,
  email varchar(255) NOT NULL,
  encrypted_password varchar(255) NOT NULL,
  username varchar(255) NOT NULL,
  name varchar(255) NOT NULL,
  public_email varchar(255),
  created_at timestamp,
  deleted_at timestamp
);

COMMENT ON TABLE "user" IS '用户表';

COMMENT ON COLUMN "user".email IS '用户邮箱';

COMMENT ON COLUMN "user".encrypted_password IS '用户密码';

COMMENT ON COLUMN "user".username IS '唯一性用户名（将用在url中）';

COMMENT ON COLUMN "user".name IS '用户昵称';

COMMENT ON COLUMN "user".public_email IS '公开的邮箱地址';

COMMENT ON COLUMN "user".created_at IS '创建的时间';

COMMENT ON COLUMN "user".deleted_at IS '删除的时间';

ALTER TABLE "user" OWNER TO growerlab;

CREATE UNIQUE INDEX user_email_uindex ON "user" (email);

CREATE INDEX user_public_email_index ON "user" (public_email);

CREATE UNIQUE INDEX user_username_uindex ON "user" (username);

CREATE TABLE IF NOT EXISTS repository (
  id bigserial NOT NULL CONSTRAINT repository_pk PRIMARY KEY,
  uuid varchar(16) NOT NULL,
  path varchar(255) NOT NULL,
  name varchar(255) NOT NULL,
  namespace_id integer NOT NULL,
  creator_id integer NOT NULL,
  description text
);

COMMENT ON TABLE repository IS '仓库表';

COMMENT ON COLUMN repository.uuid IS '仓库uuid（fork仓库相同）';

COMMENT ON COLUMN repository.path IS '仓库路径';

COMMENT ON COLUMN repository.name IS '仓库名';

COMMENT ON COLUMN repository.namespace_id IS '命名空间id';

COMMENT ON COLUMN repository.creator_id IS '仓库创建者';

COMMENT ON COLUMN repository.description IS '仓库描述';

ALTER TABLE repository OWNER TO growerlab;

CREATE INDEX IF NOT EXISTS repository_path_index ON repository (path);

CREATE INDEX IF NOT EXISTS repository_uuid_index ON repository (uuid);

CREATE TABLE IF NOT EXISTS namespace (
  id bigserial NOT NULL CONSTRAINT namespace_pk PRIMARY KEY,
  path varchar(255) NOT NULL,
  owner_id integer NOT NULL
);

COMMENT ON TABLE namespace IS ' 命名空间';

COMMENT ON COLUMN namespace.path IS '路径';

COMMENT ON COLUMN namespace.owner_id IS '所有者';

ALTER TABLE namespace OWNER TO growerlab;

CREATE UNIQUE INDEX IF NOT EXISTS namespace_id_uindex ON namespace (id);

CREATE INDEX IF NOT EXISTS namespace_owner_id_index ON namespace (owner_id);

CREATE INDEX IF NOT EXISTS namespace_path_index ON namespace (path);

