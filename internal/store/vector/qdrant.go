package vector

import (
	"context"
	"fmt"
	"strings"

	"github.com/qdrant/go-client/qdrant"
)

type VectorStore interface {
	Upsert(ctx context.Context, id string, vector []float32, payload map[string]interface{}) error
	Search(ctx context.Context, vector []float32, limit uint64) ([]map[string]interface{}, error)
}

type QdrantStore struct {
	client     *qdrant.Client
	collection string
}

func NewQdrantStore(url string, apiKey string, collection string) (*QdrantStore, error) {
	// Strip protocol from URL if present
	host := url
	host = strings.TrimPrefix(host, "https://")
	host = strings.TrimPrefix(host, "http://")
	host = strings.Split(host, ":")[0] // Strip port if present

	client, err := qdrant.NewClient(&qdrant.Config{
		Host:   host,
		Port:   6334,
		APIKey: apiKey,
		UseTLS: true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create qdrant client: %w", err)
	}

	return &QdrantStore{
		client:     client,
		collection: collection,
	}, nil
}

func (q *QdrantStore) CreateCollection(ctx context.Context, vectorSize uint64) error {
	err := q.client.CreateCollection(ctx, &qdrant.CreateCollection{
		CollectionName: q.collection,
		VectorsConfig: qdrant.NewVectorsConfig(&qdrant.VectorParams{
			Size:     vectorSize,
			Distance: qdrant.Distance_Cosine,
		}),
	})
	if err != nil {
		return fmt.Errorf("failed to create collection: %w", err)
	}
	return nil
}

func (q *QdrantStore) Upsert(ctx context.Context, id string, vector []float32, payload map[string]interface{}) error {
	qPayload := make(map[string]*qdrant.Value)
	for k, v := range payload {
		qPayload[k], _ = qdrant.NewValue(v)
	}

	point := &qdrant.PointStruct{
		Id:      qdrant.NewID(id),
		Vectors: qdrant.NewVectors(vector...),
		Payload: qPayload,
	}

	_, err := q.client.Upsert(ctx, &qdrant.UpsertPoints{
		CollectionName: q.collection,
		Points:         []*qdrant.PointStruct{point},
	})
	if err != nil {
		return fmt.Errorf("failed to upsert point: %w", err)
	}

	return nil
}

func (q *QdrantStore) Search(ctx context.Context, vector []float32, limit uint64) ([]map[string]interface{}, error) {
	res, err := q.client.Query(ctx, &qdrant.QueryPoints{
		CollectionName: q.collection,
		Query:          qdrant.NewQuery(vector...),
		Limit:          &limit,
		WithPayload:    qdrant.NewWithPayload(true),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to search points: %w", err)
	}

	results := make([]map[string]interface{}, 0, len(res))
	for _, hit := range res {
		results = append(results, convertPayload(hit.Payload))
	}

	return results, nil
}

func convertPayload(qPayload map[string]*qdrant.Value) map[string]interface{} {
	res := make(map[string]interface{})
	for k, v := range qPayload {
		switch v.GetKind().(type) {
		case *qdrant.Value_StringValue:
			res[k] = v.GetStringValue()
		case *qdrant.Value_IntegerValue:
			res[k] = v.GetIntegerValue()
		case *qdrant.Value_DoubleValue:
			res[k] = v.GetDoubleValue()
		case *qdrant.Value_BoolValue:
			res[k] = v.GetBoolValue()
		default:
			res[k] = v.GetKind()
		}
	}
	return res
}
