package menuItem

import (
	"github.com/hueat/backend/internal/pkg/hueat_log"
	"github.com/hueat/backend/internal/pkg/hueat_pubsub"
	"go.uber.org/zap"
)

type menuItemConsumerInterface interface {
	subscribe()
}

type menuItemConsumer struct {
	pubSub  *hueat_pubsub.PubSubAgent
	service menuItemServiceInterface
}

func newMenuItemConsumer(pubSub *hueat_pubsub.PubSubAgent, service menuItemServiceInterface) menuItemConsumer {
	return menuItemConsumer{
		pubSub:  pubSub,
		service: service,
	}
}

func (r menuItemConsumer) subscribe() {
	go func() {
		messageChannel := r.pubSub.Subscribe(hueat_pubsub.TopicMenuOptionV1)
		isChannelOpen := true
		for isChannelOpen {
			func() {
				defer func() {
					if r := recover(); r != nil {
						hueat_log.LogPanicError(r, "menu-item-consumer", "Panic occurred in handling a new message")
					}
				}()
				msg, channelOpen := <-messageChannel
				if !channelOpen {
					isChannelOpen = false
					zap.L().Info(
						"Channel closed. No more events to listen... quit!",
						zap.String("service", "menu-item-consumer"),
					)
					return
				}
				defer msg.Message.EventState.Done()
				zap.L().Info(
					"Received Event Message",
					zap.String("service", "menu-item-consumer"),
					zap.String("event-id", msg.Message.EventID.String()),
					zap.String("event-type", string(msg.Message.EventType)),
				)
				if msg.Message.EventType != hueat_pubsub.MenuOptionCreatedEvent {
					return
				}
				event := msg.Message.EventEntity.(*hueat_pubsub.MenuOptionEventEntity)
				if err := r.service.resetMenuItemPriceFromEvent(*event); err != nil {
					zap.L().Error("Impossible to reset menu item price from menu option created event", zap.String("service", "menu-item-consumer"))
					return
				}
			}()
		}
	}()
}
