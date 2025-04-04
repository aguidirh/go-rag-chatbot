package httpserver

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/aguidirh/go-rag-chatbot/internal/pkg/adapters"
	"github.com/aguidirh/go-rag-chatbot/internal/pkg/app"
	"github.com/aguidirh/go-rag-chatbot/internal/pkg/config"
	"github.com/aguidirh/go-rag-chatbot/internal/pkg/frameworks/util"
	"github.com/aguidirh/go-rag-chatbot/pkg/data"
	"github.com/gocolly/colly/v2"
	"github.com/sirupsen/logrus"
	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/schema"
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
	SkipKbLoad     bool
	vectorDB       adapters.VectorDB
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

	app, err := app.New(ctx, cfg, kbCfg, h.SkipKbLoad, h.Log)
	if err != nil {
		return err
	}

	embeddingLlmHandler := app.LLMHandler

	http.HandleFunc("/chat", func(w http.ResponseWriter, r *http.Request) {
		var resp, query, collectionName string
		var vectorDB adapters.VectorDB

		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			resp = fmt.Sprintf("unable to parse request body. %v", err)
			h.Log.Error(resp)
			goto chat_response
		}
		query = string(bodyBytes)

		collectionName, err = util.RequiredParameterAsString(r, "collection-name")
		if err != nil {
			resp = fmt.Sprintf("unable to get collection name. %v", err)
			h.Log.Error(resp)
			goto chat_response
		}

		vectorDB, err = util.GetVectorDBForCollection(collectionName, &cfg.Spec.VectorDB, app.Embedder)
		if err != nil {
			resp = fmt.Sprintf("unable to get vector db. %v", err)
			h.Log.Error(resp)
			goto chat_response
		}

		if query != "" {
			resp, err = app.LLMHandler.Chat(ctx, vectorDB.GetStore(), query)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	chat_response:
		fmt.Fprint(w, resp)
	})

	http.HandleFunc("/collections/initialize", func(w http.ResponseWriter, r *http.Request) {
		var resp string
		collectionName := r.URL.Query().Get("collection-name")
		if len(collectionName) == 0 {
			h.Log.Infof("Initializing database for all collections")
		} else {
			h.Log.Infof("Initializing database for collection: %s", collectionName)
		}

		go func() {
			if collectionName == "" { // Initialize all collections
				err := app.KBLoader.Load()
				if err != nil {
					h.Log.Errorf("failed to load KB: %v", err)
				}
			}
			err := app.KBLoader.Load(collectionName)
			if err != nil {
				h.Log.Errorf("failed to load KB: %v", err)
			}
		}()
		resp = "Vector database is being initialized. Impacted collections may be momentarily unavailable. Please try again later. Thank you for your patience!"

		fmt.Fprint(w, resp)
	})

	http.HandleFunc("/collection", func(w http.ResponseWriter, r *http.Request) {
		var resp string
		vectorSize := util.GetQueryParameterAsInt(r, "vector-size", cfg.Spec.VectorDB.VectorSize)
		distance := util.GetQueryParameterAsString(r, "distance", cfg.Spec.VectorDB.Distance)
		collectionName := r.URL.Query().Get("collection-name")
		if len(collectionName) == 0 {
			resp = "Please provide a collection name."

		} else {
			vectorDB, err := util.GetVectorDBForCollection(collectionName, &cfg.Spec.VectorDB, app.Embedder)
			if err != nil {
				resp = fmt.Sprintf("unable to get vector db. %v", err)
				h.Log.Error(resp)
				goto collection_result
			}

			exists, err := vectorDB.DoesCollectionExist(ctx, collectionName)
			if err != nil {
				h.Log.Errorf("failed to check if collection exists: %v", err)
				resp = "Failed to check if collection exists"
				goto collection_result
			}

			if r.Method == "GET" {
				resp = fmt.Sprintf("Collection %s exists: %t", collectionName, exists)
			} else if r.Method == "POST" && !exists {
				err := vectorDB.CreateCollection(ctx, collectionName, vectorSize, distance)
				if err != nil {
					h.Log.Errorf("failed to create collection: %v", err)
					resp = "Failed to create collection"
					goto collection_result
				} else {
					resp = fmt.Sprintf("Collection %s created successfully", collectionName)
					goto collection_result
				}
			} else if r.Method == "POST" && exists {
				resp = fmt.Sprintf("Collection %s already exists", collectionName)
			} else if r.Method == "DELETE" {
				if exists {
					err := vectorDB.DeleteCollection(ctx, collectionName)
					if err != nil {
						h.Log.Errorf("failed to delete collection: %v", err)
						resp = "Failed to delete collection"
						goto collection_result
					}
					resp = fmt.Sprintf("Collection %s deleted successfully", collectionName)
				} else {
					resp = fmt.Sprintf("Collection %s does not exist", collectionName)
				}
			} else {
				resp = "Invalid request method"
			}
		}

	collection_result:
		fmt.Fprint(w, resp)
	})

	http.HandleFunc("/docs/add", func(w http.ResponseWriter, r *http.Request) {
		var resp, collectionName string
		var vectorDB adapters.VectorDB
		var err error

		if r.Method == "POST" {
			collectionName, err = util.RequiredParameterAsString(r, "collection-name")
			if err != nil {
				resp = fmt.Sprintf("Missing required parameter: collection-name. %v", err)
				goto docs_result
			}

			vectorDB, err = util.GetVectorDBForCollection(collectionName, &cfg.Spec.VectorDB, app.Embedder)
			if err != nil {
				resp = fmt.Sprintf("unable to get vector db. %v", err)
				h.Log.Error(resp)
				goto docs_result
			}
			err = embeddingLlmHandler.LoadDocumentsFromHttpRequest(ctx, func(docs []schema.Document, e *colly.HTMLElement) error {
				err = vectorDB.AddDocuments(ctx, docs)
				if err != nil {
					resp = fmt.Sprintf("unable to add documents. %v", err)
					h.Log.Error(resp)
					return err
				}
				return nil
			}, collectionName, r)
			if err != nil {
				resp = fmt.Sprintf("unable to add documents. %v", err)
				h.Log.Error(resp)
				goto docs_result
			}

			resp = "Documents added successfully"
		} else if r.Method == "GET" {
			resp = "Documentation endpoint"
		} else {
			resp = "Invalid request method"
		}
	docs_result:
		fmt.Fprint(w, resp)
	})

	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "ok")
	})

	fmt.Printf("server is listening on http://%s:%s\n", bindAddress, cfg.Spec.Server.Port)

	addr := bindAddress + ":" + cfg.Spec.Server.Port
	return http.ListenAndServe(addr, nil)
}

// LoadConfigs loads the configuration files and returns the configuration objects with
// reasonable defaults.
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
