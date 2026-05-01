// Package protocol defines the wire format between the GSD Cloud daemon,
// the Fly.io relay, and the browser. See PROTOCOL.md for the authoritative
// specification; every change here must be mirrored in that file.
package protocol

import (
	"encoding/json"
	"time"
)

type MessageType = string

// Message type constants.
const (
	MsgTypeTask                            = "task"
	MsgTypeTaskLifecycle                   = "taskLifecycle"
	MsgTypeStop                            = "stop"
	MsgTypePermissionResponse              = "permissionResponse"
	MsgTypeQuestionResponse                = "questionResponse"
	MsgTypeBrowseDir                       = "browseDir"
	MsgTypeReadFile                        = "readFile"
	MsgTypeMkDir                           = "mkDir"
	MsgTypeMkDirResult                     = "mkDirResult"
	MsgTypeListSkills                      = "listSkills"
	MsgTypeListSkillsResult                = "listSkillsResult"
	MsgTypeCompactRequest      MessageType = "compactRequest"
	MsgTypeContextStatsRequest MessageType = "contextStatsRequest"

	MsgTypeStream                        = "stream"
	MsgTypeTaskStarted                   = "taskStarted"
	MsgTypeTaskComplete                  = "taskComplete"
	MsgTypeTaskError                     = "taskError"
	MsgTypeTaskCancelled                 = "taskCancelled"
	MsgTypePermissionRequest             = "permissionRequest"
	MsgTypeQuestion                      = "question"
	MsgTypeHeartbeat                     = "heartbeat"
	MsgTypeBrowseDirResult               = "browseDirResult"
	MsgTypeReadFileResult                = "readFileResult"
	MsgTypeContextStats      MessageType = "contextStats"
	MsgTypeCompactStatus     MessageType = "compactStatus"

	MsgTypeHello   = "hello"
	MsgTypeWelcome = "welcome"

	MsgTypeMachineStatus              = "machineStatus"
	MsgTypePreviewOpen                = "previewOpen"
	MsgTypePreviewOpenResult          = "previewOpenResult"
	MsgTypePreviewClose               = "previewClose"
	MsgTypePreviewHTTPRequest         = "previewHttpRequest"
	MsgTypePreviewHTTPResponseHead    = "previewHttpResponseHead"
	MsgTypePreviewStreamChunk         = "previewStreamChunk"
	MsgTypePreviewStreamCancel        = "previewStreamCancel"
	MsgTypePreviewWebSocketOpen       = "previewWebSocketOpen"
	MsgTypePreviewWebSocketOpenResult = "previewWebSocketOpenResult"
	MsgTypePreviewWebSocketData       = "previewWebSocketData"
	MsgTypePreviewWebSocketClose      = "previewWebSocketClose"
	MsgTypeLocalServerDetected        = "localServerDetected"

	MsgTypeTerminalOpen     = "terminalOpen"
	MsgTypeTerminalOpened   = "terminalOpened"
	MsgTypeTerminalInput    = "terminalInput"
	MsgTypeTerminalOutput   = "terminalOutput"
	MsgTypeTerminalSnapshot = "terminalSnapshot"
	MsgTypeTerminalResize   = "terminalResize"
	MsgTypeTerminalClose    = "terminalClose"
	MsgTypeTerminalExit     = "terminalExit"
	MsgTypeTerminalError    = "terminalError"

	MsgTypeAgentTerminalStarted         = "agentTerminalStarted"
	MsgTypeAgentTerminalUpdated         = "agentTerminalUpdated"
	MsgTypeAgentTerminalAttach          = "agentTerminalAttach"
	MsgTypeAgentTerminalSnapshotRequest = "agentTerminalSnapshotRequest"

	MsgTypeBrowserSessionOpen             MessageType = "browserSessionOpen"
	MsgTypeBrowserSessionOpened           MessageType = "browserSessionOpened"
	MsgTypeBrowserSessionClose            MessageType = "browserSessionClose"
	MsgTypeBrowserSessionClosed           MessageType = "browserSessionClosed"
	MsgTypeBrowserSessionError            MessageType = "browserSessionError"
	MsgTypeBrowserFrame                   MessageType = "browserFrame"
	MsgTypeBrowserRefs                    MessageType = "browserRefs"
	MsgTypeBrowserCursor                  MessageType = "browserCursor"
	MsgTypeBrowserNavigation              MessageType = "browserNavigation"
	MsgTypeBrowserAction                  MessageType = "browserAction"
	MsgTypeBrowserToolCall                MessageType = "browserToolCall"
	MsgTypeBrowserToolResult              MessageType = "browserToolResult"
	MsgTypeBrowserToolCallStarted         MessageType = "browserToolCallStarted"
	MsgTypeBrowserToolCallUpdated         MessageType = "browserToolCallUpdated"
	MsgTypeBrowserArtifactCreated         MessageType = "browserArtifactCreated"
	MsgTypeBrowserControlClaim            MessageType = "browserControlClaim"
	MsgTypeBrowserControlRelease          MessageType = "browserControlRelease"
	MsgTypeBrowserUserInput               MessageType = "browserUserInput"
	MsgTypeBrowserUserInputAck            MessageType = "browserUserInputAck"
	MsgTypeBrowserTransportStatus         MessageType = "browserTransportStatus"
	MsgTypeBrowserBridgeAccessOpen        MessageType = "browserBridgeAccessOpen"
	MsgTypeBrowserBridgeAccessOpened      MessageType = "browserBridgeAccessOpened"
	MsgTypeBrowserBridgeAccessClose       MessageType = "browserBridgeAccessClose"
	MsgTypeBrowserSensitiveActionRequest  MessageType = "browserSensitiveActionRequest"
	MsgTypeBrowserSensitiveActionResponse MessageType = "browserSensitiveActionResponse"

	MsgTypePlanningEvent    MessageType = "planningEvent"
	MsgTypePlanningEventAck MessageType = "planningEventAck"
)

const (
	BrowserOwnerAgent    = "agent"
	BrowserOwnerLex      = "lex"
	BrowserOwnerPaused   = "paused"
	BrowserOwnerApproval = "approval"

	BrowserInputKindClick         = "click"
	BrowserInputKindPointerMove   = "pointer_move"
	BrowserInputKindPointerDown   = "pointer_down"
	BrowserInputKindPointerUp     = "pointer_up"
	BrowserInputKindWheel         = "wheel"
	BrowserInputKindKeyDown       = "key_down"
	BrowserInputKindKeyUp         = "key_up"
	BrowserInputKindText          = "text"
	BrowserInputKindNavigate      = "navigate"
	BrowserInputKindBack          = "back"
	BrowserInputKindForward       = "forward"
	BrowserInputKindReload        = "reload"
	BrowserInputKindRefAction     = "ref_action"
	BrowserInputKindSetViewport   = "set_viewport"
	BrowserInputKindEmulateDevice = "emulate_device"

	BrowserCoordinateSpaceFrameCssPixels = "frame_css_pixels"

	BrowserInputRejectStaleFrame     = "stale_frame"
	BrowserInputRejectOwnerMismatch  = "owner_mismatch"
	BrowserInputRejectExpiredGrant   = "expired_grant"
	BrowserInputRejectInvalidPayload = "invalid_payload"
)

