// Package output provides structured output formatting for release-agent.
package output

import (
	"encoding/json"
	"io"
	"os"

	"github.com/grokify/release-agent/pkg/actions"
	"github.com/grokify/release-agent/pkg/interactive"
)

// MessageType defines the type of JSON message.
type MessageType string

const (
	// MessageTypeQuestion is a question requiring user input.
	MessageTypeQuestion MessageType = "question"
	// MessageTypeProposal is a proposed change for review.
	MessageTypeProposal MessageType = "proposal"
	// MessageTypeInfo is an informational message.
	MessageTypeInfo MessageType = "info"
	// MessageTypeWarning is a warning message.
	MessageTypeWarning MessageType = "warning"
	// MessageTypeError is an error message.
	MessageTypeError MessageType = "error"
	// MessageTypeResult is a result from an operation.
	MessageTypeResult MessageType = "result"
	// MessageTypeProgress is a progress update.
	MessageTypeProgress MessageType = "progress"
)

// Message is the base protocol message.
type Message struct {
	Type      MessageType `json:"type" toon:"type"`
	ID        string      `json:"id,omitempty" toon:"id,omitempty"`
	Timestamp string      `json:"timestamp,omitempty" toon:"timestamp,omitempty"`
}

// QuestionMessage represents a question for user input.
type QuestionMessage struct {
	Type       string       `json:"type" toon:"type"`
	ID         string       `json:"id,omitempty" toon:"id,omitempty"`
	Question   string       `json:"question" toon:"question"`
	InputType  string       `json:"input_type" toon:"input_type"` // single_choice, multi_choice, confirm, text
	Options    []OptionJSON `json:"options,omitempty" toon:"options,omitempty"`
	Default    string       `json:"default,omitempty" toon:"default,omitempty"`
	Context    string       `json:"context,omitempty" toon:"context,omitempty"`
	Required   bool         `json:"required" toon:"required"`
	WaitingFor string       `json:"waiting_for" toon:"waiting_for"` // Always "user_input" for questions
}

// OptionJSON represents a choice option.
type OptionJSON struct {
	ID          string `json:"id" toon:"id"`
	Label       string `json:"label" toon:"label"`
	Description string `json:"description,omitempty" toon:"description,omitempty"`
}

// AnswerMessage represents a user's answer.
type AnswerMessage struct {
	QuestionID string   `json:"question_id" toon:"question_id"`
	Selected   []string `json:"selected,omitempty" toon:"selected,omitempty"`
	Text       string   `json:"text,omitempty" toon:"text,omitempty"`
	Confirmed  *bool    `json:"confirmed,omitempty" toon:"confirmed,omitempty"`
}

// ProposalMessage represents a proposed change for review.
type ProposalMessage struct {
	Type        string            `json:"type" toon:"type"`
	Description string            `json:"description" toon:"description"`
	FilePath    string            `json:"file_path,omitempty" toon:"file_path,omitempty"`
	OldContent  string            `json:"old_content,omitempty" toon:"old_content,omitempty"`
	NewContent  string            `json:"new_content,omitempty" toon:"new_content,omitempty"`
	Diff        string            `json:"diff,omitempty" toon:"diff,omitempty"`
	Metadata    map[string]string `json:"metadata,omitempty" toon:"metadata,omitempty"`
	WaitingFor  string            `json:"waiting_for" toon:"waiting_for"` // "user_approval"
	Actions     []string          `json:"actions" toon:"actions"`         // ["apply", "skip", "abort"]
}

// InfoMessage represents an informational message.
type InfoMessage struct {
	Type string `json:"type" toon:"type"`
	Text string `json:"text" toon:"text"`
}

// WarningMessage represents a warning message.
type WarningMessage struct {
	Type string `json:"type" toon:"type"`
	Text string `json:"text" toon:"text"`
}

// ErrorMessage represents an error message.
type ErrorMessage struct {
	Type  string `json:"type" toon:"type"`
	Text  string `json:"text" toon:"text"`
	Code  string `json:"code,omitempty" toon:"code,omitempty"`
	Fatal bool   `json:"fatal" toon:"fatal"`
}

// ResultMessage represents the result of an operation.
type ResultMessage struct {
	Type    string `json:"type" toon:"type"`
	Name    string `json:"name" toon:"name"`
	Success bool   `json:"success" toon:"success"`
	Output  string `json:"output,omitempty" toon:"output,omitempty"`
	Error   string `json:"error,omitempty" toon:"error,omitempty"`
	Skipped bool   `json:"skipped" toon:"skipped"`
	Reason  string `json:"reason,omitempty" toon:"reason,omitempty"`
}

