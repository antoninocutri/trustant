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

OpenServerless web actions are Apache OpenWhisk web actions:

- Public web actions can be invoked over HTTP without an OpenWhisk API key.
- The action owner pays for the activation, so the action must implement its
  own application-level authorization when needed.
- Query parameters, form fields, and JSON object body fields can be passed as
  first-class action arguments.
- In normal OpenWhisk merging, body fields override query fields.
- HTTP context is exposed through reserved metadata keys such as
  `__ow_method`, `__ow_headers`, and `__ow_path`.
- Web actions support HTTP methods such as GET, POST, PUT, PATCH, DELETE, HEAD,
  and OPTIONS. Use `__ow_method` for method-based CRUD actions.
- Requests cannot override reserved `__ow_*` metadata names.

Trustant-generated Python actions should be defensive: some wrappers or
clients may also provide `args["body"]` as a dict or JSON string. Merge both
shapes and let top-level fields win, because a generated wrapper or previous
edit can create an empty `body = {}` while real request fields are top-level.

Use this pattern in editable modules when reading JSON fields:

```python
import json

def request_data(args):
    data = dict(args) if isinstance(args, dict) else {}
    body = data.get("body")
    if isinstance(body, str):
        try:
            body = json.loads(body)
        except Exception:
            body = {}
    merged = dict(body) if isinstance(body, dict) else {}
    ignored = {"body", "POSTGRES_URL", "__ow_method", "__ow_headers", "__ow_path"}
    merged.update({k: v for k, v in data.items() if k not in ignored})
    return merged
```

Read request metadata from OpenServerless keys first:

```python
def request_method(args):
    return (args.get("__ow_method") or args.get("method") or "GET").upper()

def request_headers(args):
    headers = args.get("__ow_headers") or args.get("headers") or {}
    return {str(k).lower(): v for k, v in headers.items()} if isinstance(headers, dict) else {}

headers = request_headers(args)
auth_header = headers.get("authorization", "")
```

If a raw or non-JSON request body is needed, handle `__ow_body` explicitly.
Most app JSON endpoints should not need raw body handling.

For REST-style item routes, do not assume `__ow_path` always contains the full
public URL. It can be a suffix or a different shape depending on the
OpenServerless web action route. Use body `id` only as a fallback, not as the
only way update/delete works.

Use this pattern or an equivalent one for item ids:

```python
def request_route_id(args, data, resource_name):
    for key in ("id", f"{resource_name}_id"):
        value = data.get(key)
        if value not in (None, ""):
            return str(value)

    raw_path = str(args.get("__ow_path") or args.get("path") or "").strip("/")
    if not raw_path:
        return ""

    parts = [part for part in raw_path.split("/") if part]
    if not parts:
        return ""

    if resource_name in parts:
        index = parts.index(resource_name)
        if index + 1 < len(parts):
            return parts[index + 1]

    return parts[-1]
```

For CRUD resources, test both update and delete through the public HTTP path:

```bash
curl -X PUT http://localhost:5173/api/my/v1/<resource>/<id> ...
curl -X DELETE http://localhost:5173/api/my/v1/<resource>/<id> ...
```

A test that only calls `/api/my/v1/<resource>` with `{"id": ...}` in the body
does not prove the REST-style item route works.

## Web Action Response Rules

OpenWhisk web actions can use top-level `headers`, `statusCode`, and `body` as
HTTP response instructions. Trustant-generated Python wrappers, however,
commonly call the editable module and return:

```python
{ "body": module.main(args, ctx=ctx) }
```

Because of that, a module return value such as:

```python
{"statusCode": 401, "body": {"error": "Token non fornito"}}
```

can reach the browser as HTTP 200 with that object nested inside JSON if the
wrapper did not pass it through.

Therefore:

- Do not edit `__main__.py` just to force HTTP status behavior.
- Treat editable module return values as application JSON unless the generated
  wrapper is known to pass web-action envelopes through.
- Prefer simple module payloads such as `{"ok": False, "error": "..."}` for
  app-level errors.
- Frontend fetch code should normalize both direct and wrapped payloads before
  reading fields.

Use this frontend pattern or an equivalent one:

```ts
const raw = await response.json();
const data = raw && typeof raw === "object" && "body" in raw ? raw.body : raw;
if (!response.ok || data?.ok === false || data?.error) {
  throw new Error(data?.error || `Request failed: ${response.status}`);
}
```

## Browser-Opened And Printable Actions

If the frontend opens an action URL directly with `window.open(...)`, an `<a>`
link, or a form target, the endpoint must return a browser-native response. Do
not return JSON that contains HTML and then claim the browser flow is complete.

For printable HTML such as invoices, receipts, labels, reports, or documents,
the correct behavior is one of these:

1. The action returns a real HTML response:

```python
return {
    "statusCode": 200,
    "headers": {"Content-Type": "text/html; charset=utf-8"},
    "body": html,
}
```

