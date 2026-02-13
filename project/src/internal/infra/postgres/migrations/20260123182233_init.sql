-- +goose Up
-- +goose StatementBegin
create type public.position_type as enum ('ADMIN', 'MANAGER', 'CUSTOMER', 'DEV');
create table if not exists public.roles(
    id serial primary key,
    role public.position_type not null default 'CUSTOMER',
    description varchar(100) not null,
    created_at timestamp not null default now(),
    updated_at timestamp not null default now(),
    deleted_at timestamp
);
insert into public.roles (id, role, description)
values (
        1,
        'CUSTOMER',
        'cliente do sistema'
    );
insert into public.roles (id, role, description)
values (
        2,
        'MANAGER',
        'gerente do sistema'
    );
insert into public.roles (id, role, description)
values (
        3,
        'ADMIN',
        'administrador do sistema'
    );
insert into public.roles (id, role, description)
values (
        4,
        'DEV',
        'desenvolvedor do sistema'
    );
create table if not exists public.users(
    id serial primary key,
    uuid uuid not null,
    name varchar(50) not null,
    email varchar(100) not null,
    password varchar(255) not null,
    role_id int not null default 1,
    enabled boolean not null default true,
    created_at timestamp not null default now(),
    updated_at timestamp not null default now(),
    deleted_at timestamp
);
ALTER TABLE users
ADD CONSTRAINT unique_user_email UNIQUE (email);
-- CONSTRAINTS
-- user -> roles
alter table if exists public.users
add constraint fk_users_role_id foreign key (role_id) references public.roles(id) on update cascade on delete cascade;
insert into public.users (uuid, name, email, "password", role_id)
values (
        '4a9b3fd5-6813-4c75-9598-5fd9ae202d88',
        'Admin',
        'admin@email.com',
        '$2a$10$zYC48a1doguo1VoCqbmQBezAUQJKVSbGnHgoPWInNFn2idbPABUoe',
        3
    ) on conflict (email) do nothing;
insert into public.users (uuid, name, email, "password", role_id)
values (
        '296446de-e045-4638-a4e5-a09e94136fee',
        'Jonas',
        'jonas.w.martins@gmail.com',
        '$2a$10$zYC48a1doguo1VoCqbmQBezAUQJKVSbGnHgoPWInNFn2idbPABUoe',
        4
    ) on conflict (email) do nothing;
create table if not exists public.links(
    id serial primary key,
    uuid uuid not null,
    data varchar(5000) not null,
    expires_at timestamp not null,
    created_at timestamp not null default now(),
    updated_at timestamp not null default now(),
    deleted_at timestamp
);
CREATE INDEX idx_links_uuid ON public.links(uuid);
create type public.link_type as enum ('RESET_PASS', 'OTHER');
create table if not exists public.user_available_links(
    id serial primary key,
    user_id int not null,
    link_uuid uuid not null,
    type link_type not null,
    created_at timestamp not null default now(),
    updated_at timestamp not null default now(),
    deleted_at timestamp
);
-- CONSTRAINTS
DROP INDEX IF EXISTS idx_links_uuid;
ALTER TABLE public.links
ADD CONSTRAINT idx_unique_links_uuid UNIQUE (uuid);
-- user_available_links -> users
alter table if exists public.user_available_links
add constraint fk_user_available_links_users_id foreign key (user_id) references public.users(id) on update cascade on delete cascade;
-- user_available_links -> links
alter table if exists public.user_available_links
add constraint fk_user_available_links_link_uuid foreign key (link_uuid) references public.links(uuid) on update cascade on delete cascade;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table IF EXISTS public.users;
drop table IF EXISTS public.roles;
drop type IF EXISTS public.position_type;
drop table if not exists public.links;
DROP INDEX IF EXISTS idx_links_uuid;
drop type if exists public.link_type cascade;
drop table if exists public.user_available_links;
DROP INDEX IF EXISTS idx_unique_links_uuid;
-- +goose StatementEnd
