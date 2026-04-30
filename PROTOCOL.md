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
| planCapability | PlanCapability? | Task-scoped bearer capability for Plan Mode tools. The daemon passes it to the Pi process as environment variables and never places it in the prompt. |

`ContextRef`:

| Field | Type | Notes |
|---|---|---|
| kind | "file" \| "folder" | |
| path | string | Project-relative path. |
| name | string | Display name for the referenced path. |
| size | int? | File size in bytes when known. |
| modifiedAt | string? | ISO timestamp when known. |

`PlanCapability`:

| Field | Type | Notes |
|---|---|---|
| token | string | Bearer token scoped to one task. |
| id | string? | Capability record identifier. |
| attemptId | uuid? | Attempt that owns this capability. |
| apiBaseUrl | string | Cloud app base URL for `/api/agent-plan/*`. |
| expiresAt | string | ISO timestamp. |
| snapshot | object? | Capability metadata snapshot used by cloud-side authorization. |

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

### `planningEvent`

Append-only planning journal event sent from a source runtime to the relay.

| Field | Type | Notes |
|---|---|---|
| type | "planningEvent" | |
| eventId | string | Unique event identifier. |
| schemaVersion | int | Journal schema version. |
| projectionVersion | int | Projection contract version. |
| projectId | uuid/string | Project that owns the event. |
| sourceId | string | Stable source identity, such as daemon or SDK instance. |
| sourceKind | "daemon" \| "sdk" \| "cloud" \| "import" | Producer class. |
| sourceSeq | int64 | Monotonic sequence for `sourceId`; consumers use this for ordering and gap detection. |
| sourceCursor | string? | Optional source-specific resume cursor. |
| runId | string | Runtime run identity. |
| workstreamId | string? | Optional workstream identity. |
| planId | uuid/string? | |
| itemId | uuid/string? | |
| actorType | string | Agent, human, runtime, verifier, or system. |
| actorId | string | Stable actor identity. |
| actorRole | string? | |
| sessionId | uuid? | |
| taskId | uuid? | |
| eventKind | string | Event-specific verb such as `plan.done`. |
| idempotencyKey | string | Stable producer-provided key. Replays with the same request body are accepted idempotently. |
| causationId | string? | Event ID or command ID that caused this event. |
| occurredAt | iso8601 string | Producer-side timestamp. |
| payload | object | Event-specific structured data. It is data, not executable instructions. |
| evidenceIds | string[]? | Evidence records attached to this event. |
| parentEventIds | string[]? | Causal parent events. |
| trace | object? | Optional diagnostic trace context. |

### `planningEventAck`

Relay acknowledgement for a `planningEvent`.

| Field | Type | Notes |
|---|---|---|
| type | "planningEventAck" | |
| eventId | string | Event being acknowledged. |
| sourceId | string | Source identity from the event. |
| sourceSeq | int64 | Source sequence from the event. |
| accepted | boolean | `true` when persisted or idempotently replayed. |
| error | string? | Human-readable rejection reason. |

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

## Shared Browser Sessions

Shared browser sessions let a browser client and daemon coordinate a task-scoped local browser controlled through `gsd-browser`.

### Capability

Browser support is advertised in `hello.capabilities`.

| Capability | Type | Meaning |
|---|---|---|
| browserSessions | bool | Daemon can open task-scoped shared browser sessions. |
| browserFrameStream | bool | Daemon can stream screenshot frames with `browserFrame`. |
| browserUserControl | bool | Daemon can arbitrate Lex and agent browser control. |
| browserIdentities | bool | Daemon can isolate browser state by identity scope and key. |
| browserSensitiveActionApproval | bool | Daemon can pause browser tool execution for approval. |
| browserMaxFrameBytes | int | Maximum encoded browser frame size the daemon can send. |

### Lifecycle

`browserSessionOpen` starts or attaches to a local browser session for one grant. `browserSessionClose` revokes that session. `browserSessionError` reports stable error codes.

#### `browserSessionOpen`

| Field | Type |
|---|---|
| type | "browserSessionOpen" |
| requestId | string |
| grantId | string |
| sessionId | uuid |
| projectId | uuid |
| taskId | uuid |
| channelId | string |
| machineId | string |
| identityId | uuid? |
| mode | string |
| expiresAt | iso8601 string |

#### `browserSessionOpened`

| Field | Type |
|---|---|
| type | "browserSessionOpened" |
| requestId | string |
| browserId | string |
| grantId | string |
| sessionId | uuid |
| channelId | string |
| url | string? |
| title | string? |
| openedAt | iso8601 string |

#### `browserSessionClose`

| Field | Type |
|---|---|
| type | "browserSessionClose" |
| browserId | string |
| grantId | string |
| sessionId | uuid |
| channelId | string |
| reason | string? |

#### `browserSessionClosed`

| Field | Type |
|---|---|
| type | "browserSessionClosed" |
| browserId | string |
| sessionId | uuid |
| channelId | string |
| reason | string? |
| closedAt | iso8601 string |

#### `browserSessionError`

| Field | Type |
|---|---|
| type | "browserSessionError" |
| browserId | string? |
| requestId | string? |
| sessionId | uuid |
| channelId | string |
| code | string |
| message | string |

### Visual State

`browserFrame` carries bounded screenshot frames. `browserCursor`, `browserNavigation`, and `browserAction` carry compact visible state and action metadata.

#### `browserFrame`

| Field | Type |
|---|---|
| type | "browserFrame" |
| browserId | string |
| sessionId | uuid |
| channelId | string |
| seq | int64 |
| contentType | string |
| dataBase64 | base64 string |
| frameRef | string? |
| width | int |
| height | int |
| viewportWidth | int? |
| viewportHeight | int? |
| devicePixelRatio | float64? |
| capturedAt | iso8601 string |
| droppedPriorCount | int? |

