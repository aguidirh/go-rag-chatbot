package data

type DocumentType string

const (
	DOC_TYPE_HTTP = DocumentType("http")
)

type DocSourceHttp struct {
	// URL is a string that holds the Uniform Resource Locator of the document.
	URL string `yaml:"url"`

	// RecursionLevels indicates whether to recursively fetch and include all linked documents within this one.
	// values greater than 0 will result in recursion.
	RecursionLevels int `yaml:"recursionLevels"`

	// FileType overrides the inferred file type. This is useful when calling a URL which doesn't have an
	// explicit file type.
	FileType string `yaml:"fileType"`

	AllowedDomains []string `yaml:"allowedDomains"`
}

type DocSourceFile struct {
	// Path is a string that holds the Uniform Resource Locator of the document.
	Path string `yaml:"path"`

	// Recurse indicates whether to recursively fetch and include all linked documents within this one.
	Recurse bool `yaml:"recurse"`

	// FileType overrides the inferred file type. This is useful when reading a file which doesn't have an
	// explicit file type.
	FileType string `yaml:"fileType"`
}

// DocSpec represents a document specification
type DocSpec struct {
	// Collection the name of the collection in to which these documents will be imported
	Collection string `yaml:"collection"`

	// Type discriminating union which defines the source type. supported options are file|http
	Type DocumentType `yaml:"type"`

	// DocSourceFiles is an array of DocSourceFile entries representing file-based sources for documents.
	DocSourceFiles []DocSourceFile `yaml:"fileSources"`

	// DocSourceHttp is an array of DocSourceHttp entries representing HTTP-based sources for documents.
	DocSourceHttp []DocSourceHttp `yaml:"httpSources"`
}

type KBConfigSpec struct {
	Docs []DocSpec `yaml:"docs"`
}

// KBConfig represents a configuration for Kubernetes resources.
type KBConfig struct {
	// ApiVersion is a string that identifies the version of the API being used.
	ApiVersion string `yaml:"apiVersion"`

	// Kind indicates the kind of resource this configuration is for.
	Kind string `yaml:"kind"`

	// Spec contains detailed information about the desired state of the Kubernetes resource.
	Spec KBConfigSpec `yaml:"spec"`
}
