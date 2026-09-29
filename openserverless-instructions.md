---
name: trustable-app-assistant
description: >-
  Build and fix user-created Trustant apps with React frontend code and Python
  OpenServerless actions. Use this when creating app features, login/register
  flows, CRUD APIs, setup actions, service-backed storage, MCP service checks,
  or debugging app runtime failures inside a Trustant workbench.
---

# Trustant App Assistant Guide

You are working inside a user-created Trustant app. This is a
TypeScript/React frontend plus Python OpenServerless actions. It is not a
conventional backend server project.

## Core Engineering Principles

Make the smallest correct change that fully satisfies the user's request.

Before changing code, understand the existing implementation and preserve its
architecture and behavior unless the requested task requires otherwise.

- Do not refactor unrelated code.
- Do not add features the user did not request.
- Do not introduce abstractions, layers, helpers, files, dependencies, or
  services for hypothetical future needs.
- Prefer local changes over broad rewrites.
- Follow existing project patterns before introducing new ones.
- Do not replace working code merely because another implementation seems
  cleaner.
- Do not fix adjacent issues unless they prevent completion of the requested
  task.
- Validate proportionally to the change: run the smallest set of checks that
  proves the requested behavior works and the modified code remains valid.

## Serverless Operating Model

Core principle: build the app through Trustant/OpenServerless primitives. Do
not replace the platform with hand-written servers, hand-written generated
wrappers, raw credentials, or guessed `ops` commands.

- Frontend code lives in `src/` and calls public actions through
  `/api/my/<package>/<action>`.
- Backend logic lives in editable action modules under
  `packages/<package>/<action>/<module>.py`.
- Generated `__main__.py` files are platform wrappers. Do not create or edit
  them.
- Setup and initialization belong in private actions in package `setup`.
- Trustant launches and manages the Vite dev server and TruACP/Pi process.
- OpenServerless web actions have specific request parameter, metadata, and
  response semantics. Follow the web action rules in this guide rather than
  assuming conventional HTTP server behavior.
- Validate backend changes proportionally to their scope, using the real
  deployed action endpoint when runtime behavior is affected.

## Critical Recovery Contract

Trustant generates `AGENTS.md` in the app root as the app-local mandatory
entrypoint. Treat the Trustant-managed block in `AGENTS.md`,
`.openserverless-contract.md`, and `.mcp.json` as the authoritative project
instruction set. This guidance is embedded directly in `AGENTS.md`;
`CLAUDE.md` contains the same managed block for Claude-compatible agents and
is not an independent source.

Ignore `CONTEXT.md`, `.cursorrules`, `.cursor/rules/*`,
`.github/copilot-instructions.md`, and generated `rules.md` files as mandatory
agent instructions. Inspect them only when the user explicitly asks or when
they are needed to understand legacy template context. They must never
override Trustant-managed instructions or runtime configuration.

Before touching actions, databases, setup, seed data, deploys, or service
state, read `.openserverless-contract.md` if it exists. It is the recovery
contract for the app and takes priority for OpenServerless workflow details.

In a live Trustant Edit session, `ops ide devel` already owns packaging and
deployment. After a coherent backend edit batch, read the canonical watcher
evidence with `trustant_runtime_status`, then run:

```bash
timeout 60 check_openserverless_actions.sh .
```

If validation fails, inspect the reported evidence and fix the underlying
source or tool sequence before retrying. Do not bypass the managed workflow
with manual ZIP handling or raw deploy commands.

If `.openserverless-contract.md` is missing or the checker is unavailable,
fall back to the managed `AGENTS.md` rules and report the limitation rather
than inventing an alternative deployment workflow.

After compaction or context recovery, do not continue editing from memory.
Re-read the active user request, this managed guidance,
`.openserverless-contract.md`, git status, and the relevant project files
before making further changes.

When the user reports a bug or says a previous fix still does not work,
reproduce the reported symptom before editing when reproduction is feasible.
Use the smallest bounded check that demonstrates the failure. If the same
check fails repeatedly, inspect the new evidence and revise the diagnosis
before making another source change.

## Non-Negotiable Rules

- Never create a backend server. Create public or private actions instead.
- Never create or edit generated `__main__.py` files.
- Do not edit `packages/**/__main__.py`, `packages/**/*.zip`, or use raw shell
  commands matching `ops action` / `ops action *`. The managed Pi extension
  enforces the selected workbench boundary; use the OpenServerless MCP action
  tools and the managed watcher/setup flow because those own the generated
  artifacts.
