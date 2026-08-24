create table if not exists public.users (
  id text primary key,
  email text not null unique,
  name text not null,
  password_hash text,
  google_id text unique,
  created_at text not null,
  updated_at text not null
);

create table if not exists public.sessions (
  token_hash text primary key,
  user_id text not null references public.users(id) on delete cascade,
  expires_at text not null,
  created_at text not null
);

create table if not exists public.forms (
  id text primary key,
  owner_id text not null references public.users(id) on delete cascade,
  title text not null,
  description text not null,
  slug text not null unique,
  status text not null check (status in ('draft', 'published')),
  fields_json text not null,
  created_at text not null,
  updated_at text not null,
  published_at text
);

create table if not exists public.responses (
  id text primary key,
  form_id text not null references public.forms(id) on delete cascade,
  answers_json text not null,
  submitted_at text not null,
  submitter_ip text not null default ''
);

create index if not exists idx_forms_owner_id on public.forms(owner_id);
create index if not exists idx_forms_slug on public.forms(slug);
create index if not exists idx_responses_form_id on public.responses(form_id);
create index if not exists idx_sessions_user_id on public.sessions(user_id);

alter table public.users enable row level security;
alter table public.sessions enable row level security;
alter table public.forms enable row level security;
alter table public.responses enable row level security;

