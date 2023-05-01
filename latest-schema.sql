--
-- PostgreSQL database dump
--

-- Dumped from database version 12.9
-- Dumped by pg_dump version 13.1

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

--
-- Name: fuzzystrmatch; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS fuzzystrmatch WITH SCHEMA public;


--
-- Name: EXTENSION fuzzystrmatch; Type: COMMENT; Schema: -; Owner: -
--

COMMENT ON EXTENSION fuzzystrmatch IS 'determine similarities and distance between strings';


--
-- Name: pg_stat_statements; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS pg_stat_statements WITH SCHEMA public;


--
-- Name: EXTENSION pg_stat_statements; Type: COMMENT; Schema: -; Owner: -
--

COMMENT ON EXTENSION pg_stat_statements IS 'track execution statistics of all SQL statements executed';


--
-- Name: pg_trgm; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS pg_trgm WITH SCHEMA public;


--
-- Name: EXTENSION pg_trgm; Type: COMMENT; Schema: -; Owner: -
--

COMMENT ON EXTENSION pg_trgm IS 'text similarity measurement and index searching based on trigrams';


SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: apply_token; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.apply_token (
    token character(27) NOT NULL PRIMARY KEY,
    job_id character(27) NOT NULL,
    created_at timestamp without time zone NOT NULL,
    confirmed_at timestamp without time zone,
    email character varying(255) NOT NULL,
    cv bytea NOT NULL
);


--
-- Name: cloudflare_browser_stats; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.cloudflare_browser_stats (
    date date NOT NULL,
    page_views bigint NOT NULL,
    ua_browser_family character varying(255) NOT NULL
);


--
-- Name: cloudflare_country_stats; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.cloudflare_country_stats (
    date date NOT NULL,
    country_code character varying(255) NOT NULL,
    requests bigint NOT NULL,
    threats bigint NOT NULL
);


--
-- Name: cloudflare_stats; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.cloudflare_stats (
    date date NOT NULL,
    bytes bigint NOT NULL,
    cached_bytes bigint NOT NULL,
    page_views bigint NOT NULL,
    requests bigint NOT NULL,
    threats bigint NOT NULL,
    uniques bigint NOT NULL
);


--
-- Name: cloudflare_status_code_stats; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.cloudflare_status_code_stats (
    date date NOT NULL,
    status_code integer NOT NULL,
    requests bigint NOT NULL
);


--
-- Name: company; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.company (
    id character(27) NOT NULL PRIMARY KEY,
    name character varying(255) NOT NULL,
    url character varying(255) NOT NULL,
    locations character varying(255) NOT NULL,
    last_job_created_at timestamp without time zone NOT NULL,
    icon_image_id character(27) NOT NULL,
    total_job_count integer NOT NULL,
    active_job_count integer NOT NULL,
    description text,
    featured_post_a_job boolean DEFAULT false,
    slug character varying(255) DEFAULT NULL::character varying NOT NULL,
    twitter character varying(255) DEFAULT NULL::character varying,
    github character varying(255) DEFAULT NULL::character varying,
    linkedin character varying(255) DEFAULT NULL::character varying,
    company_page_eligibility_expired_at TIMESTAMP DEFAULT '1970-01-01 00:00:00'
);


--
-- Name: company_event; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.company_event (
    id character(27) NOT NULL PRIMARY KEY,
    event_type character varying(128) NOT NULL,
    company_id character(27) NOT NULL,
    created_at timestamp without time zone NOT NULL
);


--
-- Name: developer_profile; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.developer_profile (
    id character(27) NOT NULL PRIMARY KEY,
    email character varying(255) NOT NULL,
    location character varying(255) NOT NULL,
    available boolean NOT NULL,
    linkedin_url character varying(255) NOT NULL,
    hourly_rate integer NOT NULL default 0,
    image_id character(27) NOT NULL,
    slug character varying(255) NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone,
    skills character varying(255) DEFAULT 'Go'::character varying NOT NULL,
    name character varying(255) NOT NULL,
    bio text NOT NULL,
    github_url character varying(255) DEFAULT NULL::character varying,
    twitter_url character varying(255) DEFAULT NULL::character varying
);


--
-- Name: developer_profile_event; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.developer_profile_event (
    id character(27) NOT NULL PRIMARY KEY,
    event_type character varying(128) NOT NULL,
    developer_profile_id character(27) NOT NULL,
    created_at timestamp without time zone NOT NULL
);


--
-- Name: developer_profile_message; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.developer_profile_message (
    id character(27) NOT NULL PRIMARY KEY,
    email character varying(255) NOT NULL,
    content text NOT NULL,
    profile_id character(27) NOT NULL,
    sender_id character(27) NOT NULL,
    created_at timestamp without time zone NOT NULL,
    sent_at timestamp without time zone
);


--
-- Name: edit_token; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.edit_token (
    token character(27) NOT NULL PRIMARY KEY,
    job_id character(27) NOT NULL,
    created_at timestamp without time zone NOT NULL
);