- Never run foreground dev servers or watchers such as `npm run dev`, `vite`,
  or `ops ide devel`.
- Never run `ops ide deploy` during a live Edit session. The existing
  Trustant-managed `ops ide devel` watcher is the sole deploy owner.
- Never kill, restart, or replace Trustant-managed processes. Diagnose the
  existing `http://localhost:5173` server and its generated proxy configuration.
- Never run unbounded commands. Use `timeout <seconds> ...` for checks that may
  hang.
- Never append `|| true` or `|| echo` to deploy, setup, login, checker, or build
  commands, and never pipe those commands through `head` or `tail`; output
  masking can turn a real failure into apparent success.
- Do not ask the user to run shell commands from inside this pod when you have
  shell access. Run bounded checks yourself. Ask the user only when shell/tool
  access is missing or the task requires credentials or physical access only
  the user has.
- Do not build or deploy the Trustant product itself. When validating app
  action changes, use the app deploy/redeploy path described in this guide.
- Put feature logic, request parsing, auth checks, and business behavior in the
  editable module file: `packages/<package>/<action>/<module>.py`.
- After every coherent action MCP/source change batch under `packages/`, wait
  for the managed watcher and run `timeout 60 check_openserverless_actions.sh .`
  before setup, runtime verification, or completion.
- If setup actions change, run `timeout 120 ops ide setup` only after the
  checker confirms the watcher-managed artifacts are current.
- When the requested change affects setup, do not claim completion until the
  required setup succeeds.
- Never create, edit, move, or delete action ZIP files manually. They are
  derived artifacts owned by the managed watcher.
- Do not leave the user with only "try it now" when you can run a bounded
  validation yourself.
- Do not declare a phase complete when the app code path is still failing,
  even if direct MCP or database commands can produce the desired data.
- Do not invent tool or `ops` command names. Use only tools and commands
  exposed by the current Trustant environment.
- Do not use shell redirection to create or replace source files. Avoid
  `cat > file`, heredocs, `tee`, `printf >`, and `sed -i` for app source or
  generated wrappers; use file edit/write tools.
- Before editing a file that changed earlier in the session, re-read its
  current contents. If an edit fails because the expected text is stale or
  already changed, do not retry blindly; re-read the file and adjust the edit
  to its current state.
- Do not write project docs, plans, rules, or examples that recommend forbidden
  commands or invalid endpoint shapes. Documentation must not contain examples
  such as `ops action deploy`, `ops action update`, `v1/auth/register`, or
  `v1/contacts/list`.
- Do not read or copy secrets from `~/.ops/config.json` into app source. Never
  paste database URLs, passwords, tokens, service hosts, buckets, or ports into
  code when platform wiring can provide them.
- If an MCP action tool fails while creating or wiring an action, stop and fix
  that tool sequence. Do not manually create nested action directories,
  generated wrappers, or hardcoded service wiring as a workaround.

## Project Layout

- `src/`: React/TypeScript frontend.
- `public/`: public web assets uploaded automatically.
- `packages/<package>/<action>/`: OpenServerless Python action directories.
- `packages/<package>/<action>/<module>.py`: editable action logic.
- `packages/<package>/<action>/__main__.py`: generated wrapper, do not edit.
- `packages/setup/<action>/`: private setup actions.
- `.agents/skills/`: app-specific skills, when installed.
- `AGENTS.md`: this managed instruction block plus preserved app-local notes.
- `.mcp.json`: authoritative standard MCP configuration for this workbench.

## Application Development Workflow

1. Inspect the existing files and project areas relevant to the requested
   change before editing them.
2. For frontend behavior, work in `src/` and follow the existing
   React/Tailwind patterns.
3. When backend behavior is required, create or update OpenServerless actions
   instead of starting a server process.
4. When the requested behavior requires a platform service, configure it
   through the appropriate action/service tools before writing code that
   depends on it.
5. When initialization is required, put schema, collection, cache, or seed
   initialization in private setup actions.
6. Use generated MCP servers and CLI wrappers when service inspection is
   relevant to implementation or debugging.
7. Validate proportionally to the requested change, using the real public
   endpoint or browser-visible application when runtime behavior is affected.

For frontend changes, validate proportionally to their scope. When runtime
behavior or routing is affected, validate against the Trustant-managed
`http://localhost:5173` using the relevant React validation and bounded HTTP
checks. Check the external `vite.<domain>` ingress only after the managed
watcher has deployed the current sources and only when ingress behavior is in
scope. Do not start another Vite server.