type ContextRef struct {
	Kind       string `json:"kind"`
	Path       string `json:"path"`
	Name       string `json:"name"`
	Size       *int64 `json:"size,omitempty"`
	ModifiedAt string `json:"modifiedAt,omitempty"`
}

type PlanCapability struct {
	ID         string          `json:"id,omitempty"`
	AttemptID  string          `json:"attemptId,omitempty"`
	Token      string          `json:"token"`
	APIBaseURL string          `json:"apiBaseUrl"`
	ExpiresAt  string          `json:"expiresAt"`
	Snapshot   json.RawMessage `json:"snapshot,omitempty"`
}

type BrowserGrantContext struct {
	GrantID   string `json:"grantId"`
	ProjectID string `json:"projectId"`
	SessionID string `json:"sessionId"`
	TaskID    string `json:"taskId"`
	ChannelID string `json:"channelId"`
	MachineID string `json:"machineId"`
	ExpiresAt string `json:"expiresAt"`
}

type PlanningEvent struct {
	Type              MessageType     `json:"type"`
	EventID           string          `json:"eventId"`
	SchemaVersion     int             `json:"schemaVersion"`
	ProjectionVersion int             `json:"projectionVersion"`
	ProjectID         string          `json:"projectId"`
	SourceID          string          `json:"sourceId"`
	SourceKind        string          `json:"sourceKind"`
	SourceSeq         int64           `json:"sourceSeq"`
	SourceCursor      string          `json:"sourceCursor,omitempty"`
	RunID             string          `json:"runId"`
	WorkstreamID      string          `json:"workstreamId,omitempty"`
	PlanID            string          `json:"planId,omitempty"`
	ItemID            string          `json:"itemId,omitempty"`
	ActorType         string          `json:"actorType"`
	ActorID           string          `json:"actorId"`
	ActorRole         string          `json:"actorRole,omitempty"`
	SessionID         string          `json:"sessionId,omitempty"`
	TaskID            string          `json:"taskId,omitempty"`
	EventKind         string          `json:"eventKind"`
	IdempotencyKey    string          `json:"idempotencyKey"`
	CausationID       string          `json:"causationId,omitempty"`
	OccurredAt        string          `json:"occurredAt"`
	PayloadJSON       json.RawMessage `json:"payload"`
	EvidenceIDs       []string        `json:"evidenceIds,omitempty"`
	ParentEventIDs    []string        `json:"parentEventIds,omitempty"`
	TraceJSON         json.RawMessage `json:"trace,omitempty"`
}

type PlanningEventAck struct {
	Type      MessageType `json:"type"`
	EventID   string      `json:"eventId"`
	SourceID  string      `json:"sourceId"`
	SourceSeq int64       `json:"sourceSeq"`
	Accepted  bool        `json:"accepted"`
	Error     string      `json:"error,omitempty"`
}

type TurnKind string

const (
	TurnKindUser         TurnKind = "user"
	TurnKindSessionTitle TurnKind = "session_title"
	TurnKindContextStats TurnKind = "context_stats"
	TurnKindCompact      TurnKind = "compact"
	TurnKindControl      TurnKind = "control"
)

type TaskDeadlines struct {
	ProcessStartMs      int `json:"processStartMs,omitempty"`
	PromptWriteMs       int `json:"promptWriteMs,omitempty"`
	FirstEventMs        int `json:"firstEventMs,omitempty"`
	FirstVisibleEventMs int `json:"firstVisibleEventMs,omitempty"`
	StreamIdleMs        int `json:"streamIdleMs,omitempty"`
	ToolIdleMs          int `json:"toolIdleMs,omitempty"`
	UserInputMs         int `json:"userInputMs,omitempty"`
	CleanupTermMs       int `json:"cleanupTermMs,omitempty"`
}

type TaskLifecyclePhase string

const (
	TaskLifecyclePhaseAccepted              TaskLifecyclePhase = "accepted"
	TaskLifecyclePhaseQueued                TaskLifecyclePhase = "queued"
	TaskLifecyclePhaseStarted               TaskLifecyclePhase = "started"
	TaskLifecyclePhasePiStarted             TaskLifecyclePhase = "pi_started"
	TaskLifecyclePhasePromptWritten         TaskLifecyclePhase = "prompt_written"
	TaskLifecyclePhaseFirstEventSeen        TaskLifecyclePhase = "first_event_seen"
	TaskLifecyclePhaseFirstVisibleEventSeen TaskLifecyclePhase = "first_visible_event_seen"
	TaskLifecyclePhaseStreaming             TaskLifecyclePhase = "streaming"
	TaskLifecyclePhaseToolStarted           TaskLifecyclePhase = "tool_started"
	TaskLifecyclePhaseToolFinished          TaskLifecyclePhase = "tool_finished"
	TaskLifecyclePhaseWaitingInput          TaskLifecyclePhase = "waiting_input"
	TaskLifecyclePhaseInputReceived         TaskLifecyclePhase = "input_received"
	TaskLifecyclePhaseCleanupStarted        TaskLifecyclePhase = "cleanup_started"
	TaskLifecyclePhaseCleanupFinished       TaskLifecyclePhase = "cleanup_finished"
	TaskLifecyclePhaseHeartbeat             TaskLifecyclePhase = "heartbeat"
	TaskLifecyclePhaseRetryScheduled        TaskLifecyclePhase = "retry_scheduled"
	TaskLifecyclePhaseCompleted             TaskLifecyclePhase = "completed"
	TaskLifecyclePhaseFailed                TaskLifecyclePhase = "failed"
	TaskLifecyclePhaseCanceled              TaskLifecyclePhase = "canceled"
	TaskLifecyclePhaseTimedOut              TaskLifecyclePhase = "timed_out"
	TaskLifecyclePhaseLost                  TaskLifecyclePhase = "lost"
)

type TaskAttemptStatus string

const (
	TaskAttemptStatusCreated               TaskAttemptStatus = "created"
	TaskAttemptStatusQueued                TaskAttemptStatus = "queued"
	TaskAttemptStatusStarted               TaskAttemptStatus = "started"
	TaskAttemptStatusPiStarted             TaskAttemptStatus = "pi_started"
	TaskAttemptStatusPromptWritten         TaskAttemptStatus = "prompt_written"
	TaskAttemptStatusFirstEventSeen        TaskAttemptStatus = "first_event_seen"
	TaskAttemptStatusFirstVisibleEventSeen TaskAttemptStatus = "first_visible_event_seen"
	TaskAttemptStatusStreaming             TaskAttemptStatus = "streaming"
	TaskAttemptStatusWaitingInput          TaskAttemptStatus = "waiting_input"
	TaskAttemptStatusToolRunning           TaskAttemptStatus = "tool_running"
	TaskAttemptStatusCleanupStarted        TaskAttemptStatus = "cleanup_started"
	TaskAttemptStatusCleanupFinished       TaskAttemptStatus = "cleanup_finished"
	TaskAttemptStatusCompleted             TaskAttemptStatus = "completed"
	TaskAttemptStatusFailed                TaskAttemptStatus = "failed"
	TaskAttemptStatusCanceled              TaskAttemptStatus = "canceled"
	TaskAttemptStatusTimedOut              TaskAttemptStatus = "timed_out"
	TaskAttemptStatusLost                  TaskAttemptStatus = "lost"
)