--
-- Name: email_notification; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.email_notification (
    id character(27) NOT NULL PRIMARY KEY,
    email character varying(255) NOT NULL,
    event_type character varying(100) NOT NULL,
    created_at timestamp without time zone NOT NULL,
    job_id character(27) NOT NULL
);


--
-- Name: email_subscribers; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.email_subscribers (
    email character varying(255) NOT NULL PRIMARY KEY,
    token character(27) NOT NULL UNIQUE,
    confirmed_at timestamp without time zone,
    created_at timestamp without time zone NOT NULL
);


--
-- Name: fx_rate; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.fx_rate (
    base character(3) NOT NULL,
    target character(3) NOT NULL,
    value double precision NOT NULL,
    updated_at timestamp without time zone NOT NULL
);


--
-- Name: image; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.image (
    id character(27) NOT NULL PRIMARY KEY,
    bytes bytea NOT NULL,
    media_type character varying(100) NOT NULL
);


--
-- Name: job; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.job (
    id character(27) NOT NULL PRIMARY KEY,
    job_title character varying(128) NOT NULL,
    job_category character varying(128) NOT NULL,
    company character varying(128) NOT NULL,
    location character varying(200) NOT NULL,
    salary_range character varying(100) NOT NULL,
    job_type character varying(100) NOT NULL,
    application_link text NOT NULL,
    subscriber_email character varying(255) NOT NULL,
    description text NOT NULL,
    created_at timestamp without time zone NOT NULL,
    approved_at timestamp without time zone,
    url_id integer NOT NULL,
    slug character varying(256),
    company_icon_image_id character varying(255) DEFAULT NULL::character varying,
    external_id character varying(28) DEFAULT ''::character varying NOT NULL,
    expired boolean DEFAULT false,
    newsletter_eligibility_expired_at TIMESTAMP DEFAULT '1970-01-01 00:00:00',
    social_media_eligibility_expired_at TIMESTAMP DEFAULT '1970-01-01 00:00:00',
    blog_eligibility_expired_at TIMESTAMP DEFAULT '1970-01-01 00:00:00',
    front_page_eligibility_expired_at TIMESTAMP DEFAULT '1970-01-01 00:00:00',
    company_page_eligibility_expired_at TIMESTAMP DEFAULT '1970-01-01 00:00:00',
    plan_expired_at TIMESTAMP DEFAULT '1970-01-01 00:00:00',
    last_week_clickouts integer DEFAULT 0 NOT NULL
);


--
-- Name: job_event; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.job_event (
    id character(27) NOT NULL PRIMARY KEY,
    event_type character varying(128) NOT NULL,
    job_id character(27) NOT NULL,
    created_at timestamp without time zone NOT NULL
);

--
-- Name: meta; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.meta (
    key character varying(255) NOT NULL PRIMARY KEY,
    value character varying(255) NOT NULL
);


--
-- Name: purchase_event; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.purchase_event (
    stripe_session_id character varying(255) NOT NULL PRIMARY KEY,
    email character varying(255) DEFAULT ''::character varying NOT NULL,
    plan_id character varying(255) DEFAULT ''::character varying NOT NULL,
    description character varying(255) NOT NULL,
    job_id character(27) NOT NULL,
    created_at timestamp without time zone NOT NULL,
    completed_at timestamp without time zone
);


--
-- Name: search_event; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.search_event (
    session_id character varying(255) NOT NULL PRIMARY KEY,
    location character varying(255) DEFAULT NULL::character varying,
    tag character varying(255) DEFAULT NULL::character varying,
    results integer NOT NULL,
    created_at timestamp without time zone NOT NULL,
    type character varying(10) DEFAULT 'job'::character varying
);


--
-- Name: seo_landing_page; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.seo_landing_page (
    uri character varying(255) NOT NULL PRIMARY KEY,
    location character varying(255) NOT NULL,
    skill character varying(255) NOT NULL
);


--
-- Name: seo_location; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.seo_location (
    name character varying(255) NOT NULL PRIMARY KEY,
    currency character varying(4) DEFAULT '$'::character varying NOT NULL,
    country character varying(255) DEFAULT NULL::character varying,
    iso2 character(2) DEFAULT NULL::bpchar,
    region character varying(255) DEFAULT NULL::character varying,
    population integer,
    lat double precision,
    long double precision,
    emoji character varying(5)
);


--
-- Name: seo_salary; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.seo_salary (
    id character varying(255) NOT NULL PRIMARY KEY,
    location character varying(255) NOT NULL,
    currency character varying(5) NOT NULL,
    uri character varying(100) NOT NULL
);


--
-- Name: seo_skill; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.seo_skill (
    name character varying(255) NOT NULL PRIMARY KEY
);


--
-- Name: sitemap; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.sitemap (
    loc character varying(255) PRIMARY KEY,
    changefreq character varying(20),
    lastmod timestamp without time zone
);


--
-- Name: user_sign_on_token; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.user_sign_on_token (
    token character(27) NOT NULL PRIMARY KEY,
    email character varying(255) NOT NULL,
    user_type VARCHAR(20) DEFAULT 'developer' NOT NULL
);