When the requested change involves OpenServerless actions, use this execution
loop:

1. Read `.openserverless-contract.md` if present.
2. Design the action endpoint names and reject invalid nested names before
   creating files.
3. Create actions with the OpenServerless MCP action tool.
4. If an action reads or writes a platform service, immediately add service
   wiring with the matching action/service tool after the action files exist
   and before editing the module logic.
5. Edit only the generated editable module, not `__main__.py`. If the module
   exists without a wrapper, stop and repair/create the action through the MCP
   action tool before continuing.
6. Add Python libraries with `action-requirements` — never with a virtualenv or
   a `requirements.txt`.
7. After a coherent action change batch, wait for the managed `ops ide devel`
   watcher, then run the checker. If setup actions changed, run
   `timeout 120 ops ide setup` only after the checker passes. Inspect failures
   before continuing with runtime validation.

When the requested feature requires backend capabilities, choose among
configured platform services according to these roles:

- Use a public `v1` action for browser-facing APIs.
- Use a private `setup` action for idempotent initialization.
- Use S3 for object/file data.
- Use PostgreSQL for relational data.
- Use Redis for cache, ephemeral state, and authenticated application sessions.
- Use MongoDB for document data only when the official MongoDB capability is
  configured.
- Use Milvus for vector search.
- Do not use Milvus as a replacement for MongoDB.
- Use the Agentic React MCP only when the Vite config imports/references
  `@agentic-react/vite` and invokes `AgenticReact()`.

## OpenServerless Action Tools

When the requested change requires creating or configuring an OpenServerless
action, use the Trustant/OpenServerless MCP action tools instead of manually
creating platform scaffolding. Tool names may appear with hyphens or
underscores depending on the client. Use the matching exposed tool:

- `action-new` / `action_new`: create public or private actions and generated
  wrappers. Repeated creation of a compatible existing action is a successful
  check/no-op; continue without retrying it.
- `action-invoke` / `action_invoke`: invoke private actions such as setup
  actions.
- `action-requirements` / `action_requirements`: add Python libraries.
- `action-add-secret` / `action_add_secret`: add an environment secret.
- `action-add-s3` / `action_add_s3`: add S3 service wiring.
- `action-add-postgresql` / `action_add_postgresql`: add PostgreSQL service
  wiring.
- `action-add-redis` / `action_add_redis`: add Redis service wiring.
- `action-add-milvus` / `action_add_milvus`: add Milvus service wiring.
- `action-add-mongodb` / `action_add_mongodb`: add MongoDB service wiring.
- `secret-unbind` / `secret_unbind`: remove an obsolete generated secret
  binding atomically without reading or deleting its value. Use it to repair a
  legacy invalid managed-variable binding; never edit `__main__.py` manually.
  After a removal, let the managed watcher update the changed endpoint. If
  previously deployed parameters remain, report that a Trustant-owned full
  redeploy is required; do not race the watcher with raw deploy commands.

If a tool call returns "Invalid Tool", stop and use only an exposed tool name.
Do not retry with guessed aliases. If a shell command reports
`no command named ...`, do not keep guessing `ops` subcommands; use the
documented MCP action tools instead.

After changing an action, let the managed watcher produce the derived
artifacts and run the action checker before runtime validation. If a setup
action changed, run `ops ide setup` only after the checker passes. Never use
raw `ops action` commands or manual ZIP manipulation as an alternative deploy
path.

For a new public HTTP endpoint, use package `v1` unless the user explicitly
asks for another package. The endpoint is reachable at
`/api/my/<package>/<action>`.

For initialization, create private actions in package `setup` with
`public: false`. After creating or changing setup actions, wait for the managed
watcher, run the checker, and then run `ops ide setup`.

## Action Endpoint Grammar

OpenWhisk action names are namespace/package/action. In Trustant app code, the
MCP action endpoint must therefore be only:

- `action`
- `package/action`

Every package and action segment must start with a letter and contain only
letters, numbers, and hyphens. Use flat hyphenated names such as
`v1/employees-photo`; underscores and spaces are invalid. Never send
`v1/employees_photo` to an action MCP tool.

For browser-facing APIs, use package `v1`. Valid examples:

- `v1/register`
- `v1/login`
- `v1/me`
- `v1/contacts`
- `v1/orders`
- `setup/database`

Invalid examples:

- `v1/auth/register`
- `v1/contacts/list`
- `v1/orders/create`
- `packages/v1/auth/register`

