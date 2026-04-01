package memory

import (
	"context"
	"encoding/json"

	"github.com/kocar/aurelia/internal/config"
	"github.com/kocar/aurelia/internal/infra"
	"github.com/redis/go-redis/v9"
)

type Chunk struct {
	Text  string  `json:"text"`
	Score float64 `json:"score"`
}

type SharedMemory interface {
	Publish(ctx context.Context, topic string, msg Message) error
	Subscribe(ctx context.Context, topic string) <-chan Message
	SearchCross(ctx context.Context, query string, k int) ([]Chunk, error)
}

type SharedMemoryRedis struct {
	redisClient *redis.Client
	cfg         *config.AppConfig
}

func NewSharedMemoryRedis(redisProv *infra.RedisProvider, cfg *config.AppConfig) *SharedMemoryRedis {
	return &SharedMemoryRedis{
		redisClient: redisProv.Client,
		cfg:         cfg,
	}
}

func (s *SharedMemoryRedis) Publish(ctx context.Context, topic string, msg Message) error {
	if s.cfg != nil && !s.cfg.Features.KAIROSEnabled {
		return nil
	}
	b, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return s.redisClient.Publish(ctx, topic, string(b)).Err()
}

func (s *SharedMemoryRedis) Subscribe(ctx context.Context, topic string) <-chan Message {
	ch := make(chan Message)
	if s.cfg != nil && !s.cfg.Features.KAIROSEnabled {
		close(ch)
		return ch
	}

	pubsub := s.redisClient.Subscribe(ctx, topic)
	go func() {
		defer close(ch)
		defer pubsub.Close()
		for {
			msg, err := pubsub.ReceiveMessage(ctx)
			if err != nil {
				return
			}
			var m Message
			if err := json.Unmarshal([]byte(msg.Payload), &m); err == nil {
				ch <- m
			}
		}
	}()
	return ch
}

func (s *SharedMemoryRedis) SearchCross(ctx context.Context, query string, k int) ([]Chunk, error) {
	if s.cfg != nil && !s.cfg.Features.KAIROSEnabled {
		return nil, nil
	}

	httpClient := NewSemanticHTTPClient(0)
	vector, err := EmbedText(ctx, httpClient, s.cfg.OllamaURL, s.cfg.QdrantEmbeddingModel, query)
	if err != nil {
		return nil, err
	}

	points, err := SearchSemantic(ctx, httpClient, s.cfg.QdrantURL, "aurelia_memory", "", vector, k)
	if err != nil {
		return nil, err
	}

	chunks := make([]Chunk, 0, len(points))
	for _, p := range points {
		chunks = append(chunks, Chunk{
			Text:  ExtractSearchableText(p.Payload),
			Score: p.Score,
		})
	}
	return chunks, nil
}
