\restrict gS2J33OkUIrV7I06wU83K5clmDBxKYL3o6tw9t3pEkEHrDf9hC4R7JvBUT8fuGh

-- Dumped from database version 18.3 (Debian 18.3-1.pgdg12+1)
-- Dumped by pg_dump version 18.3 (Debian 18.3-1.pgdg13+1)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET transaction_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: asset_type_t; Type: TYPE; Schema: public; Owner: -
--

CREATE TYPE public.asset_type_t AS ENUM (
    'image',
    'animated',
    'video',
    'audio'
);


--
-- Name: process_status_t; Type: TYPE; Schema: public; Owner: -
--

CREATE TYPE public.process_status_t AS ENUM (
    'pending',
    'processing',
    'processed'
);


SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: albums; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.albums (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    name text NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    modified_at timestamp with time zone DEFAULT now() NOT NULL,
    deleted_at timestamp with time zone,
    cover uuid
);


--
-- Name: assets; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.assets (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    album_id uuid NOT NULL,
    filename text NOT NULL,
    checksum text NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    modified_at timestamp with time zone DEFAULT now() NOT NULL,
    deleted_at timestamp with time zone,
    size bigint DEFAULT 0 NOT NULL,
    type public.asset_type_t DEFAULT 'image'::public.asset_type_t NOT NULL,
    original text DEFAULT ''::text NOT NULL,
    preview text DEFAULT ''::text NOT NULL,
    thumbnail text DEFAULT ''::text NOT NULL,
    view text DEFAULT ''::text NOT NULL,
    process_status public.process_status_t DEFAULT 'pending'::public.process_status_t NOT NULL,
    thumbnail_width integer DEFAULT 0 NOT NULL,
    thumbnail_height integer DEFAULT 0 NOT NULL,
    view_width integer DEFAULT 0 NOT NULL,
    view_height integer DEFAULT 0 NOT NULL,
    image_frames integer DEFAULT 0 NOT NULL,
    video_duration interval DEFAULT '00:00:00'::interval NOT NULL
);


--
-- Name: schema_migrations; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.schema_migrations (
    version character varying NOT NULL
);


--
-- Name: albums albums_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.albums
    ADD CONSTRAINT albums_pkey PRIMARY KEY (id);


--
-- Name: assets assets_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.assets
    ADD CONSTRAINT assets_pkey PRIMARY KEY (id);


--
-- Name: schema_migrations schema_migrations_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.schema_migrations
    ADD CONSTRAINT schema_migrations_pkey PRIMARY KEY (version);


--
-- Name: assets assets_album_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.assets
    ADD CONSTRAINT assets_album_id_fkey FOREIGN KEY (album_id) REFERENCES public.albums(id);


--
-- PostgreSQL database dump complete
--

\unrestrict gS2J33OkUIrV7I06wU83K5clmDBxKYL3o6tw9t3pEkEHrDf9hC4R7JvBUT8fuGh


--
-- Dbmate schema migrations
--

INSERT INTO public.schema_migrations (version) VALUES
    ('20260405151552'),
    ('20260405162721'),
    ('20260406133102'),
    ('20260406213309'),
    ('20260410152556'),
    ('20260411114824'),
    ('20260411141831');