For a new CRUD API, prefer one public action per resource, such as
`v1/contacts` or `v1/orders`, and branch inside the editable module using
`__ow_method` plus request data. If separate actions are clearer, keep names
flat and hyphenated, such as `v1/contacts-list` or `v1/orders-create`.

For an existing API, preserve its current valid endpoint structure unless the
requested change requires restructuring it.

Never create nested directories under `packages/<package>/<group>/<action>` to
simulate routes. They are not valid Trustant/OpenServerless endpoints.

## MCP Servers And Service Access

`.mcp.json` is regenerated at launch with an `mcpServers` object. Pi exposes
those servers through one lazy `mcp` proxy tool; use `mcp({})` to inspect
server status and `mcp({ server: "<server>" })` to list a server's tools.

A missing direct tool name or an MCP process that has not started does not
mean the server is absent. MCP servers start lazily on first use. Before
reporting a configured service as missing, check both the `mcp` proxy and the
keys of `.mcp.json.mcpServers`.

Available servers depend on the current Trustant configuration:

- `openserverless`: always present; exposes action-management tools.
- `react`: always present; read-only deterministic source validation for the
  current React/Vite workbench. Use `react_project_inspect`,
  `react_validate_routes`, `react_validate_auth_flow`, and the aggregate
  `react_validate` when relevant to the requested change.
- `agentireact`: present only when `vite.config.js` or `vite.config.ts`
  imports/references `@agentic-react/vite` and invokes `AgenticReact()` in
  executable config code. Comments, strings, wrong packages, and the obsolete
  `AgentiReact()` spelling do not enable it. Adding the plugin during a live
  session requires relaunching the app so Trustant regenerates `.mcp.json`.
- `s3`: present only when S3 is configured; companion CLI wrapper: `rclone`.
- `postgres`: present only when PostgreSQL is configured; companion CLI
  wrapper: `psql`.
- `redis`: present only when Redis is configured; companion CLI wrapper:
  `redis-cli`.
- `milvus`: present only when Milvus is configured; companion CLI wrapper:
  `milvus_cli`.
- `mongodb`: present only when MongoDB is configured as an official
  OpenServerless capability in `~/.ops/config.json`.

Service MCP access and action runtime access are separate capabilities. A
successful MCP call proves only that the assistant can inspect the service; it
does not prove that an application action can access it.

For action runtime access, use the matching OpenServerless action service tool
and its generated `ctx` binding, such as `action-add-redis`,
`action-add-postgresql`, `action-add-s3`, `action-add-milvus`, or
`action-add-mongodb`. If no matching runtime binding exists, do not invent a
backend connection.

Service MCP servers are generated from the Trustant/OpenServerless
configuration after `ops ide login`. Do not hardcode service hosts, ports,
credentials, bucket names, database names, tokens, or other connection details
when platform wiring provides them.

MongoDB is a document database capability and is separate from Milvus/vector
search. Never use Milvus as a substitute for MongoDB.

If the requested behavior requires MongoDB but the official MongoDB capability
or runtime binding is unavailable, expose a deterministic
`non configurato`/error state in the affected app path and report the
limitation. Do not invent connection details or MongoDB runtime environment
variables.

When `action-add-mongodb` / `action_add_mongodb` is available, use it before
writing module code that requires MongoDB. The generated wrapper exposes
`ctx.MONGODB_CLIENT` and `ctx.MONGODB`; editable modules should use those
context values rather than connection strings.

`MDB_MCP_CONNECTION_STRING` belongs only to the generated MongoDB MCP server
process and must not be used by application source, action modules, generated
wrappers, documentation, or frontend code.

You may inspect `~/.ops/config.json` only to determine which services are
configured. Never copy values from it into application code, wrapper code,
logs, documentation, or frontend configuration.

Service MCP servers are discovery and verification tools. Do not use direct
service mutation operations to create schema, seed records, repair application
state, or substitute for the application's setup and public action paths,
unless the user explicitly requests an administrative data repair.

For normal application work, put idempotent initialization in setup actions
and application writes in public OpenServerless actions. Read-only MCP checks
are appropriate when they provide relevant implementation or validation
evidence.

When inspecting PostgreSQL through MCP, use the exact exposed tool names. Use
`postgres_list_schemas` for schemas and `postgres_list_objects` for
tables/views in a schema. Do not invent generic tool names such as
`list_schemas`.

## Application Environment

Application `.env` and `.env.production` files are immutable agent boundaries.
Never read, create, edit, import, synchronize, or regenerate them. Only the
user may change application environment values through Trustant's
configuration interface.

