create schema if not exists core;
create schema if not exists ops;
create schema if not exists memory;
create schema if not exists app_template;

create table if not exists core.apps (
  id text primary key,
  name text not null,
  environment text not null check (environment in ('local','staging','prod')),
  owner_bot_id text not null,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now()
);

create table if not exists core.repos (
  id text primary key,
  app_id text not null references core.apps(id),
  git_remote text,
  default_branch text,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now()
);

create table if not exists core.bots (
  id text primary key,
  app_id text not null references core.apps(id),
  display_name text not null,
  persona_id text,
  role text not null,
  status text not null default 'active',
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now()
);

create table if not exists ops.audit_events (
  id uuid primary key default gen_random_uuid(),
  app_id text not null,
  repo_id text,
  canonical_bot_id text not null,
  event_type text not null,
  source_system text not null,
  source_id text,
  payload jsonb not null default '{}'::jsonb,
  created_at timestamptz not null default now()
);

create table if not exists ops.sync_jobs (
  id uuid primary key default gen_random_uuid(),
  app_id text not null,
  repo_id text,
  direction text not null,
  source_system text not null,
  target_system text not null,
  status text not null,
  started_at timestamptz,
  finished_at timestamptz,
  metadata jsonb not null default '{}'::jsonb
);

create table if not exists memory.memory_items (
  id uuid primary key default gen_random_uuid(),
  app_id text not null,
  repo_id text,
  environment text not null,
  canonical_bot_id text not null,
  kind text not null,
  domain text not null,
  source_system text not null,
  source_id text not null,
  content text not null,
  metadata jsonb not null default '{}'::jsonb,
  version int not null default 1,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now()
);

create table if not exists memory.knowledge_items (
  id uuid primary key default gen_random_uuid(),
  app_id text not null,
  repo_id text,
  environment text not null,
  canonical_bot_id text not null,
  category text not null,
  source_system text not null,
  source_id text not null,
  title text,
  content text not null,
  metadata jsonb not null default '{}'::jsonb,
  version int not null default 1,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now()
);
