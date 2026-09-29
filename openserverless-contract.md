# OpenServerless Action Contract

This file is the short recovery contract for Trustant/OpenServerless app work.
Read it before modifying actions, databases, setup, seed data, deployment, or
service state.

If `check_openserverless_actions.sh .` reports contract or source drift,
re-read this file before editing again.

Treat the Trustant-managed instructions in `AGENTS.md`, this contract, and
`.mcp.json` as authoritative for Trustant/OpenServerless work.

Do not treat `CONTEXT.md`, `.cursorrules`, `.cursor/rules/*`,
`.github/copilot-instructions.md`, or generated `rules.md` files as mandatory
agent instructions. Inspect them only when the user explicitly requests it or
when they are needed to understand relevant legacy/template context. They must
never override the authoritative Trustant/OpenServerless instructions.

Before using an MCP server, discover its actual exposed tools and schemas.
Discover only the servers relevant to the requested task; do not enumerate or
connect to unrelated services.

Use the exact tool names and schemas exposed by the active MCP configuration.
Do not guess tool names, parameters, or contracts from memory.

## Environment

- You are inside a generated Trustant app workbench, normally
  `/home/trustant/workbench/<app>`.
- The durable source of the application is its Git repository. Fix application
  behavior in repository files, not only in live runtime state.
- When validation or diagnosis requires a pod-local command and shell access is
  available, run the bounded command yourself rather than asking the user to
  run it.
- After context compaction, re-read the active user request, this contract,
  `AGENTS.md`, git status, and the relevant project files before resuming work.
- When the user reports a bug or says a previous fix does not work, reproduce
  the relevant symptom before modifying source when reproduction is possible.
  Use bounded HTTP, log, or deterministic checks. If the same check fails
  repeatedly without new evidence, stop repeating it and revise the diagnosis.
- After source changes, validate proportionally to the requested change. Run
  only the checks relevant to the files and behavior modified, and use
  `git diff --check` as a final source-integrity check. For user-visible
  frontend changes, verify the affected route or behavior.
- `ops ide devel` exposes the managed application inside the pod at
  `http://localhost:5173` and owns live action packaging and deployment.
  Never start another `ops ide devel`, `vite`, or `npm run dev` process, and
  never run `ops ide deploy` during the managed Edit session.
- Do not kill, restart, or replace Trustant-managed development processes.
- Browser/ingress hosts such as `vite.<domain>` are for external verification.
  Use them only when external browser or ingress behavior is relevant and the
  current sources have been deployed by the managed workflow.
- `OPS_APIHOST` is a Trustant/OpenServerless orchestration value, not an
  application secret or action runtime parameter. Do not bind it into an
  action, expose it through `ctx`, or read it from application action code.
- Browser code must call actions with relative
  `/api/my/<package>/<action>` URLs so it preserves the current origin.
- Actions must not call sibling actions through `OPS_APIHOST`, browser-visible
  hosts, or ingress URLs. Let the frontend call independent endpoints when
  appropriate, or give a server-side action the service bindings required to
  perform its work directly.
  
## Files

- Frontend source lives in `src/`.
- Public web assets live in `public/`.
- Editable action logic lives in
  `packages/<package>/<action>/<module>.py`.
- Generated `__main__.py` wrappers are platform-owned. Never create or edit
  them manually; create or repair actions and service wiring through the
  OpenServerless MCP tools.
- Generated action ZIP archives are platform-owned deployment artifacts, not
  editable application source.
- Setup and initialization actions live under
  `packages/setup/<action>/`.

## Action Names

Valid OpenServerless action names have one of these forms:

- `action`
- `package/action`

Each package and action segment must start with a letter and contain only
letters, numbers, and hyphens. Use flat hyphenated action names; underscores,
spaces, and additional path segments are invalid.

Use package `v1` as the default for newly created browser-facing APIs unless
the existing application uses another package or the user explicitly requests
otherwise.

Do not rename existing valid actions merely to match a preferred naming style.

