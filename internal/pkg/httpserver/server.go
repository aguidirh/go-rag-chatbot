package httpserver

import (
	"context"
	"fmt"
	"net/http"

	"github.com/aguidirh/go-rag-chatbot/internal/pkg/app"
	"github.com/aguidirh/go-rag-chatbot/internal/pkg/config"
	"github.com/aguidirh/go-rag-chatbot/pkg/data"
	"github.com/sirupsen/logrus"
)

type HttpServer struct {
	ConfigPath string
	Log        *logrus.Logger
	config     config.Config
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

	ctx := context.Background()

	app, err := app.New(cfg, kbCfg)
	if err != nil {
		return err
	}

	http.HandleFunc("/chat", func(w http.ResponseWriter, r *http.Request) {
		var resp string

		query := r.URL.Query().Get("query") //TODO change it to a payload instead

		if query != "" {
			resp, err = app.LLMHandler.Chat(ctx, app.VectorDB.GetStore(), query)
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

		if collectionName != "" {
			err = app.VectorDB.CreateCollection(ctx, collectionName, 4096, "Cosine") //TODO ALEX get the size according to the llm used

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

		}

		fmt.Fprint(w, resp)
	})

	http.HandleFunc("/add-docs", func(w http.ResponseWriter, r *http.Request) {
		var resp string

		docs, err := app.LLMHandler.DocumentLoader(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if len(docs) > 0 {
			err = app.VectorDB.AddDocuments(ctx, docs)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		fmt.Fprint(w, resp)
	})

	fmt.Println("Server is listening on http://localhost:8080")

	addr := ":" + cfg.Spec.Server.Port
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