If a required application environment variable is missing, report the missing
variable and stop the affected path rather than inventing or persisting a
value.

## Redis Authentication

When the requested application requires authentication, use Redis-backed
opaque sessions. Do not build application authentication on JWT or an
application signing secret.

For authentication flows:

- create or update only the authentication and protected-resource actions
  required by the requested flow, then call `auth-setup` / `auth_setup` once
  with the complete token, protected/session, and logout endpoint sets used by
  that flow. It atomically adds Redis wiring without reading or writing
  `.env`;
- use `action-add-redis` / `action_add_redis` only for an individual
  non-authentication endpoint that needs Redis;
- generate a cryptographically random opaque token at login or registration;
- store the token-to-identity mapping only in Redis with a bounded TTL;
- construct every session key from `ctx.REDIS_PREFIX`;
- validate the Redis record on every protected request and derive the current
  user from it, never from a browser-supplied user id;
- delete the Redis record during logout when the requested flow includes
  logout;
- return and persist only the opaque token in the browser.

When using Redis in an action, first add Redis wiring with
`action-add-redis` / `action_add_redis`. The generated wrapper exposes
`ctx.REDIS` and `ctx.REDIS_PREFIX` and adds the required Python Redis client to
the action package.

Always construct Redis keys from the generated prefix and an app-local suffix:

```python
def redis_key(ctx, name):
    return f"{getattr(ctx, 'REDIS_PREFIX', '') or ''}{name}"
```

Use `ctx.REDIS.get(redis_key(ctx, "cache:item"))`,
`ctx.REDIS.set(redis_key(ctx, "cache:item"), value)`, and the same pattern for
other Redis operations.

Do not use naked Redis keys directly with `ctx.REDIS`; Trustant/Nuvolaris Redis
ACLs only allow the configured user prefix.

## Runtime Host Rules

Pi and TruACP run inside the Trustant environment. Classify hosts before using
them:

- `localhost:5173` is the pod-local app dev server started by `ops ide devel`.
  Use it for app HTTP validation from this shell.
- `localhost:4096` is the local TruACP server.
- `trustant.<domain>` is the browser-visible Trustant UI/API host.
- `vite.<domain>` is the browser-visible app host through the Trustant
  proxy/ingress. Use it only after managed deployment is confirmed and only
  when external browser or ingress routing is in scope.
- `opencode.<domain>` is the legacy browser-visible hostname that proxies
  TruACP for ingress and WAF compatibility.
- `OPS_APIHOST` is the configured OpenServerless API host for Trustant and
  `ops ide` orchestration. It must never be bound into an action, exposed as
  `ctx.OPS_APIHOST`, read by an action module, or emitted as
  `#--param OPS_APIHOST "$OPS_APIHOST"`.

Do not rewrite the generated Vite `/api/my` proxy target. Trustant starts the
managed dev server with the current app's `OPSDEV_HOST`; browser application
code must continue to use relative `/api/my/...` URLs.

Frontend code calls actions with relative `/api/my/<package>/<action>` URLs so
the browser keeps its current origin.

Action modules must not call sibling actions through `OPS_APIHOST`,
browser-visible hosts, or ingress URLs. When server-side aggregation is
required, prefer using the required generated service bindings directly
within the responsible action rather than chaining sibling actions over HTTP.

One action may use MongoDB and Milvus independently. Use
`ctx.MONGODB_CLIENT` / `ctx.MONGODB` for document-database operations and
`ctx.MILVUS` for vector operations. Do not remove or replace either legitimate
service use merely because both appear in the same business module.

Do not invent pod IPs, raw service names, public domains, or replacement
localhost URLs.

When runtime endpoint validation is required, use the managed local app host:

```bash id="5gz4qi"
curl http://localhost:5173/api/my/<package>/<action>
```

## PostgreSQL Action Pattern

When an action requires PostgreSQL access, first add PostgreSQL wiring with
`action-add-postgresql` / `action_add_postgresql`. The generated wrapper
provides the configured client as `ctx.POSTGRESQL`.

Use the generated client directly from the editable action module:

```python
def main(args, ctx):
    db = ctx.POSTGRESQL
    rows = db.execute("SELECT id, email FROM users ORDER BY id")
    return {
        "statusCode": 200,
        "body": {"users": rows},
    }
```

Do not open a new database connection manually and do not hardcode database
connection details.

Schema creation, migrations, and seed data belong in private setup actions,
not in public application actions.

## Skills

When an installed app-specific skill is relevant to the requested task, inspect
the corresponding `.agents/skills/<skill>/SKILL.md` and follow its
instructions.

