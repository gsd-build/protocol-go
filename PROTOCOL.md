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

Receivers bound WebSocket text frames before unmarshalling them into typed
payloads. The Go bindings expose `ParseEnvelopeWithLimits` and
`ValidateEnvelopeFrame` for frame size, JSON nesting depth, object field count,
and array item count validation. `ParseEnvelope` is the raw envelope parser.

Request-scoped handlers bind responses to the originating request and session
using `requestId`, `sessionId`, and `channelId` where those fields exist. The Go
bindings expose `ExtractBinding`, `ValidateRequestBinding`, and
`ValidateSessionBinding` for these checks.

## Browser → Daemon messages

### `task`
Dispatch a user message to a session.

`engine` selects the daemon task executor. Current daemons run tasks through the Pi executor.

| Field | Type | Notes |
|---|---|---|
| type | "task" | |
| taskId | uuid | |
| sessionId | uuid | |
| channelId | string | Routes stream events back to the correct browser tab |
| attemptId | uuid? | Active task attempt created by the relay. |
| attemptNumber | int? | Monotonic attempt number for the task. |
| leaseExpiresAt | string? | ISO timestamp for the active relay lease. |
| deadlineProfile | TaskDeadlines? | Daemon supervision deadline profile in milliseconds. |
| turnKind | string? | `user`, `session_title`, `context_stats`, `compact`, or `control`. |
| prompt | string | |
| engine | "pi"? | Optional task execution engine. Empty means `"pi"`. |
| provider | string? | Optional Pi provider id. Empty means `"claude-cli"`. |
| model | string | Provider-specific model id, e.g. `claude-opus-4-6`, `gpt-5.5`, or `gpt-5.4` |
| effort | "low" \| "medium" \| "high" \| "max" | |
| permissionMode | string | e.g. `acceptEdits` |
| cwd | string | Absolute path on the daemon's machine |
| claudeSessionId | string? | Pass to `claude -p --resume` to continue an existing Claude conversation. Empty for the first turn. |
| requestId | uuid? | Optional root correlation ID for request-scoped logging. |
| traceparent | string? | W3C trace context. |
| imageUrls | string[]? | User-attached image URLs. |
| contextRefs | ContextRef[]? | Project-relative file and folder references selected in the cloud composer. |
| customInstructions | string? | Account-level instructions snapshotted onto this task. Updated daemons append this text to the Pi system prompt. |
| disableSkills | boolean? | `true` disables Claude skill discovery and explicit skill file loading for the task. |

`ContextRef`:

| Field | Type | Notes |
|---|---|---|
| kind | "file" \| "folder" | |
| path | string | Project-relative path. |
| name | string | Display name for the referenced path. |
| size | int? | File size in bytes when known. |
| modifiedAt | string? | ISO timestamp when known. |

`TaskDeadlines`:

| Field | Type | Notes |
|---|---|---|
| processStartMs | int? | Process launch deadline. |
| promptWriteMs | int? | Prompt write deadline. |
| firstEventMs | int? | Deadline for the first parsed runtime event. |
| firstVisibleEventMs | int? | Deadline for the first user-visible runtime event. |
| streamIdleMs | int? | Stream inactivity deadline. |
| toolIdleMs | int? | Tool execution inactivity deadline. |
| userInputMs | int? | User input wait deadline. |
| cleanupTermMs | int? | Grace period for process cleanup. |

`Task.contextRefs` carries project-relative file and folder references selected in the cloud composer. The relay forwards this field only to daemons that advertise `Hello.capabilities.contextRefs`.

```json
{
  "type": "task",
  "taskId": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
  "sessionId": "2c963f66-5717-4562-b3fc-3fa85f64afa6",
  "channelId": "ch_123",
  "prompt": "Inspect the project",
  "engine": "pi",
  "provider": "codex-appserver",
  "model": "gpt-5.5",
  "contextRefs": [
    { "kind": "file", "path": "apps/web/src/app/page.tsx", "name": "page.tsx" },
    { "kind": "folder", "path": "apps/web/src/components", "name": "components" }
  ]
}
```

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

### Context Compaction

#### `compactRequest`

Browser-to-daemon control message. The daemon executes Pi RPC `compact` against the session's Pi session file.

```json
{
  "type": "compactRequest",
  "sessionId": "session_123",
  "channelId": "channel_123",
  "requestId": "compact_123",
  "instructions": "preserve auth state and exact file paths"
}
```

`instructions` is optional. Empty instructions produce Pi's default compaction behavior.

#### `contextStatsRequest`

Browser-to-daemon control message. The daemon executes Pi RPC `get_session_stats` against the session's Pi session file.

```json
{
  "type": "contextStatsRequest",
  "sessionId": "session_123",
  "channelId": "channel_123",
  "requestId": "stats_123"
}
```

### `browseDir`

