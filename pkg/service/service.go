package service

import (
	"context"
	"math/rand"
	"shortener/pkg"
	"shortener/pkg/repository"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

type Repository interface {
	Insert(ctx context.Context, shortURL *pkg.ShortURL) error
}

type Service struct {
	rnd    *rand.Rand
	urlDAO Repository
}

func NewService(urlDAO *repository.UrlDAO) *Service {
	return &Service{
		rnd:    rand.New(rand.NewSource(time.Now().UnixNano())),
		urlDAO: urlDAO,
	}
}

func (s *Service) Shorten(ctx context.Context, url string, ttlDays int) (*pkg.ShortURL, error) {
	shortURL := &pkg.ShortURL{
		URL: url,
		//ExpireAt: getExpirationTime(ttlDays),
	}

	for it := 0; it < 10; it++ {
		shortURL.Id = s.generateRandomID()

		err := s.urlDAO.Insert(ctx, shortURL)
		if err == nil {
			return shortURL, nil
		}

		if !mongo.IsDuplicateKeyError(err) {
			return nil, err
		}
	}

	return nil, ErrCollision
}
