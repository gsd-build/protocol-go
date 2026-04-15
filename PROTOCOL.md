# GSD Cloud Wire Protocol

Version: 1
Transport: WebSocket (text frames, JSON payloads)

This document is the **authoritative source** for the GSD Cloud relay protocol.
Both the Go bindings in this repository and any TypeScript bindings must match
this contract exactly.

## Envelope

Every message is a JSON object with a `type` field:

```json
{ "type": "<name>", ...fields }
```

## Browser → Daemon messages

### `task`
Dispatch a user message to a session.

| Field | Type | Notes |
|---|---|---|
| type | "task" | |
| taskId | uuid | |
| sessionId | uuid | |
| channelId | string | Routes stream events back to the correct browser tab |
| prompt | string | |
| model | string | e.g. `claude-opus-4-6[1m]` |
| effort | "low" \| "medium" \| "high" \| "max" | |
| permissionMode | string | e.g. `acceptEdits` |
| personaSystemPrompt | string? | Injected via `--append-system-prompt` |
| cwd | string | Absolute path on the daemon's machine |
| claudeSessionId | string? | Pass to `claude -p --resume` to continue an existing Claude conversation. Empty for the first turn. |
| requestId | uuid? | Optional root correlation ID for request-scoped logging. |
| traceparent | string? | W3C trace context. |

### `stop`
Interrupt the current Claude process for a session.

| Field | Type |
|---|---|
| type | "stop" |
| channelId | string |
| sessionId | uuid |

### `permissionResponse`

| Field | Type |
|---|---|
| type | "permissionResponse" |
| channelId | string |
| sessionId | uuid |
| requestId | uuid |
| approved | boolean |

### `questionResponse`

| Field | Type |
|---|---|
| type | "questionResponse" |
| channelId | string |
| sessionId | uuid |
| requestId | uuid |
| answer | string |

### `browseDir`

| Field | Type |
|---|---|
| type | "browseDir" |
| requestId | uuid |
| channelId | string |
| machineId | uuid |
| path | string |

### `readFile`

| Field | Type |
|---|---|
| type | "readFile" |
| requestId | uuid |
| channelId | string |
| machineId | uuid |
| path | string |
| maxBytes | int? | Defaults to 512 KiB |

### `mkDir`

| Field | Type |
|---|---|
| type | "mkDir" |
| requestId | uuid |
| channelId | string |
| machineId | uuid |
| path | string |

## Daemon → Browser messages

### `stream`
High-frequency Claude event. The `event` field is an opaque JSON object passed through from Claude's stream-json output.

| Field | Type |
|---|---|
| type | "stream" |
| sessionId | uuid |
| channelId | string |
| sequenceNumber | int64 |
| event | object |
| requestId | uuid? | Optional root correlation ID for request-scoped logging. |
| traceparent | string? | W3C trace context. |

### Delivery semantics

Live WebSocket delivery is best-effort. The relay forwards frames immediately
and persists session-scoped history to Postgres on a separate path. Browser
reconnect recovery happens by reloading persisted messages by session and
sequence from the database. There is no daemon ↔ relay ack/replay/WAL handshake
in protocol version 1.

### `taskStarted`

| Field | Type |
|---|---|
| type | "taskStarted" |
| taskId | uuid |
| sessionId | uuid |
| channelId | string |
| startedAt | iso8601 string |
| requestId | uuid? | Optional root correlation ID for request-scoped logging. |
| traceparent | string? | W3C trace context. |

### `taskComplete`

| Field | Type |
|---|---|
| type | "taskComplete" |
| taskId | uuid |
| sessionId | uuid |
| channelId | string |
| claudeSessionId | string |
| inputTokens | int64 |
| outputTokens | int64 |
| costUsd | string | Decimal string to avoid float precision loss |
| durationMs | int |
| resultSummary | string? |
| requestId | uuid? | Optional root correlation ID for request-scoped logging. |
| traceparent | string? | W3C trace context. |

### `taskError`

| Field | Type |
|---|---|
| type | "taskError" |
| taskId | uuid |
| sessionId | uuid |
| channelId | string |
| error | string |
| requestId | uuid? | Optional root correlation ID for request-scoped logging. |
| traceparent | string? | W3C trace context. |

### `taskCancelled`
Sent when the user interrupts a running task via `stop`.

| Field | Type |
|---|---|
| type | "taskCancelled" |
| taskId | string |
| sessionId | uuid |
| channelId | string |
| requestId | uuid? | Optional root correlation ID for request-scoped logging. |
| traceparent | string? | W3C trace context. |

### `permissionRequest`

| Field | Type |
|---|---|
| type | "permissionRequest" |
| sessionId | uuid |
| channelId | string |
| requestId | uuid |
| toolName | string |
| toolInput | object |

### `question`

| Field | Type |
|---|---|
| type | "question" |
| sessionId | uuid |
| channelId | string |
| requestId | uuid |
| question | string |
| header | string? |
| multiSelect | boolean? |
| options | { label: string, description?: string, preview?: string }[]? |

### `heartbeat`

| Field | Type |
|---|---|
| type | "heartbeat" |
| machineId | uuid |
| daemonVersion | string |
| status | "online" |
| timestamp | iso8601 string |

### `browseDirResult`

