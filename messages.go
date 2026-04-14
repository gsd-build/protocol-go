// Package protocol defines the wire format between the GSD Cloud daemon,
// the Fly.io relay, and the browser. See PROTOCOL.md for the authoritative
// specification; every change here must be mirrored in that file.
package protocol

import "encoding/json"

// Message type constants.
const (
	MsgTypeTask               = "task"
	MsgTypeStop               = "stop"
	MsgTypePermissionResponse = "permissionResponse"
	MsgTypeQuestionResponse   = "questionResponse"
	MsgTypeBrowseDir          = "browseDir"
	MsgTypeReadFile           = "readFile"
	MsgTypeMkDir              = "mkDir"
	MsgTypeMkDirResult        = "mkDirResult"

	MsgTypeStream            = "stream"
	MsgTypeTaskStarted       = "taskStarted"
	MsgTypeTaskComplete      = "taskComplete"
	MsgTypeTaskError         = "taskError"
	MsgTypeTaskCancelled     = "taskCancelled"
	MsgTypePermissionRequest = "permissionRequest"
	MsgTypeQuestion          = "question"
	MsgTypeHeartbeat         = "heartbeat"
	MsgTypeBrowseDirResult   = "browseDirResult"
	MsgTypeReadFileResult    = "readFileResult"

	MsgTypeHello   = "hello"
	MsgTypeWelcome = "welcome"

	MsgTypeMachineStatus   = "machineStatus"
	MsgTypeUpdateAvailable = "updateAvailable"
)

// Task is sent from the browser to the daemon to dispatch a user message.
type Task struct {
	Type                string   `json:"type"`
	TaskID              string   `json:"taskId"`
	SessionID           string   `json:"sessionId"`
	ChannelID           string   `json:"channelId"`
	Prompt              string   `json:"prompt"`
	Model               string   `json:"model"`
	Effort              string   `json:"effort"`
	PermissionMode      string   `json:"permissionMode"`
	PersonaSystemPrompt string   `json:"personaSystemPrompt,omitempty"`
	CWD                 string   `json:"cwd"`
	ClaudeSessionID     string   `json:"claudeSessionId,omitempty"` // passed to --resume
	RequestID           string   `json:"requestId,omitempty"`
	Traceparent         string   `json:"traceparent,omitempty"` // W3C trace context
	ImageURLs           []string `json:"imageUrls,omitempty"`   // user-attached image URLs
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
	ResultSummary   string `json:"resultSummary,omitempty"`
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
	Type      string        `json:"type"`
	RequestID string        `json:"requestId"`
	ChannelID string        `json:"channelId"`
	OK        bool          `json:"ok"`
	Entries   []BrowseEntry `json:"entries,omitempty"`
	Error     string        `json:"error,omitempty"`
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

// Hello is the first frame sent by the daemon after connecting.
type Hello struct {
	Type          string   `json:"type"`
	MachineID     string   `json:"machineId"`
	DaemonVersion string   `json:"daemonVersion"`
	OS            string   `json:"os"`
	Arch          string   `json:"arch"`
	ActiveTasks   []string `json:"activeTasks,omitempty"`
}

// Welcome is the relay's response to Hello.
type Welcome struct {
	Type                string `json:"type"`
	LatestDaemonVersion string `json:"latestDaemonVersion,omitempty"`
}

// MachineStatus is pushed to all connected browsers when a daemon connects or disconnects.
type MachineStatus struct {
	Type      string `json:"type"`
	MachineID string `json:"machineId"`
	Online    bool   `json:"online"`
}

// UpdateAvailable is sent by the daemon to the relay (which forwards to browsers)
// when the daemon detects a newer version is available via the Welcome message.
type UpdateAvailable struct {
	Type           string `json:"type"`
	CurrentVersion string `json:"currentVersion"`
	LatestVersion  string `json:"latestVersion"`
}