Valid examples:

- `v1/prospect`
- `v1/issues`
- `v1/employees-photo`
- `setup/database`

Invalid examples:

- `v1/auth/register`
- `v1/contacts/list`
- `v1/employees_photo`
- `packages/v1/auth/register`

## Data, Authentication, And Service Rules

- Application `.env` and `.env.production` are immutable agent boundaries.
  Never read, create, edit, import, synchronize, or regenerate them. Report a
  required missing application variable instead of creating it.

- When the requested application requires authentication, use Redis-backed
  opaque sessions, not JWT or an application signing secret. Create only the
  authentication/session endpoints required by the requested flow and
  configure the participating authentication and protected endpoints through
  `auth-setup` / `auth_setup`.

- Generate opaque random session tokens, store token-to-identity mappings in
  Redis with a bounded TTL, validate the session on protected requests, and
  delete the session on logout when logout is part of the requested flow.
  Authentication/session Redis keys must use `ctx.REDIS_PREFIX`. The backend
  derives identity from the session and never trusts a browser-supplied user
  id.

- Trustant-managed `OPS_*` orchestration variables are not application
  secrets. Do not pass them to application secret tools or generated secret
  bindings.

- Service MCP servers are assistant-side discovery and diagnostic tools. Their
  presence does not create runtime environment variables, action parameters,
  or `ctx` bindings. Add required runtime service wiring through the
  corresponding OpenServerless action/service tool.

- For PostgreSQL actions, use generated PostgreSQL wiring and
  `ctx.POSTGRESQL`. Do not create a separate connection from environment
  variables when the generated binding is available.

- MongoDB may be used only when the official MongoDB capability is configured.
  When required, use generated MongoDB action wiring and its official
  `ctx.MONGODB_CLIENT` / `ctx.MONGODB` binding. Do not use MCP-private
  connection variables or invent MongoDB runtime environment variables.

- If a required service is not configured, report the missing capability
  rather than inventing connection details or substituting another service.

- For Redis actions, use generated Redis wiring through `ctx.REDIS` and
  construct every application key from `ctx.REDIS_PREFIX` plus an app-local
  suffix. Never use naked Redis keys.

- For S3 actions, use generated S3 wiring and its scoped application buckets.
  Never use `ctx.S3_CLIENT.list_buckets()` for service discovery.

- Do not hardcode database URLs, service hosts, users, passwords, schemas,
  bucket names, or service ports when platform wiring provides them.

- For transactional database writes, commit the intended transaction before
  reporting success.

- Create seed data only when required by the requested application behavior.
  Seed operations must be idempotent and must not duplicate existing data.

- Database migrations must be repeatable and non-destructive unless the user
  explicitly requests a destructive migration. Use idempotent operations such
  as `IF NOT EXISTS` where supported.

- When the shape of a derived database object such as a view must change,
  update it through the application's setup/migration path.

- Do not treat a live database or service-state repair as an application fix.
  If live state is inspected or repaired during debugging, put the equivalent
  reproducible change in the appropriate application/setup source and validate
  that source path.
  
## Web Action Route IDs

- The OpenServerless MCP tools create actions and service wiring; they do not
  replace this runtime contract.
- For item routes such as `/api/my/v1/contacts/123`, do not assume one fixed
  `__ow_path` shape.
- Use a helper that accepts body fallback and suffix forms such as `123`,
  `/123`, `/contacts/123`, and `/api/my/v1/contacts/123`.
- A `PUT` fix that works only because the frontend sends `id` in the JSON body
  is incomplete if `DELETE /api/my/v1/<resource>/<id>` still fails.
- For each CRUD resource, test create/list/update/delete through
  `http://localhost:5173`, including `PUT` and `DELETE` with id in the URL.

## Browser And Printable Responses

- If the frontend opens an action URL with `window.open(...)`, the target must
  be a browser response, not just app JSON.