// Task is sent from the browser to the daemon to dispatch a user message.
type Task struct {
	Type               string               `json:"type"`
	TaskID             string               `json:"taskId"`
	SessionID          string               `json:"sessionId"`
	ChannelID          string               `json:"channelId"`
	AttemptID          string               `json:"attemptId,omitempty"`
	AttemptNumber      int                  `json:"attemptNumber,omitempty"`
	LeaseExpiresAt     string               `json:"leaseExpiresAt,omitempty"`
	DeadlineProfile    TaskDeadlines        `json:"deadlineProfile,omitempty"`
	TurnKind           TurnKind             `json:"turnKind,omitempty"`
	Prompt             string               `json:"prompt"`
	Engine             string               `json:"engine,omitempty"`   // "pi"; empty defaults to pi
	Provider           string               `json:"provider,omitempty"` // Pi provider; empty defaults to claude-cli
	Model              string               `json:"model"`
	Effort             string               `json:"effort"`
	PermissionMode     string               `json:"permissionMode"`
	CWD                string               `json:"cwd"`
	ClaudeSessionID    string               `json:"claudeSessionId,omitempty"` // passed to --resume
	RequestID          string               `json:"requestId,omitempty"`
	Traceparent        string               `json:"traceparent,omitempty"` // W3C trace context
	ImageURLs          []string             `json:"imageUrls,omitempty"`   // user-attached image URLs
	ContextRefs        []ContextRef         `json:"contextRefs,omitempty"`
	CustomInstructions string               `json:"customInstructions,omitempty"`
	DisableSkills      bool                 `json:"disableSkills,omitempty"`
	PlanCapability     *PlanCapability      `json:"planCapability,omitempty"`
	BrowserGrant       *BrowserGrantContext `json:"browserGrant,omitempty"`
}

type TaskLifecycle struct {
	Type          MessageType        `json:"type"`
	TaskID        string             `json:"taskId"`
	AttemptID     string             `json:"attemptId"`
	AttemptNumber int                `json:"attemptNumber"`
	SessionID     string             `json:"sessionId"`
	ChannelID     string             `json:"channelId"`
	Phase         TaskLifecyclePhase `json:"phase"`
	Status        TaskAttemptStatus  `json:"status"`
	Retryable     bool               `json:"retryable,omitempty"`
	FailureCode   string             `json:"failureCode,omitempty"`
	Message       string             `json:"message,omitempty"`
	UserMessage   string             `json:"userMessage,omitempty"`
	ObservedAt    time.Time          `json:"observedAt"`
	DeadlineAt    *time.Time         `json:"deadlineAt,omitempty"`
	PID           int                `json:"pid,omitempty"`
	Provider      string             `json:"provider,omitempty"`
	Model         string             `json:"model,omitempty"`
	RequestID     string             `json:"requestId,omitempty"`
	Traceparent   string             `json:"traceparent,omitempty"`
}

// Stop asks the daemon to interrupt the current Claude process for a session.
type Stop struct {
	Type      string `json:"type"`
	ChannelID string `json:"channelId"`
	SessionID string `json:"sessionId"`
}

// PermissionResponse is the browser's answer to a permission request.
type PermissionResponse struct {
	Type      string `json:"type"`
	ChannelID string `json:"channelId"`
	SessionID string `json:"sessionId"`
	RequestID string `json:"requestId"`
	Approved  bool   `json:"approved"`
}

// QuestionResponse is the browser's answer to a question.
type QuestionResponse struct {
	Type      string `json:"type"`
	ChannelID string `json:"channelId"`
	SessionID string `json:"sessionId"`
	RequestID string `json:"requestId"`
	Answer    string `json:"answer"`
}

// BrowseDir lists directory contents on the daemon's machine.
type BrowseDir struct {
	Type      string `json:"type"`
	RequestID string `json:"requestId"`
	ChannelID string `json:"channelId"`
	MachineID string `json:"machineId"`
	Path      string `json:"path"`
	Limit     int    `json:"limit,omitempty"`
	Cursor    string `json:"cursor,omitempty"`
}

// ReadFile reads a file from the daemon's filesystem.
type ReadFile struct {
	Type      string `json:"type"`
	RequestID string `json:"requestId"`
	ChannelID string `json:"channelId"`
	MachineID string `json:"machineId"`
	Path      string `json:"path"`
	MaxBytes  int    `json:"maxBytes,omitempty"`
}

type CompactReason string

const (
	CompactReasonManual    CompactReason = "manual"
	CompactReasonThreshold CompactReason = "threshold"
	CompactReasonOverflow  CompactReason = "overflow"
)

type CompactLifecycleStatus string

const (
	CompactStatusStarted   CompactLifecycleStatus = "started"
	CompactStatusCompleted CompactLifecycleStatus = "completed"
	CompactStatusFailed    CompactLifecycleStatus = "failed"
)

type CompactRequest struct {
	Type         MessageType `json:"type"`
	SessionID    string      `json:"sessionId"`
	ChannelID    string      `json:"channelId"`
	RequestID    string      `json:"requestId"`
	Instructions string      `json:"instructions,omitempty"`
}

type ContextStatsRequest struct {
	Type      MessageType `json:"type"`
	SessionID string      `json:"sessionId"`
	ChannelID string      `json:"channelId"`
	RequestID string      `json:"requestId"`
}

// Stream carries a single Claude event plus a sequence number.
type Stream struct {
	Type           string          `json:"type"`
	TaskID         string          `json:"taskId,omitempty"`
	AttemptID      string          `json:"attemptId,omitempty"`
	AttemptNumber  int             `json:"attemptNumber,omitempty"`
	SessionID      string          `json:"sessionId"`
	ChannelID      string          `json:"channelId"`
	SequenceNumber int64           `json:"sequenceNumber"`
	Event          json.RawMessage `json:"event"`
	RequestID      string          `json:"requestId,omitempty"`
	Traceparent    string          `json:"traceparent,omitempty"` // W3C trace context
}

// TaskStarted signals the daemon began processing a task.
type TaskStarted struct {
	Type          string `json:"type"`
	TaskID        string `json:"taskId"`
	AttemptID     string `json:"attemptId,omitempty"`
	AttemptNumber int    `json:"attemptNumber,omitempty"`
	SessionID     string `json:"sessionId"`
	ChannelID     string `json:"channelId"`
	StartedAt     string `json:"startedAt"`
	RequestID     string `json:"requestId,omitempty"`
	Traceparent   string `json:"traceparent,omitempty"` // W3C trace context
}