// ProgressMessage represents a progress update.
type ProgressMessage struct {
	Type       string `json:"type" toon:"type"`
	Step       int    `json:"step" toon:"step"`
	TotalSteps int    `json:"total_steps" toon:"total_steps"`
	StepName   string `json:"step_name" toon:"step_name"`
	Status     string `json:"status" toon:"status"` // "running", "completed", "failed", "skipped"
}

// WorkflowResultMessage represents the final result of a workflow.
type WorkflowResultMessage struct {
	Type         string           `json:"type" toon:"type"`
	WorkflowName string           `json:"workflow_name" toon:"workflow_name"`
	Success      bool             `json:"success" toon:"success"`
	Steps        []StepResultJSON `json:"steps" toon:"steps"`
	Summary      string           `json:"summary,omitempty" toon:"summary,omitempty"`
}

// StepResultJSON represents a step result.
type StepResultJSON struct {
	Name     string `json:"name" toon:"name"`
	Status   string `json:"status" toon:"status"` // "completed", "failed", "skipped"
	Duration string `json:"duration,omitempty" toon:"duration,omitempty"`
	Output   string `json:"output,omitempty" toon:"output,omitempty"`
	Error    string `json:"error,omitempty" toon:"error,omitempty"`
}

// JSONWriter writes JSON messages to an output stream.
type JSONWriter struct {
	writer  io.Writer
	encoder *json.Encoder
}

// NewJSONWriter creates a new JSONWriter.
func NewJSONWriter(w io.Writer) *JSONWriter {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return &JSONWriter{
		writer:  w,
		encoder: encoder,
	}
}

// DefaultJSONWriter returns a JSONWriter writing to stdout.
func DefaultJSONWriter() *JSONWriter {
	return NewJSONWriter(os.Stdout)
}

// Write writes a message as JSON.
func (jw *JSONWriter) Write(msg interface{}) error {
	return jw.encoder.Encode(msg)
}

// WriteQuestion writes a question as JSON.
func (jw *JSONWriter) WriteQuestion(q interactive.Question) error {
	options := make([]OptionJSON, len(q.Options))
	for i, opt := range q.Options {
		options[i] = OptionJSON{
			ID:          opt.ID,
			Label:       opt.Label,
			Description: opt.Description,
		}
	}

	msg := QuestionMessage{
		Type:       string(MessageTypeQuestion),
		ID:         q.ID,
		Question:   q.Text,
		InputType:  q.Type.String(),
		Options:    options,
		Default:    q.Default,
		Context:    q.Context,
		Required:   true,
		WaitingFor: "user_input",
	}
	return jw.Write(msg)
}

// WriteProposal writes a proposal as JSON.
func (jw *JSONWriter) WriteProposal(p actions.Proposal) error {
	msg := ProposalMessage{
		Type:        string(MessageTypeProposal),
		Description: p.Description,
		FilePath:    p.FilePath,
		OldContent:  p.OldContent,
		NewContent:  p.NewContent,
		Metadata:    p.Metadata,
		WaitingFor:  "user_approval",
		Actions:     []string{"apply", "skip", "abort"},
	}
	return jw.Write(msg)
}

// WriteInfo writes an informational message as JSON.
func (jw *JSONWriter) WriteInfo(text string) error {
	msg := InfoMessage{
		Type: string(MessageTypeInfo),
		Text: text,
	}
	return jw.Write(msg)
}

// WriteWarning writes a warning message as JSON.
func (jw *JSONWriter) WriteWarning(text string) error {
	msg := WarningMessage{
		Type: string(MessageTypeWarning),
		Text: text,
	}
	return jw.Write(msg)
}

// WriteError writes an error message as JSON.
func (jw *JSONWriter) WriteError(text string, fatal bool) error {
	msg := ErrorMessage{
		Type:  string(MessageTypeError),
		Text:  text,
		Fatal: fatal,
	}
	return jw.Write(msg)
}

// WriteResult writes an action result as JSON.
func (jw *JSONWriter) WriteResult(r actions.Result) error {
	errStr := ""
	if r.Error != nil {
		errStr = r.Error.Error()
	}
	msg := ResultMessage{
		Type:    string(MessageTypeResult),
		Name:    r.Name,
		Success: r.Success,
		Output:  r.Output,
		Error:   errStr,
		Skipped: r.Skipped,
		Reason:  r.Reason,
	}
	return jw.Write(msg)
}

// WriteProgress writes a progress update as JSON.
func (jw *JSONWriter) WriteProgress(step, totalSteps int, stepName, status string) error {
	msg := ProgressMessage{
		Type:       string(MessageTypeProgress),
		Step:       step,
		TotalSteps: totalSteps,
		StepName:   stepName,
		Status:     status,
	}
	return jw.Write(msg)
}