This only works if the generated wrapper passes the envelope through to
OpenWhisk. Verify with:

```bash
curl -i http://localhost:5173/api/my/v1/<action>/<id>...
```

The response must show `Content-Type: text/html`, and the body must begin with
HTML, not with JSON.

2. If the wrapper nests module output as application JSON, keep the action JSON
and change the frontend flow: fetch with `Authorization`, extract `data.html`,
open a new window, write the HTML into that window, and then call print. In
that case do not use raw `window.open("/api/my/...")` as proof that printing
works.

Avoid this broken pattern for direct browser-opened endpoints:

```python
return {"ok": True, "html": html}
```

That renders as JSON in a new browser window. It is not a printable page.

Token-in-query is acceptable only when a new window cannot send the
`Authorization` header. Prefer short-lived app-session tokens, and validate with
a real session token from the app database or login flow. Do not use the
OpenServerless `~/.ops/config.json` auth value as an app session token.

## Authentication UI Rules

When an app has login or registration:

- Treat login/register as the only public UI flows.
- Replace starter placeholder screens. The root route must redirect to login,
  render login, or render the authenticated app based on session state; it must
  not keep the Trustant starter/welcome template.
- Before marking auth UI complete, inspect the router and the component used by
  `/` or `#/`. Remove or replace generated starter content such as `Welcome`,
  `Try the following prompts to start`, `Powered by Trustant`, `trustant.png`,
  or sample prompt lists. A protected app is incomplete if the browser-visible
  home page still shows the starter screen.
- Hide protected navigation items such as dashboards, contacts, orders,
  settings, admin, or profile until the user is authenticated.
- Protect direct routes too. If an unauthenticated user opens a protected hash
  route directly, redirect to the login/register route or render the auth view,
  not the protected page.
- With React Router `HashRouter`, pass logical routes such as `/login` to
  `Link`, `NavLink`, `Navigate`, and `useNavigate`. The router adds `#/` to the
  browser URL. Never pass `#/login` to those APIs, and never use root-relative
  anchors such as `<a href="/login">` for internal navigation.
- After login, persist only the opaque session token needed by the frontend.
  A cached user/profile may improve rendering but is never authoritative proof
  of authentication.
- A successful login or registration must update the live authentication
  provider/store before navigating to a protected route. Writing token/user
  data only to `localStorage` leaves the current React render unauthenticated
  and commonly causes an immediate redirect back to login.
- On every full-page load, keep an explicit authentication loading state and
  validate the persisted token through a bounded backend `me`/session endpoint
  before rendering protected routes. The backend validates expiry and derives
  identity from the token. On any validation failure, clear the token and
  cached identity and render the public authentication flow; never fall back to
  a cached localStorage user as an authenticated session.
- A successful registration must establish the same authenticated session as
  login, either by returning session/token/user data directly or by performing
  an immediate login. Do not send the newly registered user back to a separate
  login step before entering the protected area.
- Add an explicit logout path when protected navigation is shown.
- Back login, registration, `me`/session, every protected action, and logout
  with the same Redis session contract. Every one of those actions must have
  generated Redis wiring and use `ctx.REDIS_PREFIX`; JWT and application
  signing secrets are not an alternative.
- Give every form control a stable `id` and `name`, and associate each label
  with `htmlFor` matching that `id`. A placeholder is not a label. Ambiguous or
  duplicated control identity is a form-semantics bug: fix it at the source
  rather than working around it.
- Do not hardcode browser-visible identity such as `user_id=1` in fetch URLs or
  request bodies. The backend must derive the current user from authenticated
  request state, such as a token/session header, not from a user id supplied by
  the browser.

## Setup And Data Initialization

- All initialization belongs in private actions in package `setup`.
- Setup actions must be incremental, idempotent, and non-destructive.
- Table creation belongs in `setup/database`.
- Redis key preparation belongs in `setup/cache`.
- Milvus collection creation belongs in `setup/collection`.
- MongoDB collection/index preparation belongs in an idempotent setup action
  only when MongoDB is configured.
- Private S3 data preload belongs in `setup/upload`.
- Public web assets belong in `public/`, not in setup uploads.
- Wait for the managed watcher and run the action checker after creating or
  changing actions, then run `ops ide setup` when setup actions changed.
- `ops ide setup` must succeed before setup work is complete.
- If setup returns `Cannot start action. Check logs for details.`, immediately
  run `timeout <seconds> ops logs --last` and fix the first traceback. Do not
  proceed by mutating the service directly.
- Do not create missing tables or seed rows with PostgreSQL MCP write tools and
  then claim setup succeeded. The `setup/*` action must be able to recreate the
  state idempotently.
- Do not make a live DB-only schema fix with `psql`, PostgreSQL MCP, or ad hoc
  SQL and then claim the app is fixed. If you inspect or repair live state while
  debugging, put the equivalent idempotent migration in `setup/database`, run
  `ops ide setup`, and read back the schema/data through a bounded command.