// TaskComplete reports final result metadata.
type TaskComplete struct {
	Type            string `json:"type"`
	TaskID          string `json:"taskId"`
	AttemptID       string `json:"attemptId,omitempty"`
	AttemptNumber   int    `json:"attemptNumber,omitempty"`
	SessionID       string `json:"sessionId"`
	ChannelID       string `json:"channelId"`
	ClaudeSessionID string `json:"claudeSessionId"`
	InputTokens     int64  `json:"inputTokens"`
	OutputTokens    int64  `json:"outputTokens"`
	CostUSD         string `json:"costUsd"`
	DurationMs      int    `json:"durationMs"`
	RequestID       string `json:"requestId,omitempty"`
	Traceparent     string `json:"traceparent,omitempty"` // W3C trace context
}

// TaskError reports a failure.
type TaskError struct {
	Type          string `json:"type"`
	TaskID        string `json:"taskId"`
	AttemptID     string `json:"attemptId,omitempty"`
	AttemptNumber int    `json:"attemptNumber,omitempty"`
	SessionID     string `json:"sessionId"`
	ChannelID     string `json:"channelId"`
	Error         string `json:"error"`
	FailureCode   string `json:"failureCode,omitempty"`
	Retryable     bool   `json:"retryable,omitempty"`
	UserMessage   string `json:"userMessage,omitempty"`
	RequestID     string `json:"requestId,omitempty"`
	Traceparent   string `json:"traceparent,omitempty"` // W3C trace context
}

// TaskCancelled tells the relay/browser that a task was interrupted by the user.
type TaskCancelled struct {
	Type          string `json:"type"`
	TaskID        string `json:"taskId"`
	AttemptID     string `json:"attemptId,omitempty"`
	AttemptNumber int    `json:"attemptNumber,omitempty"`
	SessionID     string `json:"sessionId"`
	ChannelID     string `json:"channelId"`
	FailureCode   string `json:"failureCode,omitempty"`
	Retryable     bool   `json:"retryable,omitempty"`
	UserMessage   string `json:"userMessage,omitempty"`
	RequestID     string `json:"requestId,omitempty"`
	Traceparent   string `json:"traceparent,omitempty"` // W3C trace context
}

// PermissionRequest is Claude asking for tool approval.
type PermissionRequest struct {
	Type          string          `json:"type"`
	TaskID        string          `json:"taskId,omitempty"`
	AttemptID     string          `json:"attemptId,omitempty"`
	AttemptNumber int             `json:"attemptNumber,omitempty"`
	SessionID     string          `json:"sessionId"`
	ChannelID     string          `json:"channelId"`
	RequestID     string          `json:"requestId"`
	ToolName      string          `json:"toolName"`
	ToolInput     json.RawMessage `json:"toolInput"`
}

// QuestionOption is a structured answer choice for AskUserQuestion.
type QuestionOption struct {
	Label       string `json:"label"`
	Description string `json:"description,omitempty"`
	Preview     string `json:"preview,omitempty"`
}

// Question is Claude asking the user for input.
type Question struct {
	Type          string           `json:"type"`
	TaskID        string           `json:"taskId,omitempty"`
	AttemptID     string           `json:"attemptId,omitempty"`
	AttemptNumber int              `json:"attemptNumber,omitempty"`
	SessionID     string           `json:"sessionId"`
	ChannelID     string           `json:"channelId"`
	RequestID     string           `json:"requestId"`
	Question      string           `json:"question"`
	Header        string           `json:"header,omitempty"`
	MultiSelect   bool             `json:"multiSelect,omitempty"`
	Options       []QuestionOption `json:"options,omitempty"`
}

type ContextStats struct {
	Type                 MessageType `json:"type"`
	SessionID            string      `json:"sessionId"`
	ChannelID            string      `json:"channelId"`
	RequestID            string      `json:"requestId,omitempty"`
	Tokens               *int64      `json:"tokens"`
	ContextWindow        int64       `json:"contextWindow"`
	Percent              *float64    `json:"percent"`
	ReserveTokens        int64       `json:"reserveTokens"`
	KeepRecentTokens     int64       `json:"keepRecentTokens"`
	AutoThresholdPercent float64     `json:"autoThresholdPercent"`
	Source               string      `json:"source"`
	ObservedAt           time.Time   `json:"observedAt"`
}

type CompactStatus struct {
	Type                 MessageType            `json:"type"`
	SessionID            string                 `json:"sessionId"`
	ChannelID            string                 `json:"channelId"`
	RequestID            string                 `json:"requestId"`
	Status               CompactLifecycleStatus `json:"status"`
	Reason               CompactReason          `json:"reason"`
	Instructions         string                 `json:"instructions,omitempty"`
	TokensBefore         *int64                 `json:"tokensBefore"`
	TokensAfter          *int64                 `json:"tokensAfter"`
	ContextWindow        int64                  `json:"contextWindow"`
	ReserveTokens        int64                  `json:"reserveTokens"`
	KeepRecentTokens     int64                  `json:"keepRecentTokens"`
	AutoThresholdPercent float64                `json:"autoThresholdPercent"`
	Summary              string                 `json:"summary,omitempty"`
	FirstKeptEntryID     string                 `json:"firstKeptEntryId,omitempty"`
	Error                string                 `json:"error,omitempty"`
	Source               string                 `json:"source"`
	ObservedAt           time.Time              `json:"observedAt"`
}

// Heartbeat is the daemon's 30s health pulse.
type Heartbeat struct {
	Type          string `json:"type"`
	MachineID     string `json:"machineId"`
	DaemonVersion string `json:"daemonVersion"`
	Status        string `json:"status"`
	Timestamp     string `json:"timestamp"`
}

// BrowseEntry is one row in a directory listing.
type BrowseEntry struct {
	Name        string `json:"name"`
	Path        string `json:"path"`
	IsDirectory bool   `json:"isDirectory"`
	Size        int64  `json:"size"`
	ModifiedAt  string `json:"modifiedAt"`
}

// BrowseDirResult is the daemon's response to a BrowseDir request.
type BrowseDirResult struct {
	Type       string        `json:"type"`
	RequestID  string        `json:"requestId"`
	ChannelID  string        `json:"channelId"`
	OK         bool          `json:"ok"`
	Entries    []BrowseEntry `json:"entries,omitempty"`
	HasMore    bool          `json:"hasMore,omitempty"`
	NextCursor string        `json:"nextCursor,omitempty"`
	Error      string        `json:"error,omitempty"`
}

// MkDir asks the daemon to create a directory.
type MkDir struct {
	Type      string `json:"type"`
	RequestID string `json:"requestId"`
	ChannelID string `json:"channelId"`
	MachineID string `json:"machineId"`
	Path      string `json:"path"`
}

// Skill is a locally installed agent skill discovered on the daemon machine.
type Skill struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Path        string `json:"path"`
	Scope       string `json:"scope"`
}

