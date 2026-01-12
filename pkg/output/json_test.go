package output

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/grokify/release-agent/pkg/actions"
	"github.com/grokify/release-agent/pkg/interactive"
)

func TestJSONWriter_WriteQuestion(t *testing.T) {
	var buf bytes.Buffer
	writer := NewJSONWriter(&buf)

	q := interactive.Question{
		ID:   "test-q",
		Text: "What is your choice?",
		Type: interactive.QuestionTypeSingleChoice,
		Options: []interactive.Option{
			{ID: "opt1", Label: "Option 1", Description: "First option"},
			{ID: "opt2", Label: "Option 2", Description: "Second option"},
		},
		Default: "opt1",
		Context: "Some context here",
	}

	err := writer.WriteQuestion(q)
	if err != nil {
		t.Fatalf("WriteQuestion() error = %v", err)
	}

	// Parse the output
	var msg QuestionMessage
	if err := json.Unmarshal(buf.Bytes(), &msg); err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	if msg.Type != string(MessageTypeQuestion) {
		t.Errorf("Type = %s, want %s", msg.Type, MessageTypeQuestion)
	}
	if msg.ID != "test-q" {
		t.Errorf("ID = %s, want test-q", msg.ID)
	}
	if msg.Question != "What is your choice?" {
		t.Errorf("Question = %s, want 'What is your choice?'", msg.Question)
	}
	if msg.InputType != "single_choice" {
		t.Errorf("InputType = %s, want single_choice", msg.InputType)
	}
	if len(msg.Options) != 2 {
		t.Errorf("Options length = %d, want 2", len(msg.Options))
	}
	if msg.Default != "opt1" {
		t.Errorf("Default = %s, want opt1", msg.Default)
	}
	if msg.WaitingFor != "user_input" {
		t.Errorf("WaitingFor = %s, want user_input", msg.WaitingFor)
	}
}

func TestJSONWriter_WriteProposal(t *testing.T) {
	var buf bytes.Buffer
	writer := NewJSONWriter(&buf)

	p := actions.Proposal{
		Description: "Update README",
		FilePath:    "README.md",
		OldContent:  "old content",
		NewContent:  "new content",
		Metadata:    map[string]string{"key": "value"},
	}

	err := writer.WriteProposal(p)
	if err != nil {
		t.Fatalf("WriteProposal() error = %v", err)
	}

	var msg ProposalMessage
	if err := json.Unmarshal(buf.Bytes(), &msg); err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	if msg.Type != string(MessageTypeProposal) {
		t.Errorf("Type = %s, want %s", msg.Type, MessageTypeProposal)
	}
	if msg.Description != "Update README" {
		t.Errorf("Description = %s, want 'Update README'", msg.Description)
	}
	if msg.FilePath != "README.md" {
		t.Errorf("FilePath = %s, want README.md", msg.FilePath)
	}
	if msg.WaitingFor != "user_approval" {
		t.Errorf("WaitingFor = %s, want user_approval", msg.WaitingFor)
	}
	if len(msg.Actions) != 3 {
		t.Errorf("Actions length = %d, want 3", len(msg.Actions))
	}
}

func TestJSONWriter_WriteInfo(t *testing.T) {
	var buf bytes.Buffer
	writer := NewJSONWriter(&buf)

	err := writer.WriteInfo("Test info message")
	if err != nil {
		t.Fatalf("WriteInfo() error = %v", err)
	}

	var msg InfoMessage
	if err := json.Unmarshal(buf.Bytes(), &msg); err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	if msg.Type != string(MessageTypeInfo) {
		t.Errorf("Type = %s, want %s", msg.Type, MessageTypeInfo)
	}
	if msg.Text != "Test info message" {
		t.Errorf("Text = %s, want 'Test info message'", msg.Text)
	}
}

func TestJSONWriter_WriteWarning(t *testing.T) {
	var buf bytes.Buffer
	writer := NewJSONWriter(&buf)

	err := writer.WriteWarning("Test warning")
	if err != nil {
		t.Fatalf("WriteWarning() error = %v", err)
	}

	var msg WarningMessage
	if err := json.Unmarshal(buf.Bytes(), &msg); err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	if msg.Type != string(MessageTypeWarning) {
		t.Errorf("Type = %s, want %s", msg.Type, MessageTypeWarning)
	}
}

