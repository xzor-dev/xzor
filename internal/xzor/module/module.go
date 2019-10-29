package module

// Module defines a package of resources and commands.
type Module interface {
	Name() Name
}

// Name is a string name of a module.
type Name string