#### `browserCursor`

| Field | Type |
|---|---|
| type | "browserCursor" |
| browserId | string |
| sessionId | uuid |
| channelId | string |
| owner | string |
| x | float64 |
| y | float64 |
| updatedAt | iso8601 string |

#### `browserNavigation`

| Field | Type |
|---|---|
| type | "browserNavigation" |
| browserId | string |
| sessionId | uuid |
| channelId | string |
| url | string |
| title | string? |
| startedAt | iso8601 string? |
| endedAt | iso8601 string? |

#### `browserAction`

| Field | Type |
|---|---|
| type | "browserAction" |
| browserId | string |
| sessionId | uuid |
| channelId | string |
| taskId | uuid? |
| toolUseId | string? |
| summary | string |
| status | string |
| metadata | json? |
| at | iso8601 string |

### Control

`browserControlClaim`, `browserControlRelease`, and `browserUserInput` coordinate Lex control. Tool calls execute while control owner is `agent`.

#### `browserControlClaim`

| Field | Type |
|---|---|
| type | "browserControlClaim" |
| browserId | string |
| sessionId | uuid |
| channelId | string |
| owner | string |
| reason | string? |
| controlVersion | int64? |

#### `browserControlRelease`

| Field | Type |
|---|---|
| type | "browserControlRelease" |
| browserId | string |
| sessionId | uuid |
| channelId | string |
| owner | string |
| reason | string? |
| controlVersion | int64? |

#### `browserUserInput`

| Field | Type |
|---|---|
| type | "browserUserInput" |
| inputId | string? |
| browserId | string |
| sessionId | uuid |
| channelId | string |
| owner | string |
| kind | string |
| x | float64? |
| y | float64? |
| text | string? |
| key | string? |
| deltaX | float64? |
| deltaY | float64? |
| frameSeq | int64? |
| controlVersion | int64? |
| coordinateSpace | string? |
| viewportWidth | int? |
| viewportHeight | int? |
| frameWidth | int? |
| frameHeight | int? |
| devicePixelRatio | float64? |
| renderedLeft | float64 |
| renderedTop | float64 |
| renderedWidth | float64 |
| renderedHeight | float64 |

#### `browserUserInputAck`

| Field | Type |
|---|---|
| type | "browserUserInputAck" |
| browserId | string |
| sessionId | uuid |
| channelId | string |
| inputId | string? |
| accepted | bool |
| reason | string? |
| controlVersion | int64? |
| ackedAt | iso8601 string |

#### `browserTransportStatus`

| Field | Type |
|---|---|
| type | "browserTransportStatus" |
| browserId | string |
| sessionId | uuid |
| channelId | string |
| status | string |
| queueDepth | int? |
| droppedFrameCount | int64? |
| maxFrameBytes | int? |
| at | iso8601 string |

### Browser Input Safety

Browser input messages include the source frame sequence, control version,
viewport dimensions, captured frame dimensions, rendered panel rectangle, and
coordinate space. Receivers reject stale-frame and stale-control input with
`browserUserInputAck`.

### Browser Transport Status

Browser frame producers may drop, coalesce, downscale, or reference frames when
encoded frames exceed transport limits. Control and approval messages have
priority over frame delivery.

### Preview Bridge Access

Preview-origin bridge access is short-lived, grant-bound, and separate from
iframe preview cookies. Local-direct bridge mode navigates the browser directly
to the approved loopback target.

#### `browserBridgeAccessOpen`

| Field | Type |
|---|---|
| type | "browserBridgeAccessOpen" |
| requestId | string |
| previewId | string |
| grantId | string |
| browserId | string? |
| sessionId | uuid |
| channelId | string |
| machineId | string |
| bridgeMode | string |
| requestedAt | iso8601 string |

#### `browserBridgeAccessOpened`

| Field | Type |
|---|---|
| type | "browserBridgeAccessOpened" |
| requestId | string |
| previewId | string |
| grantId | string |
| browserId | string? |
| sessionId | uuid |
| channelId | string |
| bridgeMode | string |
| url | string |
| expiresAt | iso8601 string |

#### `browserBridgeAccessClose`

| Field | Type |
|---|---|
| type | "browserBridgeAccessClose" |
| previewId | string |
| grantId | string |
| browserId | string? |
| sessionId | uuid |
| channelId | string |
| reason | string? |

### Tool Calls

`browserToolCall` and `browserToolResult` represent agent browser tool execution. The daemon validates active grant state before executing each call.

#### `browserToolCall`

| Field | Type |
|---|---|
| type | "browserToolCall" |
| browserId | string |
| grantId | string |
| taskId | uuid |
| toolUseId | string |
| method | string |
| paramsJson | json? |

#### `browserToolResult`

| Field | Type |
|---|---|
| type | "browserToolResult" |
| browserId | string |
| grantId | string |
| taskId | uuid |
| toolUseId | string |
| ok | bool |
| resultJson | json? |
| error | string? |

### Sensitive Actions

`browserSensitiveActionRequest` pauses execution and asks Lex for approval. `browserSensitiveActionResponse` resumes or denies the action.

#### `browserSensitiveActionRequest`

| Field | Type |
|---|---|
| type | "browserSensitiveActionRequest" |
| browserId | string |
| requestId | string |
| sessionId | uuid |
| channelId | string |
| taskId | uuid |
| toolUseId | string |
| category | string |
| summary | string |
| origin | string? |
| expiresAt | iso8601 string |

#### `browserSensitiveActionResponse`

| Field | Type |
|---|---|
| type | "browserSensitiveActionResponse" |
| browserId | string |
| requestId | string |
| sessionId | uuid |
| channelId | string |
| approved | bool |
| reason | string? |