// ListSkills asks the daemon to list available skills for a working directory.
type ListSkills struct {
	Type      string `json:"type"`
	RequestID string `json:"requestId"`
	ChannelID string `json:"channelId"`
	MachineID string `json:"machineId"`
	CWD       string `json:"cwd"`
}

// ListSkillsResult is the daemon's response to a ListSkills request.
type ListSkillsResult struct {
	Type      string  `json:"type"`
	RequestID string  `json:"requestId"`
	ChannelID string  `json:"channelId"`
	OK        bool    `json:"ok"`
	Skills    []Skill `json:"skills,omitempty"`
	Error     string  `json:"error,omitempty"`
}

// MkDirResult is the daemon's response to a MkDir request.
type MkDirResult struct {
	Type      string `json:"type"`
	RequestID string `json:"requestId"`
	ChannelID string `json:"channelId"`
	OK        bool   `json:"ok"`
	Error     string `json:"error,omitempty"`
}

// ReadFileResult is the daemon's response to a ReadFile request.
type ReadFileResult struct {
	Type      string `json:"type"`
	RequestID string `json:"requestId"`
	ChannelID string `json:"channelId"`
	OK        bool   `json:"ok"`
	Content   string `json:"content,omitempty"`
	Truncated bool   `json:"truncated,omitempty"`
	Error     string `json:"error,omitempty"`
}

// TerminalOpen requests a chat-scoped PTY terminal.
type TerminalOpen struct {
	Type          string `json:"type"`
	RequestID     string `json:"requestId"`
	TerminalID    string `json:"terminalId,omitempty"`
	SessionID     string `json:"sessionId"`
	ChannelID     string `json:"channelId"`
	Token         string `json:"token,omitempty"`
	CWD           string `json:"cwd,omitempty"`
	Cols          int    `json:"cols"`
	Rows          int    `json:"rows"`
	IdleTimeoutMs int    `json:"idleTimeoutMs,omitempty"`
	MaxLifetimeMs int    `json:"maxLifetimeMs,omitempty"`
}

// TerminalOpened confirms a terminal has started.
type TerminalOpened struct {
	Type       string `json:"type"`
	RequestID  string `json:"requestId"`
	TerminalID string `json:"terminalId"`
	SessionID  string `json:"sessionId"`
	ChannelID  string `json:"channelId"`
	Shell      string `json:"shell"`
	CWD        string `json:"cwd"`
	StartedAt  string `json:"startedAt"`
}

// TerminalInput carries browser input bytes as base64.
type TerminalInput struct {
	Type       string `json:"type"`
	TerminalID string `json:"terminalId"`
	ChannelID  string `json:"channelId"`
	DataBase64 string `json:"dataBase64"`
}

// TerminalOutput carries terminal output bytes as base64.
type TerminalOutput struct {
	Type       string `json:"type"`
	TerminalID string `json:"terminalId"`
	SessionID  string `json:"sessionId"`
	ChannelID  string `json:"channelId"`
	Seq        int64  `json:"seq"`
	DataBase64 string `json:"dataBase64"`
}

// TerminalSnapshot carries bounded scrollback bytes as base64.
type TerminalSnapshot struct {
	Type       string `json:"type"`
	TerminalID string `json:"terminalId"`
	SessionID  string `json:"sessionId"`
	ChannelID  string `json:"channelId"`
	Seq        int64  `json:"seq"`
	DataBase64 string `json:"dataBase64"`
}

// TerminalResize resizes the PTY.
type TerminalResize struct {
	Type       string `json:"type"`
	TerminalID string `json:"terminalId"`
	ChannelID  string `json:"channelId"`
	Cols       int    `json:"cols"`
	Rows       int    `json:"rows"`
}

// TerminalClose terminates the PTY.
type TerminalClose struct {
	Type       string `json:"type"`
	TerminalID string `json:"terminalId"`
	ChannelID  string `json:"channelId"`
}

// TerminalExit reports terminal process completion.
type TerminalExit struct {
	Type       string `json:"type"`
	TerminalID string `json:"terminalId"`
	SessionID  string `json:"sessionId"`
	ChannelID  string `json:"channelId"`
	ExitCode   *int   `json:"exitCode,omitempty"`
	Signal     string `json:"signal,omitempty"`
	Reason     string `json:"reason"`
	EndedAt    string `json:"endedAt"`
}

// TerminalError reports a terminal lifecycle or authorization error.
type TerminalError struct {
	Type       string `json:"type"`
	RequestID  string `json:"requestId,omitempty"`
	TerminalID string `json:"terminalId,omitempty"`
	SessionID  string `json:"sessionId,omitempty"`
	ChannelID  string `json:"channelId"`
	Error      string `json:"error"`
}

type AgentTerminalReadiness struct {
	State       string `json:"state"`
	Source      string `json:"source,omitempty"`
	MatchedText string `json:"matchedText,omitempty"`
	ReadyAt     string `json:"readyAt,omitempty"`
	TimeoutMs   int    `json:"timeoutMs,omitempty"`
}

type AgentTerminalPort struct {
	Host string `json:"host"`
	Port int    `json:"port"`
	URL  string `json:"url"`
}

type AgentTerminalStarted struct {
	Type           string                 `json:"type"`
	JobID          string                 `json:"jobId"`
	TerminalID     string                 `json:"terminalId"`
	SessionID      string                 `json:"sessionId"`
	ChannelID      string                 `json:"channelId"`
	TaskID         string                 `json:"taskId,omitempty"`
	ToolCallID     string                 `json:"toolCallId,omitempty"`
	ProjectID      string                 `json:"projectId"`
	CommandPreview string                 `json:"commandPreview"`
	Title          string                 `json:"title"`
	CWD            string                 `json:"cwd"`
	Status         string                 `json:"status"`
	Readiness      AgentTerminalReadiness `json:"readiness"`
	Ports          []AgentTerminalPort    `json:"ports,omitempty"`
	URLs           []string               `json:"urls,omitempty"`
	Seq            int64                  `json:"seq,omitempty"`
	StartedAt      string                 `json:"startedAt"`
}

type AgentTerminalUpdated struct {
	Type       string                 `json:"type"`
	JobID      string                 `json:"jobId"`
	TerminalID string                 `json:"terminalId"`
	SessionID  string                 `json:"sessionId"`
	ChannelID  string                 `json:"channelId"`
	Status     string                 `json:"status"`
	Readiness  AgentTerminalReadiness `json:"readiness"`
	Ports      []AgentTerminalPort    `json:"ports,omitempty"`
	URLs       []string               `json:"urls,omitempty"`
	Seq        int64                  `json:"seq,omitempty"`
	UpdatedAt  string                 `json:"updatedAt"`
}

type AgentTerminalAttach struct {
	Type       string `json:"type"`
	TerminalID string `json:"terminalId"`
	ChannelID  string `json:"channelId"`
}

type AgentTerminalSnapshotRequest struct {
	Type       string `json:"type"`
	TerminalID string `json:"terminalId"`
	ChannelID  string `json:"channelId"`
}