func TestJSONWriter_WriteError(t *testing.T) {
	var buf bytes.Buffer
	writer := NewJSONWriter(&buf)

	err := writer.WriteError("Test error", true)
	if err != nil {
		t.Fatalf("WriteError() error = %v", err)
	}

	var msg ErrorMessage
	if err := json.Unmarshal(buf.Bytes(), &msg); err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	if msg.Type != string(MessageTypeError) {
		t.Errorf("Type = %s, want %s", msg.Type, MessageTypeError)
	}
	if msg.Text != "Test error" {
		t.Errorf("Text = %s, want 'Test error'", msg.Text)
	}
	if !msg.Fatal {
		t.Error("Fatal = false, want true")
	}
}

func TestJSONWriter_WriteResult(t *testing.T) {
	var buf bytes.Buffer
	writer := NewJSONWriter(&buf)

	r := actions.Result{
		Name:    "test-action",
		Success: true,
		Output:  "Action completed",
		Skipped: false,
	}

	err := writer.WriteResult(r)
	if err != nil {
		t.Fatalf("WriteResult() error = %v", err)
	}

	var msg ResultMessage
	if err := json.Unmarshal(buf.Bytes(), &msg); err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	if msg.Type != string(MessageTypeResult) {
		t.Errorf("Type = %s, want %s", msg.Type, MessageTypeResult)
	}
	if msg.Name != "test-action" {
		t.Errorf("Name = %s, want test-action", msg.Name)
	}
	if !msg.Success {
		t.Error("Success = false, want true")
	}
}

func TestJSONWriter_WriteProgress(t *testing.T) {
	var buf bytes.Buffer
	writer := NewJSONWriter(&buf)

	err := writer.WriteProgress(2, 5, "Building", "running")
	if err != nil {
		t.Fatalf("WriteProgress() error = %v", err)
	}

	var msg ProgressMessage
	if err := json.Unmarshal(buf.Bytes(), &msg); err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	if msg.Type != string(MessageTypeProgress) {
		t.Errorf("Type = %s, want %s", msg.Type, MessageTypeProgress)
	}
	if msg.Step != 2 {
		t.Errorf("Step = %d, want 2", msg.Step)
	}
	if msg.TotalSteps != 5 {
		t.Errorf("TotalSteps = %d, want 5", msg.TotalSteps)
	}
	if msg.StepName != "Building" {
		t.Errorf("StepName = %s, want Building", msg.StepName)
	}
	if msg.Status != "running" {
		t.Errorf("Status = %s, want running", msg.Status)
	}
}

func TestMessageTypes(t *testing.T) {
	tests := []struct {
		mt   MessageType
		want string
	}{
		{MessageTypeQuestion, "question"},
		{MessageTypeProposal, "proposal"},
		{MessageTypeInfo, "info"},
		{MessageTypeWarning, "warning"},
		{MessageTypeError, "error"},
		{MessageTypeResult, "result"},
		{MessageTypeProgress, "progress"},
	}

	for _, tt := range tests {
		if string(tt.mt) != tt.want {
			t.Errorf("MessageType = %s, want %s", tt.mt, tt.want)
		}
	}
}

func TestDefaultJSONWriter(t *testing.T) {
	writer := DefaultJSONWriter()
	if writer == nil {
		t.Error("DefaultJSONWriter() returned nil")
	}
}

func TestAnswerMessage(t *testing.T) {
	confirmed := true
	ans := AnswerMessage{
		QuestionID: "q1",
		Selected:   []string{"opt1", "opt2"},
		Text:       "custom text",
		Confirmed:  &confirmed,
	}

	data, err := json.Marshal(ans)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	var parsed AnswerMessage
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if parsed.QuestionID != "q1" {
		t.Errorf("QuestionID = %s, want q1", parsed.QuestionID)
	}
	if len(parsed.Selected) != 2 {
		t.Errorf("Selected length = %d, want 2", len(parsed.Selected))
	}
	if parsed.Text != "custom text" {
		t.Errorf("Text = %s, want 'custom text'", parsed.Text)
	}
	if parsed.Confirmed == nil || !*parsed.Confirmed {
		t.Error("Confirmed should be true")
	}
}

func TestWorkflowResultMessage(t *testing.T) {
	msg := WorkflowResultMessage{
		Type:         "workflow_result",
		WorkflowName: "Release",
		Success:      true,
		Steps: []StepResultJSON{
			{Name: "Build", Status: "completed", Duration: "1s"},
			{Name: "Test", Status: "completed", Duration: "2s"},
		},
		Summary: "All steps passed",
	}

	data, err := json.Marshal(msg)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	if !strings.Contains(string(data), "workflow_result") {
		t.Error("JSON should contain 'workflow_result'")
	}
	if !strings.Contains(string(data), "Release") {
		t.Error("JSON should contain 'Release'")
	}
}