| Field | Type |
|---|---|
| type | "browseDir" |
| requestId | uuid |
| channelId | string |
| machineId | uuid |
| path | string |
| limit | int? |
| cursor | string? |

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

### `listSkills`

List locally available daemon skills for a project working directory. The
daemon returns bounded metadata from known skill roots; it does not return skill
file bodies.

| Field | Type |
|---|---|
| type | "listSkills" |
| requestId | uuid |
| channelId | string |
| machineId | uuid |
| cwd | string | Absolute project working directory used for local `.claude/skills` ancestry lookup. |

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

### Task attempt lifecycle

The relay owns task attempts. A dispatched `task` includes the active
`attemptId`, `attemptNumber`, `leaseExpiresAt`, `deadlineProfile`, and
`turnKind`. Daemons echo attempt metadata on task-adjacent frames so the relay
can associate runtime events with the active attempt.

`taskLifecycle` is the structured lifecycle frame for attempt diagnostics:

| Field | Type | Notes |
|---|---|---|
| type | "taskLifecycle" | |
| taskId | uuid | |
| attemptId | uuid | |
| attemptNumber | int | |
| sessionId | uuid | |
| channelId | string | |
| phase | string | Lifecycle phase. |
| status | string | Durable attempt status. |
| retryable | boolean? | Terminal retry-safety hint. |
| failureCode | string? | Stable terminal failure code. |
| message | string? | Operator-facing detail. |
| userMessage | string? | User-facing failure detail. |
| observedAt | string | RFC3339 timestamp. |
| deadlineAt | string? | Deadline associated with the phase. |
| pid | int? | Local process id when known. |
| provider | string? | Pi provider id. |
| model | string? | Provider model id. |
| requestId | uuid? | Root correlation id. |
| traceparent | string? | W3C trace context. |

Lifecycle phases are `accepted`, `queued`, `started`, `pi_started`,
`prompt_written`, `first_event_seen`, `first_visible_event_seen`, `streaming`,
`tool_started`, `tool_finished`, `waiting_input`, `input_received`,
`cleanup_started`, `cleanup_finished`, `heartbeat`, `retry_scheduled`,
`completed`, `failed`, `canceled`, `timed_out`, and `lost`.

Attempt statuses are `created`, `queued`, `started`, `pi_started`,
`prompt_written`, `first_event_seen`, `first_visible_event_seen`, `streaming`,
`waiting_input`, `tool_running`, `cleanup_started`, `cleanup_finished`,
`completed`, `failed`, `canceled`, `timed_out`, and `lost`.

The relay filters stale lifecycle frames by the active tuple
`(taskId, attemptId, sessionId, machineId)`. Frames that do not match the active
attempt are ignored for aggregate task state and may still be logged for
diagnostics. Terminal lifecycle phases map to task aggregate states:
`completed`, `failed`, `canceled`, `timed_out`, and `lost`.

Retry safety is attempt-local. A timeout before visible output or side effects
can be marked `retryable`; phases after visible output or tool execution make
automatic retry unsafe unless a higher-level policy explicitly allows it.

Control turns use `turnKind` to distinguish user-visible work from local
maintenance such as session title generation, context stats, compaction, and
daemon control flows. Consumers that do not understand attempt fields ignore
them as additive JSON fields.

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

### `contextStats`

Daemon-to-browser status message. Values come from Pi.

```json
{
  "type": "contextStats",
  "sessionId": "session_123",
  "channelId": "channel_123",
  "requestId": "stats_123",
  "tokens": 270000,
  "contextWindow": 1000000,
  "percent": 27,
  "reserveTokens": 16384,
  "keepRecentTokens": 20000,
  "autoThresholdPercent": 98.3616,
  "source": "pi",
  "observedAt": "2026-04-27T12:00:00Z"
}
```

`tokens` and `percent` may be `null` immediately after compaction.

### `compactStatus`

Daemon-to-browser lifecycle message for manual and automatic compaction.

```json
{
  "type": "compactStatus",
  "sessionId": "session_123",
  "channelId": "channel_123",
  "requestId": "compact_123",
  "status": "completed",
  "reason": "manual",
  "instructions": "preserve auth state and exact file paths",
  "tokensBefore": 8951,
  "tokensAfter": 7712,
  "contextWindow": 1000000,
  "reserveTokens": 16384,
  "keepRecentTokens": 20000,
  "autoThresholdPercent": 98.3616,
  "summary": "The session is working on Pi context compaction.",
  "firstKeptEntryId": "entry_42",
  "source": "pi",
  "observedAt": "2026-04-27T12:01:00Z"
}
```

`status` is one of `started`, `completed`, or `failed`.
`reason` is one of `manual`, `threshold`, or `overflow`.

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
| hasMore | boolean? |
| nextCursor | string? |
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

