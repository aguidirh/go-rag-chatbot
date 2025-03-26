package httpserver

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/aguidirh/go-rag-chatbot/internal/pkg/app"
	"github.com/aguidirh/go-rag-chatbot/internal/pkg/config"
	"github.com/aguidirh/go-rag-chatbot/internal/pkg/frameworks/databases/qdrant"
	"github.com/aguidirh/go-rag-chatbot/pkg/data"
	"github.com/sirupsen/logrus"
	"github.com/tmc/langchaingo/embeddings"
)

type HttpServer struct {
	ConfigPath     string
	Log            *logrus.Logger
	config         config.Config
	VectorDBHost   string
	VectorDBPort   int
	BindAddress    string
	BindPort       int
	ModelServerURL string
	Embedder       embeddings.Embedder
}

func (h *HttpServer) Run() error {
	h.config = config.Config{
		ConfigPath: h.ConfigPath,
		Log:        h.Log,
	}

	cfg, kbCfg, err := h.loadConfigs()
	if err != nil {
		return err
	}
	bindAddress := "127.0.0.1"

	ctx := context.Background()

	if len(h.VectorDBHost) > 0 {
		h.Log.Infof("vector DB host overridden by --vectordb-host to %s", h.VectorDBHost)
		cfg.Spec.VectorDB.Host = h.VectorDBHost
	}
	if h.VectorDBPort > 0 {
		h.Log.Infof("vector DB port overridden by --vectordb-port to %d", h.VectorDBPort)
		cfg.Spec.VectorDB.Port = strconv.Itoa(h.VectorDBPort)
	}

	if len(h.BindAddress) > 0 {
		h.Log.Infof("bind host overridden by --bind-address to %s", h.BindAddress)
		bindAddress = h.BindAddress
	}
	if h.BindPort > 0 {
		h.Log.Infof("bind port overridden by --bind-port to %d", h.BindPort)
		cfg.Spec.Server.Port = strconv.Itoa(h.BindPort)
	}
	if len(h.ModelServerURL) > 0 {
		h.Log.Infof("Model server url overridden by --model-server-url to %s", h.ModelServerURL)
		cfg.Spec.LLM.URL = h.ModelServerURL
	}

	app, err := app.New(ctx, cfg, kbCfg, h.Log)
	if err != nil {
		return err
	}

	http.HandleFunc("/chat", func(w http.ResponseWriter, r *http.Request) {
		var resp string
		// to-do, add a collection query parameter to the URL and use it to create a collection in qdrant if it doesn't exist.
		vectorDB, err := qdrant.New(cfg.Spec.VectorDB.Host, cfg.Spec.VectorDB.Port, "", app.Embedder)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		query := r.URL.Query().Get("query") //TODO change it to a payload instead

		if query != "" {
			resp, err = app.LLMHandler.Chat(ctx, vectorDB.GetStore(), query)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		fmt.Fprint(w, resp)
	})

	http.HandleFunc("/create-collection", func(w http.ResponseWriter, r *http.Request) {
		var resp string

		collectionName := r.URL.Query().Get("collection-name")
		vectorDB, err := qdrant.New(cfg.Spec.VectorDB.Host, cfg.Spec.VectorDB.Port, collectionName, app.Embedder)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if collectionName != "" {
			err = vectorDB.CreateCollection(ctx, collectionName, 4096, "Cosine") //TODO ALEX get the size according to the llm used

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

		}

		fmt.Fprint(w, resp)
	})

	http.HandleFunc("/add-docs", func(w http.ResponseWriter, r *http.Request) {
		var resp string

		// docs, err := app.LLMHandler.DocumentLoader(ctx)
		// if err != nil {
		// 	http.Error(w, err.Error(), http.StatusInternalServerError)
		// 	return
		// }
		// vectorDB, err := qdrant.New(cfg.Spec.VectorDB.Host, cfg.Spec.VectorDB.Port, "", app.Embedder)
		// if err != nil {
		// 	http.Error(w, err.Error(), http.StatusInternalServerError)
		// 	return
		// }

		// if len(docs) > 0 {
		// 	err = vectorDB.AddDocuments(ctx, docs)
		// 	if err != nil {
		// 		http.Error(w, err.Error(), http.StatusInternalServerError)
		// 		return
		// 	}
		// }

		fmt.Fprint(w, resp)
	})

	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "ok")
	})

	fmt.Printf("server is listening on http://%s:%s\n", bindAddress, cfg.Spec.Server.Port)

	addr := bindAddress + ":" + cfg.Spec.Server.Port
	return http.ListenAndServe(addr, nil)
}

// TODO create a generic function to load configs
func (h *HttpServer) loadConfigs() (data.Config, data.KBConfig, error) {
	var cfg data.Config
	var kbCfg data.KBConfig

	cfg, err := h.config.LoadConfig()
	if err != nil {
		return cfg, kbCfg, err
	}
	kbConfig := config.KbConfig{
		ConfigPath: h.ConfigPath,
		Log:        h.Log,
	}
	kbCfg, err = kbConfig.LoadKBConfig()
	if err != nil {
		return cfg, kbCfg, err
	}

	return cfg, kbCfg, nil
}
