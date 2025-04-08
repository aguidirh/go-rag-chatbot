package util

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/aguidirh/go-rag-chatbot/internal/pkg/adapters"
	"github.com/aguidirh/go-rag-chatbot/internal/pkg/frameworks/databases/qdrant"
	"github.com/aguidirh/go-rag-chatbot/pkg/data"
	"github.com/gocolly/colly/v2"
	"github.com/sirupsen/logrus"
	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/schema"
)

type KBLoader struct {
	kbConfig   *data.KBConfig
	config     *data.Config
	log        *logrus.Logger
	embed      embeddings.Embedder
	LLMHandler adapters.LLMHandler
	mtx        sync.Mutex
	ctx        context.Context
}

// Load loads the knowledge base from the configuration. This process is handled in another thread to avoid blocking the main thread.
// only a single thread will be used for this. if there is a load in progress, Load will return an error.
func (k *KBLoader) Load(collections ...string) error {
	if !k.mtx.TryLock() {
		return fmt.Errorf("Vector database is being initialized. Impacted collections may be momentarily unavailable. Please try again later. Thank you for your patience!")
	}
	go func() {
		defer k.mtx.Unlock()
	}()

	targetDocCollections := make([]data.DocSpec, 0)
	collectionNames := make([]string, 0)
	for _, c := range k.kbConfig.Spec.Docs {
		if len(collections) > 0 {
			for _, col := range collections {
				if c.Collection == col {
					targetDocCollections = append(targetDocCollections, c)
					collectionNames = append(collectionNames, col)
					break
				}
			}
		} else {
			targetDocCollections = append(targetDocCollections, c)
			collectionNames = append(collectionNames, c.Collection)
		}
	}

	k.log.Infof("loading knowledge base with collections: %s", strings.Join(collectionNames, ","))

	for _, doc := range targetDocCollections {
		k.log.Infof("creating collection %s", doc.Collection)
		vectorDB, err := qdrant.New(k.config.Spec.VectorDB.Host, k.config.Spec.VectorDB.Port, doc.Collection, k.embed)
		if err != nil {
			k.log.Errorf("failed to create vector database for collection %s: %v", doc.Collection, err)
			return err
		}
		err = vectorDB.CreateCollection(k.ctx, doc.Collection, k.config.Spec.VectorDB.VectorSize, k.config.Spec.VectorDB.Distance)
		if err != nil {
			k.log.Errorf("failed to create collection %s: %v", doc.Collection, err)
			continue
		}
		err = k.LLMHandler.LoadDocumentsFromConfig(k.ctx, func(docs []schema.Document, e *colly.HTMLElement) error {
			for _, d := range docs {
				d.Metadata["id"] = e.Attr("id")
				d.Metadata["path"] = e.Request.URL.Path
				d.Metadata["url"] = fmt.Sprintf("%s#%s", e.Request.URL.String(), d.Metadata["id"])
				for k, v := range doc.Metadata {
					d.Metadata[k] = v
				}
				err = vectorDB.AddDocuments(k.ctx, []schema.Document{d})
				if err != nil {
					k.log.Errorf("failed to add documents to collection %s: %v", doc.Collection, err)
					continue
				}
			}
			return nil
		}, doc.Collection)
		if err != nil {
			k.log.Errorf("failed to load documents for collection %s: %v", doc.Collection, err)
			continue
		}

		k.log.Infof("documents loaded for collection %s", doc.Collection)
	}
	return nil
}

func NewKBLoader(ctx context.Context, kbConfig *data.KBConfig, config *data.Config, LLMHandler adapters.LLMHandler, embedder embeddings.Embedder, log *logrus.Logger) *KBLoader {
	return &KBLoader{
		kbConfig:   kbConfig,
		log:        log,
		mtx:        sync.Mutex{},
		config:     config,
		embed:      embedder,
		ctx:        ctx,
		LLMHandler: LLMHandler,
	}
}
