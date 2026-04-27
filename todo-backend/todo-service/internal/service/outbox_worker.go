package service

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/rs/zerolog/log"
	"github.com/todo-app/todo-service/internal/model"
	"github.com/todo-app/todo-service/internal/repository"
)

func jsonReader(b []byte) io.Reader { return bytes.NewReader(b) }

func parseESHits(res *esapi.Response) ([]model.Todo, error) {
	var result struct {
		Hits struct {
			Hits []struct {
				Source model.Todo `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}
	todos := make([]model.Todo, 0, len(result.Hits.Hits))
	for _, h := range result.Hits.Hits {
		todos = append(todos, h.Source)
	}
	return todos, nil
}

// RunOutboxWorker polls search_outbox every 5s and syncs to Elasticsearch.
func RunOutboxWorker(ctx context.Context, outboxRepo repository.OutboxRepository, es *elasticsearch.Client) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			processOutbox(ctx, outboxRepo, es)
		}
	}
}

func processOutbox(ctx context.Context, outboxRepo repository.OutboxRepository, es *elasticsearch.Client) {
	events, err := outboxRepo.FindUnprocessed(ctx, 50)
	if err != nil {
		log.Error().Err(err).Msg("outbox: fetch failed")
		return
	}
	for _, event := range events {
		var syncErr error
		switch event.Operation {
		case "upsert":
			req := esapi.IndexRequest{
				Index:      "todos",
				DocumentID: event.TodoID.String(),
				Body:       bytes.NewReader(event.Payload),
			}
			res, err := req.Do(ctx, es)
			if err == nil {
				res.Body.Close()
			}
			syncErr = err
		case "delete":
			req := esapi.DeleteRequest{
				Index:      "todos",
				DocumentID: event.TodoID.String(),
			}
			res, err := req.Do(ctx, es)
			if err == nil {
				res.Body.Close()
			}
			syncErr = err
		}
		if syncErr != nil {
			log.Error().Err(syncErr).Str("todo_id", event.TodoID.String()).Msg("outbox: sync failed, will retry")
			continue
		}
		if err := outboxRepo.MarkProcessed(ctx, event.ID); err != nil {
			log.Error().Err(err).Msg("outbox: mark processed failed")
		}
	}
}
