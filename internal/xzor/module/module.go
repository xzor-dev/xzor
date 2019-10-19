package module

import "github.com/xzor-dev/xzor/internal/xzor/command"

// Module defines a package of resources and commands.
type Module interface {
	Command(command.Name) (command.Command, error)
	Name() Name
	Resources() map[ResourceName]ResourceGetter
}

// Name is a string name of a module.
type Name string

// Resource defines a single resource belonging to a module.
type Resource interface {
	ResourceID() ResourceID
	ResourceName() ResourceName
}

// ResourceGetter handles the retrieval of individual resources.
type ResourceGetter interface {
	Resource(ResourceID) (Resource, error)
}

// ResourceID is a string used to identify single module resources.
type ResourceID string

// ResourceName is used to categorize resources.
type ResourceName string