Examples of idempotent setup:

- `CREATE TABLE IF NOT EXISTS ...`
- `ALTER TABLE ... ADD COLUMN IF NOT EXISTS ...`
- Redis `SET ... NX` for keys that should only be seeded once.
- Check Milvus collection existence before creating it.
- Check MongoDB collection/index existence before creating it.
- Upload S3 objects only when missing or changed.

## Dependencies

- Add frontend dependencies to `package.json`, then run `npm install`.
- Add Python dependencies only with `action-requirements`.
- Never create a virtualenv (`python -m venv`, `virtualenv`, `uv venv`, or any
  `.venv`/`venv` directory) and never create or edit a `requirements.txt`.
  Actions are built and deployed server-side, so a local virtualenv is never
  used at runtime — it only leaves artifacts that Clean has to remove.
  `action-requirements` is the only supported path.
- Before importing a non-stdlib Python package such as `bcrypt`, `jwt`,
  `requests`, or a database driver, add it with `action-requirements` and
  redeploy the action.
- If action logs show `ModuleNotFoundError`, fix the dependency or import before
  doing any other validation. Do not mark the feature complete.
- Add PostgreSQL, Redis, S3, Milvus, MongoDB, and secrets with the corresponding
  action/service tool. Do not hardcode credentials and do not manually edit
  generated wrapper code.
- `OPS_USER`, `OPS_PASSWORD`, `OPS_APIHOST`, `OPS_REPO`, and `OPS_SKILLS` are
  Trustant-managed orchestration variables, not application secrets. Never
  pass them to `action-add-secret`, `secret-bind`, `secret-ensure`, or
  `auth-setup`; a secret tool must reject them.

## Data And Service Restrictions

Retrieve the current user with `ops util whoami` when needed. These
restrictions are enforced by the platform:

- PostgreSQL database is named after the user; the default schema is
  `<user>_schema`.
- Milvus database is named after the user.
- MongoDB is available only when the official post-login config exposes a
  MongoDB block or derived connection string.
- Redis keys used by actions must be built with the generated
  `ctx.REDIS_PREFIX`; do not guess `<user>:` manually and do not use naked keys.
- S3 writable buckets are `<user>-data` for private app data and `<user>-web`
  for public web assets.
- The S3 MCP cannot list buckets, so assume only the two user buckets above are
  writable.

## Validation Checklist

End backend-related work with proof:

- After changing an action module, wait for the managed watcher to settle.
- Run `timeout 60 check_openserverless_actions.sh .` after that wait when the
  checker is available. On stale archives, perform only one bounded wait and
  recheck before reporting a watcher failure.
- After changing setup actions, wait for the watcher, pass the checker, and then
  run `ops ide setup`.
- Validate public actions with bounded HTTP checks against
  `http://localhost:5173/api/my/<package>/<action>` from inside this pod.
- For CRUD resources, validate the full create/list/update/delete matrix. Test
  `PUT /api/my/v1/<resource>/<id>` and
  `DELETE /api/my/v1/<resource>/<id>` without relying only on `id` in the JSON
  body, then read back to confirm the updated value or deleted absence.
- For browser-opened or printable endpoints, validate with `curl -i` and prove
  the response status and content type match the browser use case. A direct
  `window.open("/api/my/...")` target for printable HTML must not return
  `application/json`.
- Use `vite.<domain>` only after managed deployment is confirmed and only for
  explicit external browser/ingress checks.
- `ops action invoke` by itself is not enough proof when it only prints an
  activation id such as `ok: invoked ...`; inspect the action result/logs or
  validate through the HTTP endpoint.
- If any action reports `Cannot start action`, `application error`, or
  `developer error`, run `timeout <seconds> ops logs --last` before changing
  strategy.
- If `psql` or a service MCP was used to inspect or repair live database state,
  prove the source setup/action code recreates that state. Runtime state alone
  is not completion proof.
- Verify JSON request fields, method, and headers are visible to the action.
- Verify frontend fetch handling accepts the response shape actually returned.
- Treat editor, LSP, TypeScript, lint, and tool diagnostics as validation
  failures when they mention generated or edited files. Fix the diagnostic, or
  explain why it is stale with a successful bounded command that proves it.
- Use bounded checks such as `timeout <seconds> ...` and `curl`.
- Do not hide validation failures with `|| true`, forced zero exits, or output
  truncation that can mask the first error. Let checks fail loudly, then fix the
  failure.
- For frontend auth flows, validate both route shape and route behavior: root
  path, login path, register path, direct protected route while logged out, and
  protected navigation after login. After submitting valid credentials, assert
  that the protected page is visibly rendered without requiring a reload.
- If validation is impossible, state the blocker instead of asking the user to
  "try it now" with no local proof.
