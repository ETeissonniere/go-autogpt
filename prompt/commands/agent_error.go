package commands

import "fmt"

// A special error type that is used to signal that the error is not a bug in the
// code, but rather an error to be reported to the agent.
type AgentError struct {
	Err error
}

func (e *AgentError) Error() string {
	return e.Err.Error()
}

func (e *AgentError) AgentExplainer() string {
	return fmt.Sprintf("an error happened: %v", e.Err)
}

// NewAgentError creates a new AgentError.
func NewAgentError(err error) *AgentError {
	return &AgentError{Err: err}
}
