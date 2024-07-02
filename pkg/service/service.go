package service

import (
	"context"
	"fmt"
	"math/rand"
	domain "shortener/pkg/domain"
	"shortener/pkg/repository"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

// type Service interface {
// 	Insert(ctx context.Context, shortURL *pkg.ShortURL) error
// 	FindByID(ctx context.Context, id string) (*pkg.ShortURL, error)
// 	Update(ctx context.Context, shortUrl *pkg.ShortURL) error
// }

type Service interface {
	Shorten(ctx context.Context, url string, ttlDays int) (string, error)
	GetFullURL(ctx context.Context, shortURL string) (string, error)
	Update(ctx context.Context, id, url string, ttl int) (string, error)
	Delete(ctx context.Context) (string, error)
}

type service struct {
	rnd    *rand.Rand
	urlDAO repository.Repository
}

func NewService(urlDAO repository.Repository) Service {
	return &service{
		rnd:    rand.New(rand.NewSource(time.Now().UnixNano())),
		urlDAO: urlDAO,
	}
}

func (s *service) Shorten(ctx context.Context, url string, ttlDays int) (string, error) {
	shortURL := &domain.ShortURL{
		URL:       url,
		ExpiredAt: getExpirationTime(ttlDays),
	}

	for i := 0; i < 10; i++ {
		shortURL.Id = s.generateRandomId()
		err := s.urlDAO.Insert(ctx, shortURL)
		if err != nil {
			return "", err
		}
		returnedUrl := fmt.Sprintf("localhost:8080/%s", shortURL.Id)

		if !mongo.IsDuplicateKeyError(err) {
			return returnedUrl, err
		}
	}

	return "", nil
}

func (s *service) GetFullURL(ctx context.Context, shortURL string) (string, error) {
	sURL, err := s.urlDAO.FindByID(ctx, shortURL)
	if err != nil {
		return "", err
	}

	return sURL.URL, nil
}

func (s *service) Update(ctx context.Context, id, url string, ttl int) (string, error) {
	// find our full url with short url
	shortUrl, err := s.urlDAO.FindByID(ctx, id)
	if err != nil {
		return "", fmt.Errorf("Update/FindByID: %w", err)
	}

	shortUrl.URL = url
	// create new ttl if it exists
	shortUrl.ExpiredAt = getExpirationTime(ttl)

	return "http.StatusOK", s.urlDAO.Update(ctx, shortUrl)
}

func (s *service) Delete(ctx context.Context) (string, error) {

	return "", nil
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

func (s *service) generateRandomId() string {
	id := make([]rune, idLength)

	for i := range id {
		id[i] = symbols[s.rnd.Intn(len(symbols))]
	}
	return string(id)
}
