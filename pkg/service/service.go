package service

import (
	"context"
	"fmt"
	"math/rand"
	"shortener/pkg"
	"shortener/pkg/repository"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

type Repository interface {
	Insert(ctx context.Context, shortURL *pkg.ShortURL) error
	FindByID(ctx context.Context, id string) (*pkg.ShortURL, error)
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

func (s *Service) Shorten(ctx context.Context, url string, ttlDays int) (string, error) {
	shortURL := &pkg.ShortURL{
		URL:       url,
		ExpiredAt: getExpirationTime(ttlDays),
	}

	for i := 0; i < 10; i++ {
		shortURL.Id = s.generateRandomId()
		err := s.urlDAO.Insert(ctx, shortURL)
		if err != nil {
			return "shortURL", err
		}
		returnedUrl := fmt.Sprintf("localhost:8080/%s", shortURL.Id)

		if !mongo.IsDuplicateKeyError(err) {
			return returnedUrl, err
		}
	}

	return "", nil
}

func (s *Service) GetFullURL(ctx context.Context, shortURL string) (string, error) {
	sURL, err := s.urlDAO.FindByID(ctx, shortURL)
	if err != nil {
		return "", err
	}
	return sURL.URL, nil
}

func getExpirationTime(ttlDays int) *time.Time {
	if ttlDays <= 0 {
		return nil
	}
	t := time.Now().Add(time.Hour * 24 * time.Duration(ttlDays))
	return &t
}

var symbols = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

const idLength = 6

func (s *Service) generateRandomId() string {
	id := make([]rune, idLength)

	for i := range id {
		id[i] = symbols[s.rnd.Intn(len(symbols))]
	}
	return string(id)
}
