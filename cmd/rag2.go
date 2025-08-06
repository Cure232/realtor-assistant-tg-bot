package main

import (
	"context"
	"log"
	"time"

	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/textsplitter"
	"github.com/tmc/langchaingo/vectorstores/pgvector"
)

var (
	ollamaHost     = "http://localhost:11434"
	ollamaEmbModel = "mxbai-embed-large"
	ollamaGenModel = "llama3:8b"
	databaseURL    = "postgres://postgres:okocha232@localhost:5432/rag_test?sslmode=disable"
)

type ragLangchainTest struct {
	embedder       *embeddings.EmbedderImpl
	generation_llm *ollama.LLM
	vectorStore    *pgvector.Store
}

func (ai *ragLangchainTest) Init() *ragLangchainTest {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// Инициализация генератора Ollama
	generation_llm, err := ollama.New(
		ollama.WithModel(ollamaGenModel),
		ollama.WithServerURL(ollamaHost),
	)
	if err != nil {
		log.Fatalf("Ошибка инициализации генератора: %v", err)
	}
	ai.generation_llm = generation_llm
	log.Printf("Генератор инициализирован: %s", ollamaGenModel)

	// Инициализация эмбеддера Ollama
	embedder_llm, err := ollama.New(
		ollama.WithServerURL(ollamaHost),
		ollama.WithModel(ollamaEmbModel),
	)
	if err != nil {
		log.Fatalf("Ошибка инициализации нейросети эмбеддера: %v", err)
	}
	embedder, err := embeddings.NewEmbedder(
		embedder_llm,
		embeddings.WithStripNewLines(true),
		embeddings.WithBatchSize(32),
	)
	if err != nil {
		log.Fatalf("Ошибка инициализации эмбеддера: %v", err)
	}
	ai.embedder = embedder
	log.Printf("Эмбеддер инициализирован: %s", ollamaEmbModel)

	// Инициализация векторного хранилища
	store, err := pgvector.New(
		ctx,
		pgvector.WithConnectionURL(databaseURL),
		pgvector.WithEmbedder(ai.embedder),
	)
	if err != nil {
		log.Fatalf("Ошибка инициализации векторного хранилища: %v", err)
	}
	ai.vectorStore = &store
	log.Println("Векторное хранилище инициализировано")

	return ai
}

func (ai *ragLangchainTest) AddDocument(ctx context.Context, content string) {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	// Создание текст-сплиттера для чанкинга
	splitter := textsplitter.NewRecursiveCharacter(
		textsplitter.WithChunkSize(512),
		textsplitter.WithChunkOverlap(50),
	)

	// Разделение содержимого на чанки
	chunks, err := splitter.SplitText(content)
	if err != nil {
		log.Fatalf("Ошибка разделения текста: %v", err)
	}
	log.Printf("Сгенерировано %d чанков", len(chunks))

	// Подготовка документов для вставки
	docs := make([]schema.Document, len(chunks))
	for i, chunk := range chunks {
		docs[i] = schema.Document{
			PageContent: chunk,
			//Metadata:    map[string]interface{}{"source": "moon_landing"},
		}
	}

	// Вставка документов в векторное хранилище
	_, err = ai.vectorStore.AddDocuments(ctx, docs)
	if err != nil {
		log.Fatalf("Ошибка вставки документов: %v", err)
	}
	log.Println("Документы добавлены в векторное хранилище")
}

func (ai *ragLangchainTest) Test(ctx context.Context, query string) string {
	ctx, cancel := context.WithTimeout(ctx, 120*time.Second)
	defer cancel()

	if query == "" {
		query = "Когда была высадка на Луну?"
	}

	// Поиск релевантных документов
	docs, err := ai.vectorStore.SimilaritySearch(ctx, query, 3)
	if err != nil {
		log.Fatalf("Ошибка поиска документов: %v", err)
	}
	log.Printf("Найдено %d релевантных документов", len(docs))
	for _, doc := range docs {
		log.Printf("Найденный документ: %s", doc.PageContent)
	}

	// Формирование аугментированного запроса
	augmentedQuery := "Контекст: " + combineDocs(docs) + "\nВопрос: " + query

	// Создание промпта
	prompt := "Вы - ИИ-ассистент." +
		"Используйте предоставленный контекст, чтобы максимально точно ответить на вопрос пользователя. " +
		"Если предоставленный контекст не связан с вопросом пользователя, не упоминайте его. " +
		"Все ответы переводи на русский язык, независимо от языка на котором задан вопрос.\n"
	// Генерация ответа
	//response, err := llms.GenerateFromSinglePrompt(ctx, ai.generation_llm, prompt+augmentedQuery)
	response, err := ai.generation_llm.GenerateContent(ctx, []llms.MessageContent{
		{Role: "system", Parts: []llms.ContentPart{llms.TextContent{Text: prompt}}},
		{Role: "human", Parts: []llms.ContentPart{llms.TextContent{Text: augmentedQuery}}},
	}, llms.WithMaxTokens(150), llms.WithTemperature(0.7)) //, llms.WithTopP(0.9)
	if err != nil {
		log.Fatalf("Ошибка генерации ответа: %v", err)
	}
	log.Printf("Ответ RAG: %s", response.Choices[0].Content)
	return response.Choices[0].Content
}

// Вспомогательная функция для объединения документов
func combineDocs(docs []schema.Document) string {
	var combined string
	for _, doc := range docs {
		combined += doc.PageContent + "\n"
	}
	return combined
}