### `listSkillsResult`

| Field | Type |
|---|---|
| type | "listSkillsResult" |
| requestId | uuid |
| channelId | string |
| ok | boolean |
| skills | Skill[]? |
| error | string? |

`Skill`:

| Field | Type | Notes |
|---|---|---|
| name | string | Skill command name. |
| description | string? | Short description from `SKILL.md` frontmatter. |
| path | string | Absolute path to the skill `SKILL.md`. |
| scope | string | Discovery scope, e.g. `"home"` or `"project"`. |

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
| capabilities | HelloCapabilities? | Optional daemon feature support. |

`HelloCapabilities`:

| Field | Type | Notes |
|---|---|---|
| stop | boolean? | Daemon accepts stop messages for active task cancellation. |
| terminal | boolean? | Daemon accepts terminal lifecycle and PTY control messages. |
| agentTerminalJobs | boolean? | Daemon accepts daemon-owned agent terminal job lifecycle and control messages. |
| contextRefs | boolean? | Daemon resolves task context references before task execution. |
| previewTunnel | boolean? | Daemon accepts remote localhost preview messages. |
| previewMaxFrameBytes | int? | Maximum encoded preview frame size. |
| previewChunkBytes | int? | Raw preview body chunk target. |
| previewWebSocketProtocols | boolean? | Daemon forwards requested WebSocket subprotocols. |
| localServerDetection | boolean? | Daemon reports verified loopback web servers started by task tools. |
| skills | boolean? | Daemon accepts `listSkills` and can pass explicit Claude skill files into Pi. |
### `welcome` (relay → daemon, response to hello)

| Field | Type |
|---|---|
| type | "welcome" |
| latestDaemonVersion | string? Optional latest daemon version for update prompts |

## Remote Localhost Preview

Daemons advertise preview support in `hello.capabilities`.

Preview traffic is owner-approved, loopback-only, and routed as explicit protocol messages. Preview bytes are transient transport data and are not chat messages.

### Capability

`previewTunnel`, `previewMaxFrameBytes`, `previewChunkBytes`, and `previewWebSocketProtocols` describe daemon support.

### Lifecycle

`previewOpen` registers one loopback target for a preview. `previewClose` revokes it. `previewOpenResult.ok=false` includes `errorCode` and `message`.

### HTTP

`previewHttpRequest` carries method, origin-form path, and request headers. Request and response bodies use `previewStreamChunk` frames keyed by `streamId`.

### WebSocket

`previewWebSocketOpen` opens a target WebSocket with requested subprotocols. `previewWebSocketData` carries ordered text or binary payload bytes. `previewWebSocketClose` closes both sides.

### Stream Cancellation

`previewStreamCancel` cancels local IO for the stream. Receivers treat duplicate, missing, or out-of-order chunks as stream errors.

### Local Server Detection

`localServerDetected` is emitted by the daemon when task tool output identifies a reachable loopback HTTP server. The daemon verifies the port before emitting the event. The relay forwards the event to the session channel; browsers can use it to start an owner-scoped preview for the reported port.

| Field | Type | Notes |
|---|---|---|
| type | "localServerDetected" | |
| sessionId | string | Chat session that started the server. |
| channelId | string | Browser channel for the session. |
| taskId | string? | Task that produced the server output. |
| toolUseId | string? | Tool call that produced the server output. |
| host | string | Normalized loopback host, usually `127.0.0.1`. |
| port | int | Verified target port. |
| url | string | Loopback URL for the server. |
| command | string? | Shell command associated with the tool call. |
| source | string | Detection source, currently `tool_output`. |
| detectedAt | string | RFC3339 timestamp. |

## Terminal Messages

Terminal messages open and control a chat-scoped PTY on the paired daemon machine. Browser-originated `terminalOpen` carries `token`, while relay-to-daemon `terminalOpen` carries server-derived `terminalId`, `cwd`, `idleTimeoutMs`, and `maxLifetimeMs`. Terminal input and output bytes are base64-encoded live transport data.

### Capability

Daemons advertise terminal support through `hello.capabilities.terminal`.

### `terminalOpen`

| Field | Type | Notes |
|---|---|---|
| type | "terminalOpen" | |
| requestId | string | Correlates the open request with opened/error responses. |
| terminalId | string? | Relay-generated terminal id for daemon-bound opens. |
| sessionId | uuid | Chat session scope. |
| channelId | string | Owning browser channel. |
| token | string? | Browser terminal-open capability. |
| cwd | string? | Server-derived daemon working directory. |
| cols | int | Requested terminal columns. |
| rows | int | Requested terminal rows. |
| idleTimeoutMs | int? | Daemon idle/disconnect timeout. |
| maxLifetimeMs | int? | Daemon maximum terminal lifetime. |

### `terminalOpened`

