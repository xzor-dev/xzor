package action

import (
	"errors"

	"github.com/xzor-dev/xzor/internal/xzor/module"
)

// Service handles the execution and processing of actions.
type Service struct {
	Modules map[module.Name]module.Module

	actions []*Action
}

// Clear removes any actions in the service's memory.
func (s *Service) Clear() {
	s.actions = make([]*Action, 0)
}

// Execute takes an incoming action and performs its command.
// If the action is performed without an error, it gets stored in memory for later retrieval.
func (s *Service) Execute(a *Action) (*Response, error) {
	if s.actions == nil {
		s.actions = make([]*Action, 0)
	}
	if s.Modules == nil || s.Modules[a.Module] == nil {
		return nil, errors.New("invalid module provided by action")
	}
	m := s.Modules[a.Module]
	c, err := m.Command(a.Command)
	if err != nil {
		return nil, err
	}
	res, err := c.Execute(a.Arguments)
	if err != nil {
		return nil, err
	}
	s.actions = append(s.actions, a)
	return &Response{
		Action: a,
		Value:  res.Value,
	}, nil
}

// NewService creates a new service instance with the provided modules.
func NewService(modules []module.Module) *Service {
	modMap := make(map[module.Name]module.Module)
	for _, mod := range modules {
		modMap[mod.Name()] = mod
	}
	return &Service{
		Modules: modMap,
	}
}
