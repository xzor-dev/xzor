package resource

// Getter is used to get resources by their IDs.
type Getter interface {
	Resource(ID) Resource
}

// ID is a unique string used to identify resources.
type ID string

// Name is a string used to identify resource types.
type Name string

// Provider is used to return a set of resource getters.
type Provider interface {
	Resources() map[Name]Getter
}

// ProviderName is a string used to identify groups of resources.
type ProviderName string

// Resource identifies an individual piece of data.
type Resource interface {
	ResourceID() ID
	ResourceName() Name
}
