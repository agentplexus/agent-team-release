package output

import (
	"bytes"
	"strings"
	"testing"

	"github.com/toon-format/toon-go"

	"github.com/grokify/release-agent/pkg/actions"
	"github.com/grokify/release-agent/pkg/interactive"
)

func TestTOONWriter_WriteQuestion(t *testing.T) {
	var buf bytes.Buffer
	writer := NewTOONWriter(&buf)

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

	output := buf.String()

	// TOON format checks
	if !strings.Contains(output, "type: question") {
		t.Error("Output should contain 'type: question'")
	}
	if !strings.Contains(output, "id: test-q") {
		t.Error("Output should contain 'id: test-q'")
	}
	if !strings.Contains(output, "waiting_for: user_input") {
		t.Error("Output should contain 'waiting_for: user_input'")
	}
}

func TestTOONWriter_WriteProposal(t *testing.T) {
	var buf bytes.Buffer
	writer := NewTOONWriter(&buf)

	p := actions.Proposal{
		Description: "Update README",
		FilePath:    "README.md",
		OldContent:  "old content",
		NewContent:  "new content",
	}

	err := writer.WriteProposal(p)
	if err != nil {
		t.Fatalf("WriteProposal() error = %v", err)
	}

	output := buf.String()

	if !strings.Contains(output, "type: proposal") {
		t.Error("Output should contain 'type: proposal'")
	}
	if !strings.Contains(output, "waiting_for: user_approval") {
		t.Error("Output should contain 'waiting_for: user_approval'")
	}
}

func TestTOONWriter_WriteInfo(t *testing.T) {
	var buf bytes.Buffer
	writer := NewTOONWriter(&buf)

	err := writer.WriteInfo("Test info message")
	if err != nil {
		t.Fatalf("WriteInfo() error = %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "type: info") {
		t.Error("Output should contain 'type: info'")
	}
	if !strings.Contains(output, "Test info message") {
		t.Error("Output should contain 'Test info message'")
	}
}

func TestTOONWriter_WriteWarning(t *testing.T) {
	var buf bytes.Buffer
	writer := NewTOONWriter(&buf)

	err := writer.WriteWarning("Test warning")
	if err != nil {
		t.Fatalf("WriteWarning() error = %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "type: warning") {
		t.Error("Output should contain 'type: warning'")
	}
}

func TestTOONWriter_WriteError(t *testing.T) {
	var buf bytes.Buffer
	writer := NewTOONWriter(&buf)

	err := writer.WriteError("Test error", true)
	if err != nil {
		t.Fatalf("WriteError() error = %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "type: error") {
		t.Error("Output should contain 'type: error'")
	}
	if !strings.Contains(output, "fatal: true") {
		t.Error("Output should contain 'fatal: true'")
	}
}

func TestTOONWriter_WriteResult(t *testing.T) {
	var buf bytes.Buffer
	writer := NewTOONWriter(&buf)

	r := actions.Result{
		Name:    "test-action",
		Success: true,
		Output:  "Action completed",
	}

	err := writer.WriteResult(r)
	if err != nil {
		t.Fatalf("WriteResult() error = %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "type: result") {
		t.Error("Output should contain 'type: result'")
	}
	if !strings.Contains(output, "success: true") {
		t.Error("Output should contain 'success: true'")
	}
}

func TestTOONWriter_WriteProgress(t *testing.T) {
	var buf bytes.Buffer
	writer := NewTOONWriter(&buf)

	err := writer.WriteProgress(2, 5, "Building", "running")
	if err != nil {
		t.Fatalf("WriteProgress() error = %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "type: progress") {
		t.Error("Output should contain 'type: progress'")
	}
	if !strings.Contains(output, "step: 2") {
		t.Error("Output should contain 'step: 2'")
	}
	if !strings.Contains(output, "total_steps: 5") {
		t.Error("Output should contain 'total_steps: 5'")
	}
}

func TestDefaultTOONWriter(t *testing.T) {
	writer := DefaultTOONWriter()
	if writer == nil {
		t.Error("DefaultTOONWriter() returned nil")
	}
}

func TestTOONTokenEfficiency(t *testing.T) {
	// Compare JSON vs TOON output sizes
	q := QuestionMessage{
		Type:      string(MessageTypeQuestion),
		ID:        "test-q",
		Question:  "What is your choice?",
		InputType: "single_choice",
		Options: []OptionJSON{
			{ID: "opt1", Label: "Option 1", Description: "First option"},
			{ID: "opt2", Label: "Option 2", Description: "Second option"},
		},
		Default:    "opt1",
		Required:   true,
		WaitingFor: "user_input",
	}

	// TOON output
	toonData, err := toon.Marshal(q, toon.WithIndent(2))
	if err != nil {
		t.Fatalf("TOON Marshal error: %v", err)
	}

	// JSON output
	var jsonBuf bytes.Buffer
	jsonWriter := NewJSONWriter(&jsonBuf)
	if err := jsonWriter.Write(q); err != nil {
		t.Fatalf("JSON Write error: %v", err)
	}

	toonLen := len(toonData)
	jsonLen := jsonBuf.Len()

	t.Logf("TOON size: %d bytes", toonLen)
	t.Logf("JSON size: %d bytes", jsonLen)
	t.Logf("TOON is %.1f%% of JSON size", float64(toonLen)/float64(jsonLen)*100)

	// TOON should be smaller than JSON
	if toonLen >= jsonLen {
		t.Logf("Warning: TOON output (%d) not smaller than JSON (%d)", toonLen, jsonLen)
	}
}
