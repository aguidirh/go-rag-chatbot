package data

type VectorDBType string
type LLMProviderType string

const (
	QdrantVectorDBType    VectorDBType    = "qdrant"
	LLMProviderTypeOllama LLMProviderType = "ollama"
	LLMProviderTypeOpenAI LLMProviderType = "openai"
)

type VectorDBQdrant struct {
}

type VectorDB struct {
	Host           string          `yaml:"host"`
	Port           string          `yaml:"port"`
	Collection     string          `yaml:"collection"`
	Distance       string          `yaml:"distance"`
	VectorSize     int             `yaml:"vectorSize"`
	Type           VectorDBType    `yaml:"type"`
	VectorDBQdrant *VectorDBQdrant `yaml:"qdrant"`
}

type Config struct {
	ApiVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Spec       struct {
		VectorDB VectorDB `yaml:"vectorDB"`

		// LLM configuration
		LLM struct {
			ChatModel      string          `yaml:"chatModel"`      // The model used by the language learning system.
			EmbeddingModel string          `yaml:"embeddingModel"` // The embedding model used by the language learning system.
			ScoreThreshold float32         `yaml:"scoreThreshold"` // The threshold score for generating responses.
			Temperature    float64         `yaml:"temperature"`    // The temperature parameter influencing response diversity.
			URL            string          `yaml:"url"`            // The URL of the model server
			ProviderType   LLMProviderType `yaml:"providerType"`   // The type of language learning system provider.

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
