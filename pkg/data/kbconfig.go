package data

// KBConfig represents a configuration for Kubernetes resources.
type KBConfig struct {
	// ApiVersion is a string that identifies the version of the API being used.
	ApiVersion string `yaml:"apiVersion"`

	// Kind indicates the kind of resource this configuration is for.
	Kind string `yaml:"kind"`

	// Spec contains detailed information about the desired state of the Kubernetes resource.
	Spec struct {
		// Docs is a slice of strings that provides documentation or additional information related to the resource.
		Docs []string `yaml:"docs"`
	} `yaml:"spec"`
}
