package tools

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/vxcontrol/langchaingo/schema"
	"github.com/vxcontrol/langchaingo/vectorstores"
	"github.com/vxcontrol/langchaingo/vectorstores/pgvector"
)

// vectorSearch executes multiple similarity-search queries against a pgvector store,
// deduplicates results by document content, and returns up to limit documents sorted
// by descending score. Errors from individual queries are logged and skipped so that a
// single bad query does not abort the whole batch.
func vectorSearch(
	ctx context.Context,
	logger *logrus.Entry,
	store *pgvector.Store,
	questions []string,
	limit int,
	threshold float32,
	filters map[string]any,
) []schema.Document {
	var allDocs []schema.Document
	for i, query := range questions {
		queryLogger := logger.WithFields(logrus.Fields{
			"query_index": i + 1,
			"query":       query[:min(len(query), 1000)],
		})

		docs, err := store.SimilaritySearch(
			ctx,
			query,
			limit,
			vectorstores.WithScoreThreshold(threshold),
			vectorstores.WithFilters(filters),
		)
		if err != nil {
			queryLogger.WithError(err).Error("failed to search for similar documents")
			continue
		}

		queryLogger.WithField("docs_found", len(docs)).Debug("query executed")
		allDocs = append(allDocs, docs...)
	}

	logger.WithField("total_docs_before_dedup", len(allDocs)).Debug("all queries completed")

	docs := MergeAndDeduplicateDocs(allDocs, limit)

	logger.WithField("docs_after_dedup", len(docs)).Debug("documents deduplicated and sorted")

	return docs
}