- For printable HTML such as an invoice, return HTML as the HTTP body with
  `Content-Type: text/html; charset=utf-8`. Do not return
  `{"ok": true, "html": "<!DOCTYPE html>..."}` when the browser opens the URL
  directly.
- If the generated wrapper nests module returns under JSON and cannot pass
  headers/body through, use a frontend route that fetches JSON with
  `Authorization`, extracts the HTML, writes it to a new window/document, and
  then prints. Do not pretend that raw `window.open(/api/my/...)` will render
  embedded JSON HTML as a page.
- For opened/downloaded/printable URLs, verify with `curl -i` from inside the
  pod and check both status and `Content-Type`.
- Token-in-query is acceptable only when a new browser window cannot send the
  `Authorization` header; prefer short-lived or app-session tokens and validate
  with a real session token.

## Deploy And Verification

- After one or more successful `action_new` creations, finish the coherent
  action/wiring/source batch and call `trustant_runtime_redeploy` exactly
  once. It invokes the same safe Trustant host workflow as the UI Redeploy
  action: stop the watcher, run the full deploy, restart the watcher, and wait
  for readiness. A compatible idempotent `action_new` no-op does not require
  it. Do not run a concurrent `ops ide deploy`.
- After the required redeploy succeeds, read the canonical watcher state with
  `trustant_runtime_status`. Watcher status, checker, HTTP, and browser
  verification are blocked while redeploy remains required.
- Run `timeout 60 check_openserverless_actions.sh .` once after that evidence
  and before setup, runtime verification, or completion. In managed live mode,
  the checker validates source and contract invariants without treating sibling
  ZIP existence or freshness as watcher state.
- If setup actions changed, run `timeout 120 ops ide setup` only after watcher
  evidence shows a successful action update and the checker passes. Do not
  claim completion until setup runs successfully.
- Never create or update action ZIP files manually. The managed watcher owns
  the sibling `packages/<package>/<action>.zip` artifacts.
- Never inspect, list, search, stat, or poll those sibling ZIPs. If
  `trustant_runtime_status` reports an error or no progress, use that exact
  evidence to repair the source/tool sequence or report a managed failure. Do
  not repeat the checker without a relevant mutation or watcher change, run a
  manual deploy, increase the timeout, or repair a ZIP.
- If the checker reports an action module without `__main__.py`, create or
  repair that action with the OpenServerless MCP action tool before editing the
  module logic.
- If the checker reports wrapper drift, do not patch `__main__.py` by hand:
  recreate or repair the action/service wiring with the OpenServerless MCP
  tools, then edit only the module file.
- Verify app HTTP endpoints from inside the pod with:
  `curl http://localhost:5173/api/my/<package>/<action>`.
- For write paths, write and then read back the changed value.
- For delete paths, delete through the public HTTP route and then confirm the
  record is no longer returned.
- If you used `psql` or a service MCP to inspect/repair live data during
  debugging, also prove the source setup/action code recreates the same state.
- For printable/browser-opened paths, prove the response shape with
  `curl -i http://localhost:5173/...` and check that JSON endpoints return JSON
  while printable HTML endpoints return `text/html`.
- Use `vite.<domain>` only after managed deployment is confirmed and only for
  external browser/ingress verification.
- Do not hide failures with `|| true` or output truncation that masks the first
  actionable error.
- Do not use `|| true`, `|| echo`, or `head`/`tail` pipelines on deploy, setup,
  login, checker, and frontend-build commands; output masking can turn a real
  failure into apparent success.

## If Blocked

- If an MCP action tool is missing or returns invalid tool, call `mcp({})` and
  inspect generated `.mcp.json` plus that server's exact tool list; do not
  invent raw `ops action` commands.
- After three semantically equivalent failures with no successful relevant
  source or wiring mutation, stop that strategy and revise the diagnosis.
  There is no numeric global step/turn budget while work makes real progress.
- If a service MCP write would "fix" state, fix the setup/action code instead
  unless the user explicitly requested administrative data repair.
- If validation is impossible, state the blocker and the command that failed.
