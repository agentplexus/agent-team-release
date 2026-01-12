package output

import (
	"io"
	"os"

	"github.com/toon-format/toon-go"

	"github.com/grokify/release-agent/pkg/actions"
	"github.com/grokify/release-agent/pkg/interactive"
)

// TOONWriter writes TOON-formatted messages to an output stream.
type TOONWriter struct {
	writer  io.Writer
	encoder *toon.Encoder
}

// NewTOONWriter creates a new TOONWriter.
func NewTOONWriter(w io.Writer) *TOONWriter {
	return &TOONWriter{
		writer:  w,
		encoder: toon.NewEncoder(toon.WithIndent(2)),
	}
}

// DefaultTOONWriter returns a TOONWriter writing to stdout.
func DefaultTOONWriter() *TOONWriter {
	return NewTOONWriter(os.Stdout)
}

// Write writes a message as TOON.
func (tw *TOONWriter) Write(msg interface{}) error {
	data, err := toon.Marshal(msg, toon.WithIndent(2))
	if err != nil {
		return err
	}
	_, err = tw.writer.Write(data)
	if err != nil {
		return err
	}
	// Add newline separator between messages
	_, err = tw.writer.Write([]byte("\n"))
	return err
}

// WriteQuestion writes a question as TOON.
func (tw *TOONWriter) WriteQuestion(q interactive.Question) error {
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
	return tw.Write(msg)
}

// WriteProposal writes a proposal as TOON.
func (tw *TOONWriter) WriteProposal(p actions.Proposal) error {
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
	return tw.Write(msg)
}

// WriteInfo writes an informational message as TOON.
func (tw *TOONWriter) WriteInfo(text string) error {
	msg := InfoMessage{
		Type: string(MessageTypeInfo),
		Text: text,
	}
	return tw.Write(msg)
}

// WriteWarning writes a warning message as TOON.
func (tw *TOONWriter) WriteWarning(text string) error {
	msg := WarningMessage{
		Type: string(MessageTypeWarning),
		Text: text,
	}
	return tw.Write(msg)
}

// WriteError writes an error message as TOON.
func (tw *TOONWriter) WriteError(text string, fatal bool) error {
	msg := ErrorMessage{
		Type:  string(MessageTypeError),
		Text:  text,
		Fatal: fatal,
	}
	return tw.Write(msg)
}

// WriteResult writes an action result as TOON.
func (tw *TOONWriter) WriteResult(r actions.Result) error {
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
	return tw.Write(msg)
}

// WriteProgress writes a progress update as TOON.
func (tw *TOONWriter) WriteProgress(step, totalSteps int, stepName, status string) error {
	msg := ProgressMessage{
		Type:       string(MessageTypeProgress),
		Step:       step,
		TotalSteps: totalSteps,
		StepName:   stepName,
		Status:     status,
	}
	return tw.Write(msg)
}