| Field | Type |
|---|---|
| type | "terminalOpened" |
| requestId | string |
| terminalId | string |
| sessionId | uuid |
| channelId | string |
| shell | string |
| cwd | string |
| startedAt | iso8601 string |

### `terminalInput`

| Field | Type |
|---|---|
| type | "terminalInput" |
| terminalId | string |
| channelId | string |
| dataBase64 | string |

### `terminalOutput`

| Field | Type |
|---|---|
| type | "terminalOutput" |
| terminalId | string |
| sessionId | uuid |
| channelId | string |
| seq | int64 |
| dataBase64 | string |

### `terminalSnapshot`

| Field | Type |
|---|---|
| type | "terminalSnapshot" |
| terminalId | string |
| sessionId | uuid |
| channelId | string |
| seq | int64 |
| dataBase64 | string |

### `terminalResize`

| Field | Type |
|---|---|
| type | "terminalResize" |
| terminalId | string |
| channelId | string |
| cols | int |
| rows | int |

### `terminalClose`

| Field | Type |
|---|---|
| type | "terminalClose" |
| terminalId | string |
| channelId | string |

### `terminalExit`

| Field | Type |
|---|---|
| type | "terminalExit" |
| terminalId | string |
| sessionId | uuid |
| channelId | string |
| exitCode | int? |
| signal | string? |
| reason | string |
| endedAt | iso8601 string |

### `terminalError`

| Field | Type |
|---|---|
| type | "terminalError" |
| requestId | string? |
| terminalId | string? |
| sessionId | uuid? |
| channelId | string |
| error | string |

## Agent Terminal Jobs

Agent terminal jobs are daemon-owned PTY processes started by agent tools and surfaced to browsers as attachable terminal streams. Daemons advertise support through `hello.capabilities.agentTerminalJobs`.

The daemon creates an agent terminal route by sending `agentTerminalStarted`. The relay validates the session, paired machine, owning user, and channel before registering the route and forwarding the event to the browser. The daemon sends `agentTerminalUpdated` whenever job metadata advances.

Browser attach uses `agentTerminalAttach`; browser snapshot refresh uses `agentTerminalSnapshotRequest`. Agent terminal log bytes use the existing terminal data plane: `terminalOutput`, `terminalSnapshot`, `terminalExit`, and `terminalError`.

### `agentTerminalStarted`

| Field | Type | Notes |
|---|---|---|
| type | "agentTerminalStarted" | |
| jobId | string | Daemon job id. |
| terminalId | string | Terminal stream id for the job PTY. |
| sessionId | uuid | Chat session scope. |
| channelId | string | Browser channel for the session. |
| taskId | string? | Agent task that started the job. |
| toolCallId | string? | Tool call that started the job. |
| projectId | uuid | Project scope. |
| commandPreview | string | Redacted command summary for UI metadata. |
| title | string | Human-readable terminal title. |
| cwd | string | Normalized daemon working directory. |
| status | string | `starting`, `running`, `ready`, `exited`, `failed`, or `killed`. |
| readiness | AgentTerminalReadiness | Readiness state. |
| ports | AgentTerminalPort[]? | Detected loopback ports. |
| urls | string[]? | Detected loopback URLs. |
| seq | int64? | Current output sequence. |
| startedAt | iso8601 string | Job start timestamp. |

### `agentTerminalUpdated`

| Field | Type | Notes |
|---|---|---|
| type | "agentTerminalUpdated" | |
| jobId | string | Daemon job id. |
| terminalId | string | Terminal stream id for the job PTY. |
| sessionId | uuid | Chat session scope. |
| channelId | string | Browser channel for the session. |
| status | string | Current lifecycle status. |
| readiness | AgentTerminalReadiness | Current readiness state. |
| ports | AgentTerminalPort[]? | Detected loopback ports. |
| urls | string[]? | Detected loopback URLs. |
| seq | int64? | Current output sequence. |
| updatedAt | iso8601 string | Update timestamp. |

### `AgentTerminalReadiness`

| Field | Type | Notes |
|---|---|---|
| state | string | `unknown`, `waiting`, `ready`, `timed_out`, or `failed`. |
| source | string? | `pattern`, `port`, `url`, `process_exit`, or `heuristic`. |
| matchedText | string? | Output text that satisfied readiness. |
| readyAt | iso8601 string? | Readiness timestamp. |
| timeoutMs | int? | Readiness wait timeout. |

### `AgentTerminalPort`

| Field | Type |
|---|---|
| host | string |
| port | int |
| url | string |

### `agentTerminalAttach`

| Field | Type |
|---|---|
| type | "agentTerminalAttach" |
| terminalId | string |
| channelId | string |

### `agentTerminalSnapshotRequest`

| Field | Type |
|---|---|
| type | "agentTerminalSnapshotRequest" |
| terminalId | string |
| channelId | string |
