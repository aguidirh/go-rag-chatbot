package app

// func (app App) CreateCollection(ctx context.Context) error {
// 	//TODO ALEX create only if there is no collection with the name specified
// 	err := app.VectorDB.CreateCollection(ctx, app.cfg.Spec.VectorDB.Collection, 4096, "Cosine") //TODO ALEX get the size according to the llm used
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (app App) AddDocs(ctx context.Context) error {
// 	docs, err := app.LLMHandler.DocumentLoader(ctx)
// 	if err != nil {
// 		return err
// 	}

// 	if len(docs) > 0 {
// 		err = app.VectorDB.AddDocuments(ctx, docs)
// 		if err != nil {
// 			return err
// 		}
// 	}
// }

// func (app App) Chat(ctx context.Context, askMeSomething string) (string, error) {
// 	resp, err := app.LLMHandler.Chat(ctx, app.VectorDB.GetStore(), askMeSomething)

// 	return resp, err
// }
