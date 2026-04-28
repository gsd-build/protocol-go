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

`engine` selects the daemon task executor. Daemons that do not support this field ignore it through normal JSON decoding.

| Field | Type | Notes |
|---|---|---|
| type | "task" | |
| taskId | uuid | |
| sessionId | uuid | |
| channelId | string | Routes stream events back to the correct browser tab |
| prompt | string | |
| engine | "claude" \| "pi"? | Optional task execution engine. Empty means `"claude"`. |
| model | string | e.g. `claude-opus-4-6[1m]` |
| effort | "low" \| "medium" \| "high" \| "max" | |
| permissionMode | string | e.g. `acceptEdits` |
| cwd | string | Absolute path on the daemon's machine |
| claudeSessionId | string? | Pass to `claude -p --resume` to continue an existing Claude conversation. Empty for the first turn. |
| requestId | uuid? | Optional root correlation ID for request-scoped logging. |
| traceparent | string? | W3C trace context. |
| imageUrls | string[]? | User-attached image URLs. |

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
| previewTunnel | boolean? | Daemon accepts remote localhost preview messages. |
| previewMaxFrameBytes | int? | Maximum encoded preview frame size. |
| previewChunkBytes | int? | Raw preview body chunk target. |
| previewWebSocketProtocols | boolean? | Daemon forwards requested WebSocket subprotocols. |

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
