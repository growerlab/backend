--
-- PostgreSQL database dump
--

-- Dumped from database version 12.0
-- Dumped by pg_dump version 12.0

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: activate_code; Type: TABLE; Schema: public; Owner: growerlab
--

CREATE TABLE public.activate_code (
    id integer NOT NULL,
    user_id integer NOT NULL,
    code character varying(16) NOT NULL,
    created_at bigint NOT NULL,
    used_at bigint,
    expired_at bigint NOT NULL
);


ALTER TABLE public.activate_code OWNER TO growerlab;

--
-- Name: TABLE activate_code; Type: COMMENT; Schema: public; Owner: growerlab
--

COMMENT ON TABLE public.activate_code IS '用户激活码';


--
-- Name: activate_code_id_seq; Type: SEQUENCE; Schema: public; Owner: growerlab
--

CREATE SEQUENCE public.activate_code_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.activate_code_id_seq OWNER TO growerlab;

--
-- Name: activate_code_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: growerlab
--

ALTER SEQUENCE public.activate_code_id_seq OWNED BY public.activate_code.id;


--
-- Name: namespace; Type: TABLE; Schema: public; Owner: growerlab
--

CREATE TABLE public.namespace (
    id bigint NOT NULL,
    path character varying(255) NOT NULL,
    owner_id integer NOT NULL
);


ALTER TABLE public.namespace OWNER TO growerlab;

--
-- Name: TABLE namespace; Type: COMMENT; Schema: public; Owner: growerlab
--

COMMENT ON TABLE public.namespace IS ' 命名空间';


--
-- Name: COLUMN namespace.path; Type: COMMENT; Schema: public; Owner: growerlab
--

COMMENT ON COLUMN public.namespace.path IS '路径';


--
-- Name: COLUMN namespace.owner_id; Type: COMMENT; Schema: public; Owner: growerlab
--

COMMENT ON COLUMN public.namespace.owner_id IS '所有者';


--
-- Name: namespace_id_seq; Type: SEQUENCE; Schema: public; Owner: growerlab
--

CREATE SEQUENCE public.namespace_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.namespace_id_seq OWNER TO growerlab;

--
-- Name: namespace_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: growerlab
--

ALTER SEQUENCE public.namespace_id_seq OWNED BY public.namespace.id;


--
-- Name: repository; Type: TABLE; Schema: public; Owner: growerlab
--

CREATE TABLE public.repository (
    id bigint NOT NULL,
    uuid character varying(16) NOT NULL,
    path character varying(255) NOT NULL,
    name character varying(255) NOT NULL,
    namespace_id bigint NOT NULL,
    owner_id bigint NOT NULL,
    description text,
    created_at bigint NOT NULL
);


ALTER TABLE public.repository OWNER TO growerlab;

--
-- Name: TABLE repository; Type: COMMENT; Schema: public; Owner: growerlab
--

COMMENT ON TABLE public.repository IS '仓库表';


--
-- Name: COLUMN repository.uuid; Type: COMMENT; Schema: public; Owner: growerlab
--

COMMENT ON COLUMN public.repository.uuid IS '仓库uuid（fork仓库相同）';


--
-- Name: COLUMN repository.path; Type: COMMENT; Schema: public; Owner: growerlab
--

COMMENT ON COLUMN public.repository.path IS '仓库路径';


--
-- Name: COLUMN repository.name; Type: COMMENT; Schema: public; Owner: growerlab
--

COMMENT ON COLUMN public.repository.name IS '仓库名';


--
-- Name: COLUMN repository.namespace_id; Type: COMMENT; Schema: public; Owner: growerlab
--

COMMENT ON COLUMN public.repository.namespace_id IS '命名空间id';


--
-- Name: COLUMN repository.owner_id; Type: COMMENT; Schema: public; Owner: growerlab
--

COMMENT ON COLUMN public.repository.owner_id IS '仓库创建者,fork后不变';


--
-- Name: COLUMN repository.description; Type: COMMENT; Schema: public; Owner: growerlab
--

COMMENT ON COLUMN public.repository.description IS '仓库描述';


--
-- Name: repository_id_seq; Type: SEQUENCE; Schema: public; Owner: growerlab
--

CREATE SEQUENCE public.repository_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.repository_id_seq OWNER TO growerlab;

--
-- Name: repository_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: growerlab
--

ALTER SEQUENCE public.repository_id_seq OWNED BY public.repository.id;


--
-- Name: session; Type: TABLE; Schema: public; Owner: growerlab
--

CREATE TABLE public.session (
    id bigint NOT NULL,
    user_id bigint NOT NULL,
    token character varying(36) NOT NULL,
    created_at bigint NOT NULL,
    expired_at bigint NOT NULL
);


ALTER TABLE public.session OWNER TO growerlab;

--
-- Name: session_id_seq; Type: SEQUENCE; Schema: public; Owner: growerlab
--

CREATE SEQUENCE public.session_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.session_id_seq OWNER TO growerlab;

--
-- Name: session_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: growerlab
--

ALTER SEQUENCE public.session_id_seq OWNED BY public.session.id;


--
-- Name: user; Type: TABLE; Schema: public; Owner: growerlab
--

CREATE TABLE public."user" (
    id bigint NOT NULL,
    email character varying(255) NOT NULL,
    encrypted_password character varying(255) NOT NULL,
    username character varying(40) NOT NULL,
    name character varying(255) NOT NULL,
    public_email character varying(255) NOT NULL,
    last_login_ip character varying(46),
    created_at bigint NOT NULL,
    deleted_at bigint,
    verified_at bigint,
    last_login_at bigint
);


ALTER TABLE public."user" OWNER TO growerlab;

