package repositories

import (
	"context"
	"fmt"
	domain "parser/internal/domain/models"
	"parser/pkg/postgres"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

type SubscriberRepository interface {
	InsertSubscriber(ctx context.Context, sub domain.Subscriber) error
	InsertSubscription(ctx context.Context, sub *domain.Subscription) error
	GetAdvertSubscribers(ctx context.Context, advertID string) ([]*domain.Subscriber, error)
}

type subscriberRepo struct {
	db *postgres.Postgres
}

func NewRepo(db *postgres.Postgres) SubscriberRepository {
	return &subscriberRepo{db: db}
}

func (s *subscriberRepo) GetAdvertSubscribers(ctx context.Context, advertID string) ([]*domain.Subscriber, error) {

	sql, args, err := sq.Select("sub.subscriber_id, sub.telegram_id").
		From("subscriptions sp").
		Join("subscribers sub on sub.subscriber_id = sp.subscriber_id").
		Join("adverts ads on sp.advert_id = ads.advert_id").
		Where("ads.advert_id = $1", advertID).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return nil, err
	}

	rows, release, err := s.db.Query(ctx, sql, args)
	if err != nil {
		return nil, fmt.Errorf("internal error: %w", err)
	}

	defer release()

	var subscribers []*domain.Subscriber
	err = s.db.ScanAll(rows, &subscribers)
	if err != nil {
		return nil, fmt.Errorf("internal error: %w", err)
	}

	return subscribers, nil
}

func (s *subscriberRepo) InsertSubscription(ctx context.Context, sub *domain.Subscription) error {

	sql, args, err := sq.Insert("subscriptions").
		Columns("advert_id", "subscriber_id").
		Values(sub.AdvertID, sub.SubscriberID).
		ToSql()

	if err != nil {
		return err
	}
	_, release, err := s.db.Query(ctx, sql, args)
	if err != nil {
		return fmt.Errorf("internal error: %w", err)
	}
	defer release()

	return nil
}

func (s *subscriberRepo) InsertSubscriber(ctx context.Context, sub domain.Subscriber) error {

	sub.ID = uuid.NewString()

	sql, args, err := sq.Insert("subscribers").
		Columns("id", "telegram_id").
		Values(sub.ID, sub.TelegramID).
		ToSql()

	if err != nil {
		return err
	}
	_, release, err := s.db.Query(ctx, sql, args)
	if err != nil {
		return fmt.Errorf("internal error: %w", err)
	}
	defer release()

	return nil
}