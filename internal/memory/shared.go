package memory

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// SharedMemory gerencia a comunicação Pub/Sub baseada em Redis
// para a coordenação do swarm local de agentes.
type SharedMemory struct {
	client *redis.Client
}

// NewSharedMemory cria uma nova instância de SharedMemory.
func NewSharedMemory(client *redis.Client) *SharedMemory {
	return &SharedMemory{
		client: client,
	}
}

// Publish envia uma mensagem msg para o canal channel especificado.
func (sm *SharedMemory) Publish(ctx context.Context, channel, msg string) error {
	return sm.client.Publish(ctx, channel, msg).Err()
}

// Subscribe assina um canal e retorna um canal de leitura (<-chan string)
// por onde as mensagens recebidas serão entregues.
func (sm *SharedMemory) Subscribe(ctx context.Context, channel string) <-chan string {
	pubsub := sm.client.Subscribe(ctx, channel)
	ch := make(chan string)

	go func() {
		defer pubsub.Close()
		defer close(ch)

		redisCh := pubsub.Channel()
		for {
			select {
			case <-ctx.Done():
				return
			case msg, ok := <-redisCh:
				if !ok {
					return
				}
				select {
				case <-ctx.Done():
					return
				case ch <- msg.Payload:
				}
			}
		}
	}()

	return ch
}
