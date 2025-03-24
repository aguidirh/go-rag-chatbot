package data

type Config struct {
	ApiVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Spec       struct {
		VectorDB struct {
			Host       string `yaml:"host"`       // The host address for the vector database.
			Port       string `yaml:"port"`       // The port number for the vector database.
			Collection string `yaml:"collection"` // The name of the collection in the vector database.
		} `yaml:"qdrant"`

		// LLM configuration
		LLM struct {
			Model          string  `yaml:"model"`          // The model used by the language learning system.
			ScoreThreshold float32 `yaml:"scoreThreshold"` // The threshold score for generating responses.
			Temperature    float64 `yaml:"temperature"`    // The temperature parameter influencing response diversity.
		} `yaml:"llm"`

		Server struct {
			// bindAddress address on which the server will listen
			// default: 127.0.0.1
			Address string `yaml:"bindAddress"`
			// port port on which the server will listen
			// default: 8080
			Port string `yaml:"port"`
		} `yaml:"server"`
	} `yaml:"spec"`
}