Do not inspect or load unrelated skills.

## Web Action Request Rules

OpenServerless web actions expose request data through action arguments and
`__ow_*` metadata. Parse only the request inputs required by the action. Do not
copy every helper below into an action unless that action needs it.

When an action accepts a JSON request body, use a small local helper that
handles the supported OpenServerless body representations:

```python
import base64
import json


def body_json(args):
    body = args.get("__ow_body")
    if body is None:
        return {}

    if args.get("__ow_isBase64Encoded"):
        body = base64.b64decode(body).decode("utf-8")

    if isinstance(body, dict):
        return body

    if isinstance(body, str):
        try:
            return json.loads(body)
        except json.JSONDecodeError:
            return {}

    return {}
```

When an action accepts query parameters, read them from the OpenServerless
request arguments rather than assuming a conventional web framework request
object:

```python
from urllib.parse import parse_qs


def query_params(args):
    query = args.get("__ow_query")
    if not query:
        return {}

    parsed = parse_qs(query, keep_blank_values=True)
    return {
        key: values[-1] if values else ""
        for key, values in parsed.items()
    }
```

When an action needs an HTTP request header, read it from the OpenServerless
request metadata:

```python
def header_value(args, name):
    headers = args.get("__ow_headers") or {}
    wanted = name.lower()

    for key, value in headers.items():
        if key.lower() == wanted:
            return value

    return None
```

Keep small request-parsing helpers local to the action module unless the
existing project already provides an appropriate shared utility. Do not
introduce shared abstractions solely to avoid a small amount of local parsing
code.

## Web Action Response Rules

OpenWhisk web actions can use top-level `headers`, `statusCode`, and `body` as
HTTP response instructions. Trustant-generated Python wrappers, however, may
call the editable module and return:

```python
{"body": module.main(args, ctx=ctx)}
```

Because of that, a module return value such as:

```python
{"statusCode": 401, "body": {"error": "Token non fornito"}}
```

may reach the browser as HTTP 200 with that object nested inside JSON when the
generated wrapper does not pass web-action envelopes through.

Therefore:

- Do not edit `__main__.py` to force HTTP status behavior.
- Treat editable module return values as application JSON unless the generated
  wrapper is known to pass web-action envelopes through.
- When the wrapper does not pass envelopes through, use simple application
  payloads such as `{"ok": False, "error": "..."}` for app-level errors.
- Preserve the existing response contract of working endpoints unless the
  requested change requires changing it.
- When wrapped responses must be handled by the frontend, follow the existing
  project's response-normalization pattern if one exists.

If the frontend does not already normalize wrapped responses, use a minimal
pattern such as:

```ts
const raw = await response.json();
const data = raw && typeof raw === "object" && "body" in raw ? raw.body : raw;

if (!response.ok || data?.ok === false || data?.error) {
  throw new Error(data?.error || `Request failed: ${response.status}`);
}
```

Do not add response-normalization code to unrelated frontend paths that already
handle their response contract correctly.

## Browser-Opened And Printable Actions

When the frontend opens an action URL directly with `window.open(...)`, an
`<a>` link, or a form target, the endpoint must return a browser-native
response. Returning application JSON that contains HTML does not make the
endpoint directly browser-renderable.

Preserve an existing working browser or print flow unless the requested change
requires changing it.

For printable HTML such as invoices, receipts, labels, reports, or documents,
use the response strategy supported by the generated wrapper and the existing
application flow.

If the generated wrapper passes web-action envelopes through to OpenWhisk, the
action may return a real HTML response:

```python id="87p2oh"
return {
    "statusCode": 200,
    "headers": {"Content-Type": "text/html; charset=utf-8"},
    "body": html,
}
```

When implementing or fixing this direct browser-opened flow, verify the actual
HTTP response with a bounded check such as:

```bash id="5pjv66"
curl -i http://localhost:5173/api/my/v1/<action>/<id>...
```

The response must have the expected HTML content type and return HTML rather
than application JSON.

If the generated wrapper nests module output as application JSON, keep the
action response as JSON and let the frontend perform the authenticated fetch.
Extract the returned HTML, open a new window, write the HTML into that window,
and print from there. Do not treat a raw
`window.open("/api/my/...")` call as proof that this flow works.

For example, this is not a directly printable browser response:

```python id="z3k4sw"
return {"ok": True, "html": html}
```

Do not put application session tokens in query parameters when an authenticated
fetch can be used instead.

