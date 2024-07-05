package service

import (
	"context"
	"fmt"
	"math/rand"
	configs "shortener/configs"
	domain "shortener/internal/domain"
	"shortener/internal/repository"

	"time"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
)

type Service interface {
	Shorten(ctx context.Context, url string, ttlDays int) ([]byte, error)
	GetFullURL(ctx context.Context, shortURL string) (*domain.ShortURL, error)
	Update(ctx context.Context, id, url string, ttl int) error
	Delete(ctx context.Context, shortUrl string) error
}

type service struct {
	rnd        *rand.Rand
	repository repository.Repository
}

func NewService(repository repository.Repository) Service {
	return &service{
		rnd:        rand.New(rand.NewSource(time.Now().UnixNano())),
		repository: repository,
	}
}

func (s *service) Shorten(ctx context.Context, url string, ttlDays int) ([]byte, error) {
	shortURL := &domain.ShortURL{
		URL:       url,
		ExpiredAt: getExpirationTime(ttlDays),
	}

	for i := 0; i < 10; i++ {
		shortURL.Id = s.generateRandomId()

		err := s.repository.Insert(ctx, shortURL)
		if err != nil {
			return nil, fmt.Errorf("Shorten.Insert: %w", err)
		}

		if !mongo.IsDuplicateKeyError(err) {
			break
		} else {
			return nil, fmt.Errorf("Shorten.IsDuplicateKeyError: %w", err)
		}
	}

	// display our short URL
	displayedMessage := []byte(fmt.Sprintf(viper.GetString("db.host") + viper.GetString("port") + "/" + shortURL.Id))

	return displayedMessage, nil
}

func (s *service) GetFullURL(ctx context.Context, shortURL string) (*domain.ShortURL, error) {
	return s.repository.FindByID(ctx, shortURL)
}

func (s *service) Update(ctx context.Context, shortURL, url string, ttl int) error {
	// find our full url with short url
	shortUrl, err := s.repository.FindByID(ctx, shortURL)
	if err != nil {
		return fmt.Errorf("Update.FindByID: %w", err)
	}

	shortUrl.URL = url
	// create new ttl if it exists
	shortUrl.ExpiredAt = getExpirationTime(ttl)

	return s.repository.Update(ctx, shortUrl)
}

func (s *service) Delete(ctx context.Context, shortUrl string) error {
	return s.repository.Delete(ctx, shortUrl)
}

func getExpirationTime(ttlDays int) *time.Time {
	if ttlDays <= 0 {
		return nil
	}
	t := time.Now().Add(time.Hour * configs.Hours * time.Duration(ttlDays))

	return &t
}

func (s *service) generateRandomId() string {
	id := make([]rune, configs.IdLength)

	for i := range id {
		id[i] = configs.Symbols[s.rnd.Intn(len(configs.Symbols))]
	}
	return string(id)
}
