package service

import (
	"context"
	"fmt"

	"github.com/khiemnd777/noah_framework/internal/legacy/shared/logger"
	"github.com/khiemnd777/noah_framework/internal/legacy/shared/modules/notification/model"
	"github.com/khiemnd777/noah_framework/internal/legacy/shared/pubsub"
)

func (s *NotificationService) InitPubSubEvents() {
	pubsub.SubscribeAsync("notification:notify", func(payload *model.NotifyRequest) error {
		ctx := context.Background()
		if _, err := s.Create(ctx, payload.MessageID, payload.UserID, payload.NotifierID, payload.Type, payload.Data); err != nil {
			logger.Error(fmt.Sprintf("❌ Failed to create notification: %v", err))
		}
		return nil
	})
}