type BrowserSessionOpen struct {
	Type       MessageType `json:"type"`
	RequestID  string      `json:"requestId"`
	GrantID    string      `json:"grantId"`
	SessionID  string      `json:"sessionId"`
	ProjectID  string      `json:"projectId"`
	TaskID     string      `json:"taskId"`
	ChannelID  string      `json:"channelId"`
	MachineID  string      `json:"machineId"`
	IdentityID string      `json:"identityId,omitempty"`
	Mode       string      `json:"mode"`
	InitialURL string      `json:"initialUrl,omitempty"`
	BridgeMode string      `json:"bridgeMode,omitempty"`
	PreviewID  string      `json:"sourcePreviewId,omitempty"`
	ExpiresAt  string      `json:"expiresAt"`
}

type BrowserSessionOpened struct {
	Type      MessageType `json:"type"`
	RequestID string      `json:"requestId"`
	BrowserID string      `json:"browserId"`
	GrantID   string      `json:"grantId"`
	SessionID string      `json:"sessionId"`
	ChannelID string      `json:"channelId"`
	URL       string      `json:"url,omitempty"`
	Title     string      `json:"title,omitempty"`
	OpenedAt  string      `json:"openedAt"`
}

type BrowserSessionClose struct {
	Type      MessageType `json:"type"`
	BrowserID string      `json:"browserId"`
	GrantID   string      `json:"grantId"`
	SessionID string      `json:"sessionId"`
	ChannelID string      `json:"channelId"`
	Reason    string      `json:"reason,omitempty"`
}

type BrowserSessionClosed struct {
	Type      MessageType `json:"type"`
	BrowserID string      `json:"browserId"`
	SessionID string      `json:"sessionId"`
	ChannelID string      `json:"channelId"`
	Reason    string      `json:"reason,omitempty"`
	ClosedAt  string      `json:"closedAt"`
}

type BrowserSessionError struct {
	Type      MessageType `json:"type"`
	BrowserID string      `json:"browserId,omitempty"`
	RequestID string      `json:"requestId,omitempty"`
	SessionID string      `json:"sessionId"`
	ChannelID string      `json:"channelId"`
	Code      string      `json:"code"`
	Message   string      `json:"message"`
}

type BrowserFrame struct {
	Type              MessageType `json:"type"`
	BrowserID         string      `json:"browserId"`
	SessionID         string      `json:"sessionId"`
	ChannelID         string      `json:"channelId"`
	Seq               int64       `json:"seq"`
	ContentType       string      `json:"contentType"`
	DataBase64        string      `json:"dataBase64,omitempty"`
	FrameRef          string      `json:"frameRef,omitempty"`
	Width             int         `json:"width"`
	Height            int         `json:"height"`
	ViewportWidth     int         `json:"viewportWidth,omitempty"`
	ViewportHeight    int         `json:"viewportHeight,omitempty"`
	DevicePixelRatio  float64     `json:"devicePixelRatio,omitempty"`
	CapturedAt        string      `json:"capturedAt"`
	DroppedPriorCount int         `json:"droppedPriorCount,omitempty"`
}

type BrowserRefs struct {
	Type       MessageType  `json:"type"`
	BrowserID  string       `json:"browserId"`
	SessionID  string       `json:"sessionId"`
	ChannelID  string       `json:"channelId"`
	Version    int          `json:"version"`
	Refs       []BrowserRef `json:"refs"`
	CapturedAt string       `json:"capturedAt"`
}

type BrowserRef struct {
	Ref  string  `json:"ref"`
	Key  string  `json:"key"`
	Role string  `json:"role"`
	Name string  `json:"name,omitempty"`
	X    float64 `json:"x"`
	Y    float64 `json:"y"`
	W    float64 `json:"w"`
	H    float64 `json:"h"`
}

type BrowserCursor struct {
	Type      MessageType `json:"type"`
	BrowserID string      `json:"browserId"`
	SessionID string      `json:"sessionId"`
	ChannelID string      `json:"channelId"`
	Owner     string      `json:"owner"`
	X         float64     `json:"x"`
	Y         float64     `json:"y"`
	UpdatedAt string      `json:"updatedAt"`
}

type BrowserNavigation struct {
	Type      MessageType `json:"type"`
	BrowserID string      `json:"browserId"`
	SessionID string      `json:"sessionId"`
	ChannelID string      `json:"channelId"`
	URL       string      `json:"url"`
	Title     string      `json:"title,omitempty"`
	StartedAt string      `json:"startedAt,omitempty"`
	EndedAt   string      `json:"endedAt,omitempty"`
}

type BrowserAction struct {
	Type      MessageType     `json:"type"`
	BrowserID string          `json:"browserId"`
	SessionID string          `json:"sessionId"`
	ChannelID string          `json:"channelId"`
	TaskID    string          `json:"taskId,omitempty"`
	ToolUseID string          `json:"toolUseId,omitempty"`
	Summary   string          `json:"summary"`
	Status    string          `json:"status"`
	Metadata  json.RawMessage `json:"metadata,omitempty"`
	At        string          `json:"at"`
}

type BrowserToolCall struct {
	Type       MessageType     `json:"type"`
	BrowserID  string          `json:"browserId"`
	GrantID    string          `json:"grantId"`
	TaskID     string          `json:"taskId"`
	ToolUseID  string          `json:"toolUseId"`
	Method     string          `json:"method"`
	ParamsJSON json.RawMessage `json:"paramsJson,omitempty"`
}

type BrowserToolResult struct {
	Type                 MessageType     `json:"type"`
	BrowserID            string          `json:"browserId"`
	GrantID              string          `json:"grantId"`
	SessionID            string          `json:"sessionId,omitempty"`
	ChannelID            string          `json:"channelId,omitempty"`
	TaskID               string          `json:"taskId"`
	ToolUseID            string          `json:"toolUseId"`
	OK                   bool            `json:"ok"`
	ResultJSON           json.RawMessage `json:"resultJson,omitempty"`
	Error                string          `json:"error,omitempty"`
	ErrorCode            string          `json:"errorCode,omitempty"`
	RecoveryHint         string          `json:"recoveryHint,omitempty"`
	Sensitivity          string          `json:"sensitivity,omitempty"`
	RedactionStatus      string          `json:"redactionStatus,omitempty"`
	LocalArtifactPointer string          `json:"localArtifactPointer,omitempty"`
}

type BrowserToolCallStarted struct {
	Type      MessageType     `json:"type"`
	BrowserID string          `json:"browserId,omitempty"`
	GrantID   string          `json:"grantId,omitempty"`
	SessionID string          `json:"sessionId"`
	ChannelID string          `json:"channelId"`
	TaskID    string          `json:"taskId,omitempty"`
	ToolUseID string          `json:"toolUseId,omitempty"`
	Method    string          `json:"method"`
	Category  string          `json:"category,omitempty"`
	Summary   string          `json:"summary"`
	Intent    string          `json:"intent,omitempty"`
	Metadata  json.RawMessage `json:"metadata,omitempty"`
	At        string          `json:"at"`
}

