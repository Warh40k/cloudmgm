BEGIN;


CREATE TABLE IF NOT EXISTS public.users
(
    id uuid primary key,
    name character varying(255) NOT NULL,
    username character varying(255) NOT NULL,
    password_hash character varying NOT NULL,
    balance numeric NOT NULL DEFAULT 0,
    created timestamp without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS public.storages
(
    id uuid,
    title character varying(255) NOT NULL,
    geolocation character varying NOT NULL,
    total_size bigint NOT NULL default 0,
    created timestamp without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS public.volumes
(
    id uuid primary key,
    storage_id uuid references public.storages(id) on delete cascade ,
    label character varying(255) NOT NULL,
    description character varying(1024),
    size bigint NOT NULL DEFAULT 0
);

CREATE TABLE IF NOT EXISTS public.files
(
    id uuid primary key ,
    volume_id uuid NOT NULL references volumes(id) on delete cascade ,
    name character varying NOT NULL,
    link character varying NOT NULL,
    created timestamp without time zone NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS public.users_volumes
(
    id uuid primary key ,
    user_id uuid NOT NULL references public.users(id) on delete  cascade,
    volume_id uuid NOT NULL references public.volumes(id) on delete cascade
);

CREATE TABLE IF NOT EXISTS public.volumes_storages
(
    id uuid primary key ,
    volume_id uuid NOT NULL references public.volumes(id) on delete cascade ,
    storage_id uuid NOT NULL references public.storages(id) on delete cascade
);

END;