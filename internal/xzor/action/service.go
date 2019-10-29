package action

import (
	"errors"

	"github.com/xzor-dev/xzor/internal/xzor/command"
)

// Service handles the execution and processing of actions.
type Service struct {
	actions          map[Hash]*Action
	commandProviders map[command.ProviderName]command.Provider
}

// NewService creates a new service instance with the provided modules.
func NewService(commandProviders []command.Provider) *Service {
	cpMap := make(map[command.ProviderName]command.Provider)
	for _, cp := range commandProviders {
		cpMap[cp.CommandProviderName()] = cp
	}
	return &Service{
		actions:          make(map[Hash]*Action),
		commandProviders: cpMap,
	}
}

// Clear removes any actions in the service's memory.
func (s *Service) Clear() {
	s.actions = make(map[Hash]*Action)
}

// ExecuteAction takes an incoming action and performs its command.
// If the action is performed without an error, it gets stored in memory for later retrieval.
func (s *Service) ExecuteAction(a *Action) (*Response, error) {
	if s.actions[a.Hash] != nil {
		return nil, ErrDuplicateAction
	}
	if s.commandProviders[a.CommandProvider] == nil {
		return nil, errors.New("invalid command provider name")
	}
	provider := s.commandProviders[a.CommandProvider]
	commands := provider.Commands()
	if commands[a.Command] == nil {
		return nil, errors.New("invalid command name")
	}
	c := commands[a.Command]
	res, err := c.Execute(a.Parameters)
	if err != nil {
		return nil, err
	}
	s.actions[a.Hash] = a
	return &Response{
		Action: a,
		Value:  res.Value,
	}, nil
}
