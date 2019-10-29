package resource

// Service handles resource operations.
type Service struct {
	providers map[ProviderName]Provider
}

// NewService creates a new Service instance with the supplied providers.
func NewService(providers []Provider) *Service {
	pMap := make(map[ProviderName]Provider)
	for _, p := range providers {
		pMap[p.ResourceProviderName()] = p
	}
	return &Service{
		providers: pMap,
	}
}

// Provider returns a resource provider by its name.
func (s *Service) Provider(providerName ProviderName) (Provider, error) {
	if s.providers[providerName] == nil {
		return nil, ErrInvalidProvider
	}
	return s.providers[providerName], nil
}

// Resource returns a single resource using the supplied arguments.
func (s *Service) Resource(providerName ProviderName, resourceName Name, resourceID ID) (Resource, error) {
	provider, err := s.Provider(providerName)
	if err != nil {
		return nil, err
	}

	resources := provider.Resources()
	if resources[resourceName] == nil {
		return nil, ErrInvalidResourceName
	}

	return resources[resourceName].Resource(resourceID)
}
