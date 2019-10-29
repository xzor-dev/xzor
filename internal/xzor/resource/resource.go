package resource

// Getter is used to get resources by their IDs.
type Getter interface {
	Resource(ID) (Resource, error)
}

// GetterMap maps resource getters to resource names.
type GetterMap map[Name]Getter

// ID is a unique string used to identify resources.
type ID string

// Name is a string used to identify resource types.
type Name string

// Provider is used to return a set of resource getters.
type Provider interface {
	Resources() GetterMap
	ResourceProviderName() ProviderName
}

// ProviderName is a string used to identify groups of resources.
type ProviderName string

// Resource identifies an individual piece of data.
type Resource interface {
	ResourceID() ID
	ResourceName() Name
}
