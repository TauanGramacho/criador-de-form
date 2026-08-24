# Form Builder Challenge

Aplicacao full stack para criacao e publicacao de formularios, criada para o desafio tecnico Full Stack Developer 2026.

## Stack

- Backend: Go, Huma, SQLite local, Postgres/Supabase em producao, cookie sessions, Google OAuth2, bcrypt.
- Frontend: React, TypeScript, Vite, MUI, React Router, TanStack Query.
- API contract: OpenAPI gerada pelo backend.
- Client TypeScript: gerado a partir de `api/openapi.json` com `@hey-api/openapi-ts`.

## Requisitos Locais

- Go 1.26 ou superior.
- Node.js 24 ou superior.
- npm.

SQLite nao precisa ser instalado como servidor. A aplicacao ja esta integrada com banco de dados SQLite e usa um arquivo local em `backend/data/formbuilder.db`.

Em producao, use Supabase Postgres com uma `DATABASE_URL` no formato:

```text
postgresql://postgres:SUA-SENHA@db.gptvpmdacbhgnyisvzsl.supabase.co:5432/postgres?sslmode=require
```

## Backend

```bash
cd backend
go run ./cmd/server migrate
go run ./cmd/server run
```

Por padrao, o backend sobe em `http://localhost:8080`.

Endpoints uteis:

- API: `http://localhost:8080/api`
- OpenAPI JSON: `http://localhost:8080/openapi.json`
- Docs interativas: `http://localhost:8080/docs`

## Frontend

```bash
cd frontend
npm install
npm run generate:client
npm run dev
```

Por padrao, o frontend sobe em `http://localhost:5173`.

## Configuracao

O backend aceita configuracao por arquivo, variaveis de ambiente e argumentos.

Arquivo:

```bash
cd backend
go run ./cmd/server run --config=./config.example.yaml
```

Argumentos:

```bash
cd backend
go run ./cmd/server run --address=localhost:8080 --database-url="file:./data/formbuilder.db?_pragma=foreign_keys(1)"
```

Variaveis de ambiente:

```bash
ADDRESS=localhost:8080 DATABASE_URL="file:./data/formbuilder.db?_pragma=foreign_keys(1)" go run ./cmd/server run
```

No PowerShell:

```powershell
$env:ADDRESS="localhost:8080"
$env:DATABASE_URL="file:./data/formbuilder.db?_pragma=foreign_keys(1)"
go run ./cmd/server run
```

Principais opcoes:

- `ADDRESS`: endereco do backend.
- `DATABASE_URL`: arquivo SQLite.
- `FRONTEND_ORIGIN`: origem aceita no CORS, padrao `http://localhost:5173`.
- `FRONTEND_BASE_URL`: base usada para gerar links publicos.
- `COOKIE_SECURE`: use `true` em HTTPS.
- `GOOGLE_CLIENT_ID`: client ID OAuth do Google.
- `GOOGLE_CLIENT_SECRET`: client secret OAuth do Google.
- `GOOGLE_REDIRECT_URL`: callback OAuth, padrao `http://localhost:8080/api/auth/google/callback`.

## Configurando Login Com Google

O fluxo Google OAuth esta implementado, mas precisa de credenciais reais para funcionar.

No Google Cloud Console:

1. Crie ou selecione um projeto.
2. Configure a OAuth consent screen.
3. Crie uma credencial do tipo OAuth Client ID para Web application.
4. Adicione este redirect URI autorizado:

```text
http://localhost:8080/api/auth/google/callback
```

5. Configure o backend com as credenciais:

```powershell
$env:GOOGLE_CLIENT_ID="seu-client-id"
$env:GOOGLE_CLIENT_SECRET="seu-client-secret"
$env:GOOGLE_REDIRECT_URL="http://localhost:8080/api/auth/google/callback"
go run ./cmd/server run
```

Se voce baixou o JSON do Google, tambem pode rodar:

```powershell
cd backend
.\run-google-local.ps1 -ClientSecretJson "C:\caminho\para\client_secret.json"
```

Se essas variaveis nao estiverem configuradas, o frontend desabilita o botao de login com Google e mostra um aviso.

Para producao na URL Render sugerida, configure no Google Cloud:

Origens JavaScript autorizadas:

```text
https://criador-de-form.onrender.com
```

URIs de redirecionamento autorizados:

```text
https://criador-de-form.onrender.com/api/auth/google/callback
```

Links de consentimento:

```text
https://criador-de-form.onrender.com
https://criador-de-form.onrender.com/privacy
https://criador-de-form.onrender.com/terms
```

## Banco E Migrations

As migrations ficam em `backend/migrations`.

O banco usado pela aplicacao e criado automaticamente em:

```text
backend/data/formbuilder.db
```

Esse arquivo guarda usuarios, sessoes, formularios publicados e respostas. Ele nao deve ser versionado no Git.

Para Supabase, a migration equivalente fica em:

```text
supabase/migrations/20260824152000_init_form_builder.sql
```

Para aplicar localmente via backend:

```bash
cd backend
go run ./cmd/server migrate
```

Para aplicar no projeto Supabase `gptvpmdacbhgnyisvzsl`:

```bash
npx supabase login
npx supabase link --project-ref gptvpmdacbhgnyisvzsl
npx supabase db push --linked
```

Se preferir sem linkar o projeto, use `npx supabase db push --project-ref gptvpmdacbhgnyisvzsl --password "SUA-SENHA-DO-BANCO"`.

O comando `run` tambem aplica migrations na inicializacao para melhorar a experiencia local.

## OpenAPI

A especificacao OpenAPI e gerada a partir dos tipos e handlers do backend Go.

Para atualizar o arquivo versionado:

```bash
cd backend
go run ./cmd/server openapi > ../api/openapi.json
```

## Client TypeScript

O frontend nao mantem manualmente o client da API. Depois de atualizar a OpenAPI:

```bash
cd frontend
npm run generate:client
```

Os arquivos gerados ficam em `frontend/src/api/generated` e nao devem ser editados manualmente.

## Fluxo Principal

1. Acesse `http://localhost:5173`.
2. Crie uma conta com e-mail e senha ou configure Google OAuth.
3. Crie um formulario.
4. Configure campos.
5. Publique o formulario.
6. Abra o link publico `/f/:slug`.
7. Envie uma resposta.
8. Veja as respostas na area administrativa.

## Decisoes De Arquitetura

- SQLite foi escolhido para reduzir dependencias externas e garantir execucao local multiplataforma. A camada de persistencia esta isolada em `Store` e tambem suporta Postgres/Supabase via `DATABASE_URL`.
- Huma foi escolhido porque gera OpenAPI a partir do backend, reduzindo risco de contrato desatualizado.
- Sessoes usam cookie `HttpOnly` e tokens persistidos por hash no banco.
- A definicao dos campos e as respostas ficam como JSON no SQLite. Isso simplifica o modelo para um form builder dinamico e mantem validacao no backend.
- O frontend usa TanStack Query para estado remoto e o SDK TypeScript gerado para todas as chamadas de API comuns.

## Limitacoes E Trade-offs

- Nao ha versionamento de formularios publicados. Editar um formulario publicado muda a validacao das proximas respostas.
- Google OAuth depende de credenciais externas configuradas pelo avaliador.
- Nao ha controle de equipes ou multiplos administradores por formulario.
- O armazenamento de respostas e simples; exportacao CSV seria uma evolucao natural.

## Verificacao

Backend:

```bash
cd backend
go test ./...
```

Frontend:

```bash
cd frontend
npm run build
npm audit --omit=dev
```

## Deploy

O arquivo `render.yaml` publica a aplicacao como um unico servico Docker no Render. O Dockerfile compila o frontend, compila o backend Go e o backend serve tanto a API quanto a SPA React.

Variaveis obrigatorias no Render:

- `DATABASE_URL`: connection string do Supabase Postgres com senha real e `sslmode=require`.
- `GOOGLE_CLIENT_ID`: client ID OAuth do Google.
- `GOOGLE_CLIENT_SECRET`: client secret OAuth do Google.
- `GOOGLE_REDIRECT_URL`: `https://criador-de-form.onrender.com/api/auth/google/callback`.
- `FRONTEND_BASE_URL`: `https://criador-de-form.onrender.com`.
- `FRONTEND_ORIGIN`: `https://criador-de-form.onrender.com`.
- `COOKIE_SECURE`: `true`.