| Field | Type |
|---|---|
| type | "browseDirResult" |
| requestId | uuid |
| channelId | string |
| ok | boolean |
| entries | []BrowseEntry? |
| error | string? |

`BrowseEntry`:
```json
{ "name": "...", "path": "...", "isDirectory": bool, "size": int, "modifiedAt": "iso8601" }
```

### `readFileResult`

| Field | Type |
|---|---|
| type | "readFileResult" |
| requestId | uuid |
| channelId | string |
| ok | boolean |
| content | string? |
| truncated | boolean? |
| error | string? |

### `mkDirResult`

| Field | Type |
|---|---|
| type | "mkDirResult" |
| requestId | uuid |
| channelId | string |
| ok | boolean |
| error | string? |

## Daemon ↔ Relay control messages

### `hello` (daemon → relay, first frame after connect)

| Field | Type |
|---|---|
| type | "hello" |
| machineId | uuid |
| daemonVersion | string |
| os | string |
| arch | string |
| activeTasks | string[]? | Task IDs the daemon still considers in flight |

### `welcome` (relay → daemon, response to hello)

| Field | Type |
|---|---|
| type | "welcome" |
| latestDaemonVersion | string? | Optional latest daemon version for update prompts |

### `syncCrons` (relay → daemon)

Sent after daemon connect and after cron config mutations so the daemon can
reconcile its local cron config directory to the server-owned config set.

| Field | Type |
|---|---|
| type | "syncCrons" |
| machineId | uuid |
| jobs | CronSpec[] |
| sentAt | iso8601 string |

`CronSpec`:

| Field | Type |
|---|---|
| id | uuid |
| name | string |
| cronExpression | string |
| prompt | string |
| mode | string |
| model | string |
| effort | string |
| projectId | uuid |
| targetSessionId | uuid? |
| enabled | boolean |

### `cronInventory` (daemon → relay)

Sent after daemon reconciliation and on future local inventory changes so the
relay/browser can show the daemon's current cron state.

| Field | Type |
|---|---|
| type | "cronInventory" |
| machineId | uuid |
| items | CronLocalState[] |
| timestamp | iso8601 string |

`CronLocalState`:

| Field | Type |
|---|---|
| id | uuid |
| name | string |
| cronExpression | string |
| enabled | boolean |
| syncedAt | iso8601 string |
| lastRunAt | iso8601 string? |
| nextRunAt | iso8601 string? |
| locallyModified | boolean |

### `cronExecResult` (daemon → relay)

Sent after a daemon-managed cron execution finishes so the relay can persist the
result and fan it out to browser subscribers.

| Field | Type |
|---|---|
| type | "cronExecResult" |
| machineId | uuid |
| cronJobId | uuid |
| taskId | uuid |
| status | string |
| error | string? |
| durationMs | int64? |
| timestamp | iso8601 string |

### `skillInventory` (daemon → relay)

Sent after connect and whenever the daemon's discovered skill inventory changes.

| Field | Type |
|---|---|
| type | "skillInventory" |
| machineId | uuid |
| entries | `SkillInventoryEntry[]` |

`SkillInventoryEntry`:

| Field | Type | Notes |
|---|---|---|
| slug | string | Stable skill identifier |
| displayName | string | Human-readable name |
| description | string | Short summary |
| scope | string | `project`, `global`, or `installed` |
| runtime | string | Discovery runtime such as `claude` or `codex` |
| root | string | Absolute root directory |
| projectRoot | string? | Project root for project-scoped skills |
| relativePath | string | Path from `root` to the skill file |
| sourceKind | string | Discovery source classification |
| machineFingerprint | string | Machine-side content + metadata hash used for inventory sync |
| editable | boolean | True only for cloud-managed global skills |

### `skillContentRequest` (relay → daemon)

Requests the body for one managed skill file on a machine.

| Field | Type |
|---|---|
| type | "skillContentRequest" |
| machineId | uuid |
| slug | string |
| root | string |
| relativePath | string |

### `skillContentUpload` (daemon → relay)

Returns the file body for a requested managed skill.

| Field | Type | Notes |
|---|---|---|
| type | "skillContentUpload" |
| machineId | uuid |
| slug | string |
| root | string |
| relativePath | string |
| content | string | Raw file contents |
| machineFingerprint | string | Machine fingerprint for the uploaded file |
| baseCloudRevision | int64 | Cloud revision the machine based the upload on |

### `skillPush` (relay → daemon)

Pushes the cloud version of a managed skill back to a daemon.

| Field | Type | Notes |
|---|---|---|
| type | "skillPush" |
| machineId | uuid |
| slug | string |
| root | string |
| relativePath | string |
| content | string | Raw file contents |
| cloudFingerprint | string | Cloud fingerprint for the pushed file |
| cloudRevision | int64 | Cloud revision being applied |

### `skillDelete` (relay → daemon)

Deletes a managed skill from an allowed managed-global root.

| Field | Type | Notes |
|---|---|---|
| type | "skillDelete" |
| machineId | uuid |
| slug | string |
| root | string |
| relativePath | string |
| cloudRevision | int64 | Cloud revision that authorized the delete |

`skillContentRequest`, `skillContentUpload`, `skillPush`, and `skillDelete`
are only valid for editable managed-global skills. Project and installed skills
stay inventory-only and should never receive a push or delete command.