type BrowserToolCallUpdated struct {
	Type      MessageType     `json:"type"`
	BrowserID string          `json:"browserId,omitempty"`
	GrantID   string          `json:"grantId,omitempty"`
	SessionID string          `json:"sessionId"`
	ChannelID string          `json:"channelId"`
	TaskID    string          `json:"taskId,omitempty"`
	ToolUseID string          `json:"toolUseId,omitempty"`
	Status    string          `json:"status"`
	Summary   string          `json:"summary,omitempty"`
	Metadata  json.RawMessage `json:"metadata,omitempty"`
	At        string          `json:"at"`
}

type BrowserArtifactCreated struct {
	Type       MessageType     `json:"type"`
	BrowserID  string          `json:"browserId"`
	GrantID    string          `json:"grantId,omitempty"`
	SessionID  string          `json:"sessionId"`
	ChannelID  string          `json:"channelId"`
	TaskID     string          `json:"taskId,omitempty"`
	ArtifactID string          `json:"artifactId"`
	Kind       string          `json:"kind"`
	Title      string          `json:"title,omitempty"`
	URL        string          `json:"url,omitempty"`
	Metadata   json.RawMessage `json:"metadata,omitempty"`
	CreatedAt  string          `json:"createdAt"`
}

type BrowserControlClaim struct {
	Type           MessageType `json:"type"`
	BrowserID      string      `json:"browserId"`
	SessionID      string      `json:"sessionId"`
	ChannelID      string      `json:"channelId"`
	Owner          string      `json:"owner"`
	Reason         string      `json:"reason,omitempty"`
	ControlVersion int64       `json:"controlVersion,omitempty"`
}

type BrowserControlRelease struct {
	Type           MessageType `json:"type"`
	BrowserID      string      `json:"browserId"`
	SessionID      string      `json:"sessionId"`
	ChannelID      string      `json:"channelId"`
	Owner          string      `json:"owner"`
	Reason         string      `json:"reason,omitempty"`
	ControlVersion int64       `json:"controlVersion,omitempty"`
}

type BrowserUserInput struct {
	Type             MessageType `json:"type"`
	InputID          string      `json:"inputId,omitempty"`
	BrowserID        string      `json:"browserId"`
	SessionID        string      `json:"sessionId"`
	ChannelID        string      `json:"channelId"`
	Owner            string      `json:"owner"`
	Kind             string      `json:"kind"`
	X                *float64    `json:"x,omitempty"`
	Y                *float64    `json:"y,omitempty"`
	Text             string      `json:"text,omitempty"`
	Key              string      `json:"key,omitempty"`
	DeltaX           *float64    `json:"deltaX,omitempty"`
	DeltaY           *float64    `json:"deltaY,omitempty"`
	FrameSeq         int64       `json:"frameSeq,omitempty"`
	ControlVersion   int64       `json:"controlVersion,omitempty"`
	CoordinateSpace  string      `json:"coordinateSpace,omitempty"`
	ViewportWidth    int         `json:"viewportWidth,omitempty"`
	ViewportHeight   int         `json:"viewportHeight,omitempty"`
	FrameWidth       int         `json:"frameWidth,omitempty"`
	FrameHeight      int         `json:"frameHeight,omitempty"`
	DevicePixelRatio float64     `json:"devicePixelRatio,omitempty"`
	RenderedLeft     float64     `json:"renderedLeft"`
	RenderedTop      float64     `json:"renderedTop"`
	RenderedWidth    float64     `json:"renderedWidth"`
	RenderedHeight   float64     `json:"renderedHeight"`
}

type BrowserUserInputAck struct {
	Type           MessageType `json:"type"`
	BrowserID      string      `json:"browserId"`
	SessionID      string      `json:"sessionId"`
	ChannelID      string      `json:"channelId"`
	InputID        string      `json:"inputId,omitempty"`
	Accepted       bool        `json:"accepted"`
	Reason         string      `json:"reason,omitempty"`
	ControlVersion int64       `json:"controlVersion,omitempty"`
	AckedAt        string      `json:"ackedAt"`
}

type BrowserTransportStatus struct {
	Type              MessageType `json:"type"`
	BrowserID         string      `json:"browserId"`
	SessionID         string      `json:"sessionId"`
	ChannelID         string      `json:"channelId"`
	Status            string      `json:"status"`
	QueueDepth        int         `json:"queueDepth,omitempty"`
	DroppedFrameCount int64       `json:"droppedFrameCount,omitempty"`
	MaxFrameBytes     int         `json:"maxFrameBytes,omitempty"`
	At                string      `json:"at"`
}

type BrowserBridgeAccessOpen struct {
	Type        MessageType `json:"type"`
	RequestID   string      `json:"requestId"`
	PreviewID   string      `json:"previewId"`
	GrantID     string      `json:"grantId"`
	BrowserID   string      `json:"browserId,omitempty"`
	SessionID   string      `json:"sessionId"`
	ChannelID   string      `json:"channelId"`
	MachineID   string      `json:"machineId"`
	BridgeMode  string      `json:"bridgeMode"`
	RequestedAt string      `json:"requestedAt"`
}

type BrowserBridgeAccessOpened struct {
	Type       MessageType `json:"type"`
	RequestID  string      `json:"requestId"`
	PreviewID  string      `json:"previewId"`
	GrantID    string      `json:"grantId"`
	BrowserID  string      `json:"browserId,omitempty"`
	SessionID  string      `json:"sessionId"`
	ChannelID  string      `json:"channelId"`
	BridgeMode string      `json:"bridgeMode"`
	URL        string      `json:"url"`
	ExpiresAt  string      `json:"expiresAt"`
}

type BrowserBridgeAccessClose struct {
	Type      MessageType `json:"type"`
	PreviewID string      `json:"previewId"`
	GrantID   string      `json:"grantId"`
	BrowserID string      `json:"browserId,omitempty"`
	SessionID string      `json:"sessionId"`
	ChannelID string      `json:"channelId"`
	Reason    string      `json:"reason,omitempty"`
}

type BrowserSensitiveActionRequest struct {
	Type      MessageType `json:"type"`
	BrowserID string      `json:"browserId"`
	RequestID string      `json:"requestId"`
	SessionID string      `json:"sessionId"`
	ChannelID string      `json:"channelId"`
	TaskID    string      `json:"taskId"`
	ToolUseID string      `json:"toolUseId"`
	Category  string      `json:"category"`
	Summary   string      `json:"summary"`
	Origin    string      `json:"origin,omitempty"`
	ExpiresAt string      `json:"expiresAt"`
}

type BrowserSensitiveActionResponse struct {
	Type      MessageType `json:"type"`
	BrowserID string      `json:"browserId"`
	RequestID string      `json:"requestId"`
	SessionID string      `json:"sessionId"`
	ChannelID string      `json:"channelId"`
	Approved  bool        `json:"approved"`
	Reason    string      `json:"reason,omitempty"`
}