--
-- Name: TABLE "user"; Type: COMMENT; Schema: public; Owner: growerlab
--

COMMENT ON TABLE public."user" IS '用户表';


--
-- Name: COLUMN "user".email; Type: COMMENT; Schema: public; Owner: growerlab
--

COMMENT ON COLUMN public."user".email IS '用户邮箱';


--
-- Name: COLUMN "user".encrypted_password; Type: COMMENT; Schema: public; Owner: growerlab
--

COMMENT ON COLUMN public."user".encrypted_password IS '用户密码';


--
-- Name: COLUMN "user".username; Type: COMMENT; Schema: public; Owner: growerlab
--

COMMENT ON COLUMN public."user".username IS '唯一性用户名（将用在url中）';


--
-- Name: COLUMN "user".name; Type: COMMENT; Schema: public; Owner: growerlab
--

COMMENT ON COLUMN public."user".name IS '用户昵称';


--
-- Name: COLUMN "user".public_email; Type: COMMENT; Schema: public; Owner: growerlab
--

COMMENT ON COLUMN public."user".public_email IS '公开的邮箱地址';


--
-- Name: COLUMN "user".last_login_ip; Type: COMMENT; Schema: public; Owner: growerlab
--

COMMENT ON COLUMN public."user".last_login_ip IS '最后的登录ip（兼容ipv6长度）';


--
-- Name: user_id_seq; Type: SEQUENCE; Schema: public; Owner: growerlab
--

CREATE SEQUENCE public.user_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.user_id_seq OWNER TO growerlab;

--
-- Name: user_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: growerlab
--

ALTER SEQUENCE public.user_id_seq OWNED BY public."user".id;


--
-- Name: activate_code id; Type: DEFAULT; Schema: public; Owner: growerlab
--

ALTER TABLE ONLY public.activate_code ALTER COLUMN id SET DEFAULT nextval('public.activate_code_id_seq'::regclass);


--
-- Name: namespace id; Type: DEFAULT; Schema: public; Owner: growerlab
--

ALTER TABLE ONLY public.namespace ALTER COLUMN id SET DEFAULT nextval('public.namespace_id_seq'::regclass);


--
-- Name: repository id; Type: DEFAULT; Schema: public; Owner: growerlab
--

ALTER TABLE ONLY public.repository ALTER COLUMN id SET DEFAULT nextval('public.repository_id_seq'::regclass);


--
-- Name: session id; Type: DEFAULT; Schema: public; Owner: growerlab
--

ALTER TABLE ONLY public.session ALTER COLUMN id SET DEFAULT nextval('public.session_id_seq'::regclass);


--
-- Name: user id; Type: DEFAULT; Schema: public; Owner: growerlab
--

ALTER TABLE ONLY public."user" ALTER COLUMN id SET DEFAULT nextval('public.user_id_seq'::regclass);


--
-- Name: activate_code activate_code_pk; Type: CONSTRAINT; Schema: public; Owner: growerlab
--

ALTER TABLE ONLY public.activate_code
    ADD CONSTRAINT activate_code_pk PRIMARY KEY (id);


--
-- Name: namespace namespace_pk; Type: CONSTRAINT; Schema: public; Owner: growerlab
--

ALTER TABLE ONLY public.namespace
    ADD CONSTRAINT namespace_pk PRIMARY KEY (id);


--
-- Name: repository repository_pk; Type: CONSTRAINT; Schema: public; Owner: growerlab
--

ALTER TABLE ONLY public.repository
    ADD CONSTRAINT repository_pk PRIMARY KEY (id);


--
-- Name: session session_pk; Type: CONSTRAINT; Schema: public; Owner: growerlab
--

ALTER TABLE ONLY public.session
    ADD CONSTRAINT session_pk PRIMARY KEY (id);


--
-- Name: user user_pkey; Type: CONSTRAINT; Schema: public; Owner: growerlab
--

ALTER TABLE ONLY public."user"
    ADD CONSTRAINT user_pkey PRIMARY KEY (id);


--
-- Name: activate_code_code_uindex; Type: INDEX; Schema: public; Owner: growerlab
--

CREATE UNIQUE INDEX activate_code_code_uindex ON public.activate_code USING btree (code);


--
-- Name: namespace_id_uindex; Type: INDEX; Schema: public; Owner: growerlab
--

CREATE UNIQUE INDEX namespace_id_uindex ON public.namespace USING btree (id);


--
-- Name: namespace_owner_id_index; Type: INDEX; Schema: public; Owner: growerlab
--

CREATE INDEX namespace_owner_id_index ON public.namespace USING btree (owner_id);


--
-- Name: namespace_path_index; Type: INDEX; Schema: public; Owner: growerlab
--

CREATE INDEX namespace_path_index ON public.namespace USING btree (path);


--
-- Name: repository_path_index; Type: INDEX; Schema: public; Owner: growerlab
--

CREATE INDEX repository_path_index ON public.repository USING btree (path);


--
-- Name: repository_uuid_index; Type: INDEX; Schema: public; Owner: growerlab
--

CREATE INDEX repository_uuid_index ON public.repository USING btree (uuid);


--
-- Name: session_user_id_token_uniq; Type: INDEX; Schema: public; Owner: growerlab
--

CREATE UNIQUE INDEX session_user_id_token_uniq ON public.session USING btree (user_id, token);


--
-- Name: user_email_uindex; Type: INDEX; Schema: public; Owner: growerlab
--

CREATE UNIQUE INDEX user_email_uindex ON public."user" USING btree (email);


--
-- Name: user_username_uindex; Type: INDEX; Schema: public; Owner: growerlab
--

CREATE UNIQUE INDEX user_username_uindex ON public."user" USING btree (username);


--
-- PostgreSQL database dump complete
--

