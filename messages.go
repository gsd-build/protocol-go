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
	MsgTypeStop                            = "stop"
	MsgTypePermissionResponse              = "permissionResponse"
	MsgTypeQuestionResponse                = "questionResponse"
	MsgTypeBrowseDir                       = "browseDir"
	MsgTypeReadFile                        = "readFile"
	MsgTypeMkDir                           = "mkDir"
	MsgTypeMkDirResult                     = "mkDirResult"
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
)

type ContextRef struct {
	Kind       string `json:"kind"`
	Path       string `json:"path"`
	Name       string `json:"name"`
	Size       *int64 `json:"size,omitempty"`
	ModifiedAt string `json:"modifiedAt,omitempty"`
}

// Task is sent from the browser to the daemon to dispatch a user message.
type Task struct {
	Type            string       `json:"type"`
	TaskID          string       `json:"taskId"`
	SessionID       string       `json:"sessionId"`
	ChannelID       string       `json:"channelId"`
	Prompt          string       `json:"prompt"`
	Engine          string       `json:"engine,omitempty"` // "pi"; empty defaults to pi
	Model           string       `json:"model"`
	Effort          string       `json:"effort"`
	PermissionMode  string       `json:"permissionMode"`
	CWD             string       `json:"cwd"`
	ClaudeSessionID string       `json:"claudeSessionId,omitempty"` // passed to --resume
	RequestID       string       `json:"requestId,omitempty"`
	Traceparent     string       `json:"traceparent,omitempty"` // W3C trace context
	ImageURLs       []string     `json:"imageUrls,omitempty"`   // user-attached image URLs
	ContextRefs     []ContextRef `json:"contextRefs,omitempty"`
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
	SessionID      string          `json:"sessionId"`
	ChannelID      string          `json:"channelId"`
	SequenceNumber int64           `json:"sequenceNumber"`
	Event          json.RawMessage `json:"event"`
	RequestID      string          `json:"requestId,omitempty"`
	Traceparent    string          `json:"traceparent,omitempty"` // W3C trace context
}

// TaskStarted signals the daemon began processing a task.
type TaskStarted struct {
	Type        string `json:"type"`
	TaskID      string `json:"taskId"`
	SessionID   string `json:"sessionId"`
	ChannelID   string `json:"channelId"`
	StartedAt   string `json:"startedAt"`
	RequestID   string `json:"requestId,omitempty"`
	Traceparent string `json:"traceparent,omitempty"` // W3C trace context
}

// TaskComplete reports final result metadata.
type TaskComplete struct {
	Type            string `json:"type"`
	TaskID          string `json:"taskId"`
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
	Type        string `json:"type"`
	TaskID      string `json:"taskId"`
	SessionID   string `json:"sessionId"`
	ChannelID   string `json:"channelId"`
	Error       string `json:"error"`
	RequestID   string `json:"requestId,omitempty"`
	Traceparent string `json:"traceparent,omitempty"` // W3C trace context
}

// TaskCancelled tells the relay/browser that a task was interrupted by the user.
type TaskCancelled struct {
	Type        string `json:"type"`
	TaskID      string `json:"taskId"`
	SessionID   string `json:"sessionId"`
	ChannelID   string `json:"channelId"`
	RequestID   string `json:"requestId,omitempty"`
	Traceparent string `json:"traceparent,omitempty"` // W3C trace context
}

// PermissionRequest is Claude asking for tool approval.
type PermissionRequest struct {
	Type      string          `json:"type"`
	SessionID string          `json:"sessionId"`
	ChannelID string          `json:"channelId"`
	RequestID string          `json:"requestId"`
	ToolName  string          `json:"toolName"`
	ToolInput json.RawMessage `json:"toolInput"`
}

// QuestionOption is a structured answer choice for AskUserQuestion.
type QuestionOption struct {
	Label       string `json:"label"`
	Description string `json:"description,omitempty"`
	Preview     string `json:"preview,omitempty"`
}

// Question is Claude asking the user for input.
type Question struct {
	Type        string           `json:"type"`
	SessionID   string           `json:"sessionId"`
	ChannelID   string           `json:"channelId"`
	RequestID   string           `json:"requestId"`
	Question    string           `json:"question"`
	Header      string           `json:"header,omitempty"`
	MultiSelect bool             `json:"multiSelect,omitempty"`
	Options     []QuestionOption `json:"options,omitempty"`
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

// HelloCapabilities describes optional daemon protocol support.
type HelloCapabilities struct {
	Stop                      bool `json:"stop,omitempty"`
	Terminal                  bool `json:"terminal,omitempty"`
	ContextRefs               bool `json:"contextRefs,omitempty"`
	PreviewTunnel             bool `json:"previewTunnel,omitempty"`
	PreviewMaxFrameBytes      int  `json:"previewMaxFrameBytes,omitempty"`
	PreviewChunkBytes         int  `json:"previewChunkBytes,omitempty"`
	PreviewWebSocketProtocols bool `json:"previewWebSocketProtocols,omitempty"`
	LocalServerDetection      bool `json:"localServerDetection,omitempty"`
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