// HelloCapabilities describes optional daemon protocol support.
type HelloCapabilities struct {
	Stop                           bool   `json:"stop,omitempty"`
	Terminal                       bool   `json:"terminal,omitempty"`
	AgentTerminalJobs              bool   `json:"agentTerminalJobs,omitempty"`
	ContextRefs                    bool   `json:"contextRefs,omitempty"`
	PreviewTunnel                  bool   `json:"previewTunnel,omitempty"`
	PreviewMaxFrameBytes           int    `json:"previewMaxFrameBytes,omitempty"`
	PreviewChunkBytes              int    `json:"previewChunkBytes,omitempty"`
	PreviewWebSocketProtocols      bool   `json:"previewWebSocketProtocols,omitempty"`
	LocalServerDetection           bool   `json:"localServerDetection,omitempty"`
	Skills                         bool   `json:"skills,omitempty"`
	BrowserSessions                bool   `json:"browserSessions,omitempty"`
	BrowserFrameStream             bool   `json:"browserFrameStream,omitempty"`
	BrowserUserControl             bool   `json:"browserUserControl,omitempty"`
	BrowserIdentities              bool   `json:"browserIdentities,omitempty"`
	BrowserSensitiveActionApproval bool   `json:"browserSensitiveActionApproval,omitempty"`
	BrowserMaxFrameBytes           int    `json:"browserMaxFrameBytes,omitempty"`
	BrowserRuntimeInstalled        bool   `json:"browserRuntimeInstalled,omitempty"`
	BrowserRuntimeVersion          string `json:"browserRuntimeVersion,omitempty"`
	BrowserRuntimeMinVersion       string `json:"browserRuntimeMinVersion,omitempty"`
	BrowserRuntimeMinVersionOK     bool   `json:"browserRuntimeMinVersionOk,omitempty"`
	BrowserRuntimePath             string `json:"browserRuntimePath,omitempty"`
	BrowserRuntimeErrorCode        string `json:"browserRuntimeErrorCode,omitempty"`
	BrowserRuntimeErrorMessage     string `json:"browserRuntimeErrorMessage,omitempty"`
	BrowserCloudMethodsVersion     int    `json:"browserCloudMethodsVersion,omitempty"`
	BrowserChromeAvailable         bool   `json:"browserChromeAvailable,omitempty"`
}

// Hello is the first frame sent by the daemon after connecting.
type Hello struct {
	Type          string             `json:"type"`
	MachineID     string             `json:"machineId"`
	DaemonVersion string             `json:"daemonVersion"`
	OS            string             `json:"os"`
	Arch          string             `json:"arch"`
	ActiveTasks   []string           `json:"activeTasks,omitempty"`
	Capabilities  *HelloCapabilities `json:"capabilities,omitempty"`
}

// Welcome is the relay's response to Hello.
type Welcome struct {
	Type                string `json:"type"`
	LatestDaemonVersion string `json:"latestDaemonVersion,omitempty"`
}

// MachineStatus is pushed to all connected browsers when machine presence changes.
type MachineStatus struct {
	Type          string `json:"type"`
	MachineID     string `json:"machineId"`
	State         string `json:"state"`
	PreviousState string `json:"previousState,omitempty"`
	Reason        string `json:"reason,omitempty"`
	OccurredAt    string `json:"occurredAt"`
}

type PreviewOpen struct {
	Type       string `json:"type"`
	RequestID  string `json:"requestId"`
	PreviewID  string `json:"previewId"`
	SessionID  string `json:"sessionId"`
	ChannelID  string `json:"channelId"`
	MachineID  string `json:"machineId"`
	TargetHost string `json:"targetHost"`
	TargetPort int    `json:"targetPort"`
	ExpiresAt  string `json:"expiresAt"`
}

type PreviewOpenResult struct {
	Type      string `json:"type"`
	RequestID string `json:"requestId"`
	PreviewID string `json:"previewId"`
	OK        bool   `json:"ok"`
	ErrorCode string `json:"errorCode,omitempty"`
	Message   string `json:"message,omitempty"`
}

type PreviewClose struct {
	Type      string `json:"type"`
	PreviewID string `json:"previewId"`
	Reason    string `json:"reason"`
}

type PreviewHTTPRequest struct {
	Type      string              `json:"type"`
	RequestID string              `json:"requestId"`
	StreamID  string              `json:"streamId"`
	PreviewID string              `json:"previewId"`
	Method    string              `json:"method"`
	Path      string              `json:"path"`
	Headers   map[string][]string `json:"headers,omitempty"`
}

type PreviewHTTPResponseHead struct {
	Type       string              `json:"type"`
	RequestID  string              `json:"requestId"`
	StreamID   string              `json:"streamId"`
	PreviewID  string              `json:"previewId"`
	StatusCode int                 `json:"statusCode"`
	Headers    map[string][]string `json:"headers,omitempty"`
}

type PreviewStreamChunk struct {
	Type       string `json:"type"`
	StreamID   string `json:"streamId"`
	Sequence   int64  `json:"sequence"`
	BodyBase64 string `json:"bodyBase64"`
	Final      bool   `json:"final"`
}

type PreviewStreamCancel struct {
	Type     string `json:"type"`
	StreamID string `json:"streamId"`
	Reason   string `json:"reason"`
}

type PreviewWebSocketOpen struct {
	Type      string              `json:"type"`
	StreamID  string              `json:"streamId"`
	PreviewID string              `json:"previewId"`
	Path      string              `json:"path"`
	Headers   map[string][]string `json:"headers,omitempty"`
	Protocols []string            `json:"protocols,omitempty"`
}

type PreviewWebSocketOpenResult struct {
	Type      string `json:"type"`
	StreamID  string `json:"streamId"`
	PreviewID string `json:"previewId"`
	OK        bool   `json:"ok"`
	Protocol  string `json:"protocol,omitempty"`
	ErrorCode string `json:"errorCode,omitempty"`
	Message   string `json:"message,omitempty"`
}

type PreviewWebSocketData struct {
	Type       string `json:"type"`
	StreamID   string `json:"streamId"`
	Sequence   int64  `json:"sequence"`
	IsBinary   bool   `json:"isBinary"`
	BodyBase64 string `json:"bodyBase64"`
}

type PreviewWebSocketClose struct {
	Type     string `json:"type"`
	StreamID string `json:"streamId"`
	Code     int    `json:"code,omitempty"`
	Reason   string `json:"reason,omitempty"`
}

type LocalServerDetected struct {
	Type       string `json:"type"`
	SessionID  string `json:"sessionId"`
	ChannelID  string `json:"channelId"`
	TaskID     string `json:"taskId,omitempty"`
	ToolUseID  string `json:"toolUseId,omitempty"`
	Host       string `json:"host"`
	Port       int    `json:"port"`
	URL        string `json:"url"`
	Command    string `json:"command,omitempty"`
	Source     string `json:"source"`
	DetectedAt string `json:"detectedAt"`
}