If the existing application explicitly requires a direct browser-opened URL
and no header-based authenticated flow is possible, use only its existing
short-lived application-session mechanism. Never use OpenServerless platform
credentials, including authentication values from `~/.ops/config.json`, as
application session tokens.

## Authentication UI Rules

Apply these rules when the requested application flow includes authentication.

Implement only the authentication UI and routes required by the application.
Do not add registration, logout, profile, recovery, or other authentication
flows unless they are required by the requested behavior or already exist in
the application.

- Public authentication routes must remain accessible without an authenticated
  session.
- Protected content and navigation must not be shown as authenticated merely
  because cached user data exists in the browser.
- Protect direct routes as well as navigation links. If an unauthenticated user
  opens a protected route directly, redirect to the appropriate authentication
  route or render the application's unauthenticated state.
- When the application is intended to be fully protected, the root route must
  resolve to the appropriate authentication or authenticated application state
  rather than leaving the Trustant starter screen visible.
- Replace starter placeholder content only where it conflicts with the
  requested application flow. Do not refactor unrelated screens solely to
  remove starter content.
- With React Router `HashRouter`, pass logical routes such as `/login` to
  `Link`, `NavLink`, `Navigate`, and `useNavigate`. The router adds `#/` to the
  browser URL. Do not pass `#/login` to those APIs or use root-relative
  anchors such as `<a href="/login">` for internal HashRouter navigation.
- After successful login, persist only the opaque session token required by
  the frontend. Cached user or profile data may be used for rendering but is
  not authoritative proof of authentication.
- A successful login must update the live authentication provider or store
  before navigating to a protected route. Persisting token or user data only
  to `localStorage` must not leave the current React state unauthenticated.
- When restoring a persisted authenticated session on page load, keep an
  explicit loading state while the backend validates the token through the
  application's session-validation endpoint. Do not render protected content
  from cached browser identity before validation succeeds.
- If persisted-session validation fails, clear the invalid session state and
  render the appropriate unauthenticated flow.
- After successful registration, follow the authentication behavior required
  by the application. Do not add automatic login, an additional login step, or
  a redirect unless required by the requested flow.
- When logout is part of the requested flow, clear both the backend session
  according to the application's Redis session contract and the corresponding
  frontend authentication state.
- Give form controls stable `id` and `name` values and associate labels with
  matching `htmlFor` values. Do not use placeholders as substitutes for
  labels.
- Do not hardcode browser-visible identity such as `user_id=1` in fetch URLs or
  request bodies. Protected backend actions must derive the current user from
  validated authenticated request state rather than trusting a user id supplied
  by the browser.

Use the Redis authentication contract defined earlier for every authentication
or protected action that participates in the requested flow. Do not create
additional authentication endpoints solely to complete a predefined
login/register/session/logout set.

## Setup And Data Initialization

Persistent application initialization belongs in private actions in package
`setup`.

Create or modify only the setup actions required by the requested change and
the services actually used by the application. Do not create setup actions for
unused services or for hypothetical future needs.

Setup actions must be incremental, idempotent, and non-destructive. They must
preserve existing application data unless the user explicitly requests a data
reset or destructive migration.

Use the setup action appropriate to the required initialization:

- PostgreSQL schema and table creation belong in `setup/database`.
- Redis key or cache initialization belongs in `setup/cache`.
- Milvus collection creation belongs in `setup/collection`.
- MongoDB collection and index preparation belongs in an idempotent setup
  action only when MongoDB is configured and required by the application.
- Private S3 data preload belongs in `setup/upload`.
- Public web assets belong in `public/`, not in setup uploads.

Do not perform schema creation, collection creation, persistent cache
initialization, or seed-data initialization from public application actions.

Do not add seed, demo, or example data unless it is required by the requested
application behavior.

When setup actions are created or changed, follow the managed watcher,
checker, and `ops ide setup` workflow defined earlier in this guide.

## Dependencies

Add dependencies only when they are required by the requested implementation
and the existing project does not already provide the needed capability.

For a required frontend dependency that is not already present, add it to
`package.json` and use the project's existing npm workflow.

For Python action dependencies, use only
`action-requirements` / `action_requirements`. Never create or edit an action
`requirements.txt` manually.

Before importing a non-standard-library Python package, determine whether it is
actually required by the requested implementation. If it is required, pass it
to `action-requirements` and let the tool determine whether the library is
already provided by the OpenServerless runtime or must be added to the action's
generated `requirements.txt`.

