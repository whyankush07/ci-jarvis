package rag

import (
	"context"
	"crypto/sha256"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"ci-jarvis/internal/llm"
	"ci-jarvis/internal/store/vector"
)

type Crawler struct {
	llmClient   *llm.Client
	vectorStore vector.VectorStore
}

func NewCrawler(llmClient *llm.Client, vectorStore vector.VectorStore) *Crawler {
	return &Crawler{
		llmClient:   llmClient,
		vectorStore: vectorStore,
	}
}

// IndexRepository walks through the root directory and indexes all relevant files.
func (c *Crawler) IndexRepository(ctx context.Context, root string) error {
	return filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip directories and hidden files
		if d.IsDir() || strings.HasPrefix(d.Name(), ".") {
			if d.IsDir() && (d.Name() == "vendor" || d.Name() == "node_modules" || d.Name() == ".git") {
				return filepath.SkipDir
			}
			return nil
		}

		// Only index code files (expand this list as needed)
		ext := filepath.Ext(path)
		if ext != ".go" && ext != ".txt" && ext != ".md" && ext != ".sql" {
			return nil
		}

		return c.indexFile(ctx, path)
	})
}

func (c *Crawler) indexFile(ctx context.Context, path string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", path, err)
	}

	chunks := c.chunkContent(string(content), 1000) // 1000 characters per chunk

	for i, chunk := range chunks {
		// Generate a unique ID for the chunk
		id := fmt.Sprintf("%x", sha256.Sum256([]byte(fmt.Sprintf("%s-%d", path, i))))

		vector, err := c.llmClient.EmbedText(ctx, chunk)
		if err != nil {
			return fmt.Errorf("failed to embed chunk from %s: %w", path, err)
		}

		payload := map[string]interface{}{
			"path":    path,
			"chunk":   i,
			"content": chunk,
		}

		err = c.vectorStore.Upsert(ctx, id, vector, payload)
		if err != nil {
			return fmt.Errorf("failed to upsert chunk to vector store: %w", err)
		}
	}

	return nil
}

// splits text into smaller pieces of roughly maxChunkSize.
func (c *Crawler) chunkContent(content string, maxChunkSize int) []string {
	var chunks []string
	lines := strings.Split(content, "\n")
	var currentChunk strings.Builder

	for _, line := range lines {
		if currentChunk.Len()+len(line) > maxChunkSize && currentChunk.Len() > 0 {
			chunks = append(chunks, currentChunk.String())
			currentChunk.Reset()
		}
		currentChunk.WriteString(line)
		currentChunk.WriteString("\n")
	}

	if currentChunk.Len() > 0 {
		chunks = append(chunks, currentChunk.String())
	}

	return chunks
}
