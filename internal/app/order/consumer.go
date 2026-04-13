package order

import (
	"github.com/hueat/backend/internal/pkg/hueat_log"
	"github.com/hueat/backend/internal/pkg/hueat_pubsub"
	"go.uber.org/zap"
)

type orderConsumerInterface interface {
	subscribe()
}

type orderConsumer struct {
	pubSub  *hueat_pubsub.PubSubAgent
	service orderServiceInterface
}

func newOrderConsumer(pubSub *hueat_pubsub.PubSubAgent, service orderServiceInterface) orderConsumer {
	consumer := orderConsumer{
		pubSub:  pubSub,
		service: service,
	}
	return consumer
}

func (r orderConsumer) subscribe() {
	go func() {
		messageChannel := r.pubSub.Subscribe(hueat_pubsub.TopicTableV1)
		isChannelOpen := true
		for isChannelOpen {
			func() {
				defer func() {
					if r := recover(); r != nil {
						hueat_log.LogPanicError(r, "order-consumer", "Panic occurred in handling a new message")
					}
				}()
				msg, channelOpen := <-messageChannel
				if !channelOpen {
					isChannelOpen = false
					zap.L().Info(
						"Channel closed. No more events to listen... quit!",
						zap.String("service", "order-consumer"),
					)
					return
				}
				// ACK message
				defer msg.Message.EventState.Done()
				zap.L().Info(
					"Received Event Message",
					zap.String("service", "order-consumer"),
					zap.String("event-id", msg.Message.EventID.String()),
					zap.String("event-type", string(msg.Message.EventType)),
				)
				if msg.Message.EventType != hueat_pubsub.TableCreatedEvent {
					return
				}
				event := msg.Message.EventEntity.(*hueat_pubsub.TableEventEntity)
				// Create order on Table creation
				if err := r.service.createOrderFromEvent(*event); err != nil {
					zap.L().Error("Impossible to update orders from RS Engine event", zap.String("service", "order-consumer"))
					return
				}
			}()
		}
	}()
}