--
-- Name: users; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.users (
    id character(27) NOT NULL PRIMARY KEY,
    email character varying(255) NOT NULL,
    user_type varchar(20) NOT NULL default 'developer',
    role_level VARCHAR(20) NOT NULL DEFAULT 'mid-level',
    search_status VARCHAR(20) NOT NULL DEFAULT 'casually-looking',
    role_types VARCHAR(60) NOT NULL DEFAULT 'full-time',
    detected_location_id VARCHAR(255) DEFAULT NULL,
    created_at timestamp without time zone
);

--
-- Name: cloudflare_browser_stats cloudflare_browser_stats_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cloudflare_browser_stats
    ADD CONSTRAINT cloudflare_browser_stats_pkey PRIMARY KEY (date, ua_browser_family);


--
-- Name: cloudflare_country_stats cloudflare_country_stats_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cloudflare_country_stats
    ADD CONSTRAINT cloudflare_country_stats_pkey PRIMARY KEY (date, country_code);


--
-- Name: cloudflare_stats cloudflare_stats_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cloudflare_stats
    ADD CONSTRAINT cloudflare_stats_pkey PRIMARY KEY (date);


--
-- Name: cloudflare_status_code_stats cloudflare_status_code_stats_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.cloudflare_status_code_stats
    ADD CONSTRAINT cloudflare_status_code_stats_pkey PRIMARY KEY (date, requests);

--
-- Name: fx_rate fx_rate_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.fx_rate
    ADD CONSTRAINT fx_rate_pkey PRIMARY KEY (base, target);

--
-- Name: company_slug_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX company_slug_idx ON public.company USING btree (slug);


--
-- Name: developer_profile_email_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX developer_profile_email_idx ON public.developer_profile USING btree (email);


--
-- Name: developer_profile_slug_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX developer_profile_slug_idx ON public.developer_profile USING btree (slug);


--
-- Name: purchase_event_job_id_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX purchase_event_job_id_idx ON public.purchase_event USING btree (job_id);


--
-- Name: slug_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX slug_idx ON public.job USING btree (slug);


--
-- Name: url_id_idx; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX url_id_idx ON public.job USING btree (url_id);

--
-- Name: apply_token apply_token_job_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.apply_token
    ADD CONSTRAINT apply_token_job_id_fkey FOREIGN KEY (job_id) REFERENCES public.job(id);


--
-- Name: company_event company_event_company_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.company_event
    ADD CONSTRAINT company_event_company_id_fkey FOREIGN KEY (company_id) REFERENCES public.company(id);


--
-- Name: developer_profile_event developer_profile_event_developer_profile_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.developer_profile_event
    ADD CONSTRAINT developer_profile_event_developer_profile_id_fkey FOREIGN KEY (developer_profile_id) REFERENCES public.developer_profile(id);


--
-- Name: developer_profile developer_profile_image_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.developer_profile
    ADD CONSTRAINT developer_profile_image_id_fk FOREIGN KEY (image_id) REFERENCES public.image(id);


--
-- Name: developer_profile_message developer_profile_message_profile_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.developer_profile_message
    ADD CONSTRAINT developer_profile_message_profile_id_fk FOREIGN KEY (profile_id) REFERENCES public.developer_profile(id);

--
-- Name: developer_profile_message developer_profile_message_sender_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.developer_profile_message
    ADD CONSTRAINT developer_profile_message_sender_id_fk FOREIGN KEY (sender_id) REFERENCES public.users(id);


--
-- Name: edit_token edit_token_job_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.edit_token
    ADD CONSTRAINT edit_token_job_id_fkey FOREIGN KEY (job_id) REFERENCES public.job(id);


--
-- Name: job_event job_event_job_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.job_event
    ADD CONSTRAINT job_event_job_id_fkey FOREIGN KEY (job_id) REFERENCES public.job(id);


--
-- Name: purchase_event purchase_event_job_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.purchase_event
    ADD CONSTRAINT purchase_event_job_id_fkey FOREIGN KEY (job_id) REFERENCES public.job(id);


CREATE TABLE IF NOT EXISTS public.blog_post (
	id CHAR(27) NOT NULL PRIMARY KEY,
	title VARCHAR(255) NOT NULL,
	description VARCHAR(255) NOT NULL,
	slug VARCHAR(255) NOT NULL,
	tags VARCHAR(255) NOT NULL,
	text TEXT NOT NULL,
	created_at TIMESTAMP NOT NULL,
	updated_at TIMESTAMP NOT NULL,
	created_by CHAR(27) NOT NULL,
	published_at TIMESTAMP DEFAULT NULL
);

 CREATE UNIQUE INDEX blog_post_slug_idx on public.blog_post (slug);
 
 CREATE TABLE "public"."recruiter_profile" (
    "id" char(27) NOT NULL PRIMARY KEY,
    "email" varchar(255) NOT NULL,
    "company_url" varchar(255) NOT NULL,
    "slug" varchar(255) NOT NULL,
    "created_at" timestamp NOT NULL,
    "updated_at" timestamp,
    "name" varchar(255)
);


--
-- PostgreSQL database dump complete
--