Never create a Python virtual environment with `python -m venv`, `virtualenv`,
`uv venv`, or a `.venv`/`venv` directory, and do not install, vendor, download,
or otherwise introduce Python action dependencies outside the
`action-requirements` workflow.

If an action fails with `ModuleNotFoundError`, determine whether the import is
required. If it is, process the missing dependency through
`action-requirements`; otherwise correct or remove the invalid import. Do not
continue runtime validation while a required import is unresolved.

Do not add, upgrade, replace, or reinstall unrelated dependencies while
implementing a requested change.

## Data And Service Restrictions

Apply these restrictions to the services involved in the requested change.
Do not inspect or configure unrelated services merely to verify these
restrictions.

Retrieve the current user with `ops util whoami` only when the user-specific
service name or namespace is needed.

- PostgreSQL database is named after the user; the default schema is
  `<user>_schema`.
- Milvus database is named after the user.
- MongoDB may be used only when the official MongoDB capability is configured.
- Redis keys used by actions must be built with the generated
  `ctx.REDIS_PREFIX`. Do not guess `<user>:` manually and do not use naked
  Redis keys.
- S3 writable application buckets are `<user>-data` for private app data and
  `<user>-web` for public web assets.
- The S3 MCP cannot list buckets. Treat `<user>-data` and `<user>-web` as the
  writable application buckets defined by the platform; do not probe for or
  invent additional writable buckets.

`OPS_USER`, `OPS_PASSWORD`, `OPS_APIHOST`, `OPS_REPO`, and `OPS_SKILLS` are
Trustant-managed orchestration variables, not application secrets. Never bind
or expose them as application secrets or action runtime values, and never copy
their values into application source, frontend code, or documentation.

## Validation Checklist

Validate proportionally to the requested change. Run the smallest set of
checks that proves the requested behavior works and that the modified code
remains valid. Do not expand validation into unrelated features or services.

For OpenServerless action changes:

- After a coherent action change batch, wait for the managed watcher and run
  `timeout 60 check_openserverless_actions.sh .` when the checker is
  available.
- When setup actions changed, run `ops ide setup` only after the checker
  passes. Required setup must succeed before setup-related work is complete.
- When runtime behavior changed, validate the affected public action with a
  bounded HTTP check against
  `http://localhost:5173/api/my/<package>/<action>`.
- If an action reports `Cannot start action`, `application error`, or
  `developer error`, inspect the bounded action logs before changing the
  implementation strategy.
- An invocation that returns only an activation id is not proof of application
  behavior. When behavior must be verified, inspect the result or validate the
  relevant public endpoint.

For request and API behavior:

- Validate the HTTP methods, request fields, headers, route parameters, and
  response shapes affected by the requested change.
- For CRUD work, validate the operations affected by the requested change.
  When implementing or restructuring a complete CRUD resource, validate the
  complete create/list/update/delete flow.
- When REST-style item routes are part of the change, validate the actual item
  path such as `PUT /api/my/v1/<resource>/<id>` or
  `DELETE /api/my/v1/<resource>/<id>` rather than proving only a body-based
  fallback.
- When frontend code consumes a changed backend response, verify that it
  handles the response shape actually returned.

For frontend behavior:

- Validate the routes and UI states affected by the requested change using the
  relevant React validation and bounded checks against the Trustant-managed
  app.
- For authentication changes, validate only the authentication and protected
  states involved in the requested flow, including direct access to affected
  protected routes when relevant.
- For browser-opened or printable endpoints, verify the actual HTTP status,
  content type, and body shape when those properties are part of the requested
  behavior.
- Use `vite.<domain>` only when external browser or ingress behavior is in
  scope and managed deployment is confirmed.

For service-backed behavior:

- If live service state was modified during debugging, verify that the
  corresponding application or setup source can reproduce the required state.
  Runtime state alone is not completion proof.
- When the requested change affects S3 read/write behavior, validate it through
  the configured application action path. If a direct read/write probe is
  needed, use a unique temporary object in `ctx.S3_DATA`, read it back, compare
  its contents, and delete it afterward. Do not use bucket listing as proof of
  read/write access.

For source validation:

- Treat diagnostics in files changed by the task as failures when they indicate
  a real problem in the modified code.
- Do not fix unrelated pre-existing diagnostics, warnings, or lint issues
  unless they prevent validation of the requested change.
- Use bounded checks and let failures remain visible. Do not hide failures with
  forced successful exits or output truncation that can mask the relevant
  error.

If the required validation cannot be performed with the available tools or
environment, report the specific blocker and the validation that remains
unproven rather than claiming completion.
