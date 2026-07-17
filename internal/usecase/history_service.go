package usecase

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/marcell322/api-test-studio/internal/domain/models"
	"github.com/marcell322/api-test-studio/internal/domain/repository"
)

// SendRequestResult is what gets returned to the caller after executing a request.
type SendRequestResult struct {
	StatusCode     int               `json:"status_code"`
	Headers        map[string]string `json:"headers"`
	Body           string            `json:"body"`
	ResponseTimeMs int64             `json:"response_time_ms"`
}

type HistoryService interface {
	Send(userID uint, method, rawURL string, headers map[string]string, body string) (*SendRequestResult, error)
	List(userID uint) ([]models.RequestHistory, error)
	Get(userID, id uint) (*models.RequestHistory, error)
	Delete(userID, id uint) error
}

type historyService struct {
	repo   repository.HistoryRepository
	client *http.Client
}

func NewHistoryService(r repository.HistoryRepository) HistoryService {
	return &historyService{
		repo:   r,
		client: &http.Client{Timeout: 15 * time.Second},
	}
}

func validateSendRequest(method, rawURL string) error {
	if !validHTTPMethods[strings.ToUpper(strings.TrimSpace(method))] {
		return errors.New("invalid http method")
	}
	if strings.TrimSpace(rawURL) == "" {
		return errors.New("url is required")
	}
	u, err := url.Parse(rawURL)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") {
		return errors.New("url must be a valid http or https url")
	}
	return nil
}

// Send executes the given HTTP request and records the outcome in history,
// including failed attempts (network errors, timeouts), so history reflects
// what actually happened rather than only successful calls.
func (s *historyService) Send(userID uint, method, rawURL string, headers map[string]string, body string) (*SendRequestResult, error) {
	if err := validateSendRequest(method, rawURL); err != nil {
		return nil, err
	}
	method = strings.ToUpper(strings.TrimSpace(method))

	req, err := http.NewRequest(method, rawURL, bytes.NewBufferString(body))
	if err != nil {
		return nil, err
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	start := time.Now()
	resp, err := s.client.Do(req)
	elapsed := time.Since(start).Milliseconds()

	entry := &models.RequestHistory{
		UserID:         userID,
		Method:         method,
		URL:            rawURL,
		ResponseTimeMs: elapsed,
	}

	if err != nil {
		entry.StatusCode = 0
		entry.Error = err.Error()
		_ = s.repo.Create(entry) // best-effort log even on failure
		return nil, errors.New("request failed: " + err.Error())
	}
	defer resp.Body.Close()

	respBody, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		respBody = []byte("")
	}

	entry.StatusCode = resp.StatusCode
	if err := s.repo.Create(entry); err != nil {
		return nil, err
	}

	respHeaders := map[string]string{}
	for k, v := range resp.Header {
		if len(v) > 0 {
			respHeaders[k] = v[0]
		}
	}

	return &SendRequestResult{
		StatusCode:     resp.StatusCode,
		Headers:        respHeaders,
		Body:           string(respBody),
		ResponseTimeMs: elapsed,
	}, nil
}

func (s *historyService) List(userID uint) ([]models.RequestHistory, error) {
	return s.repo.ListByUser(userID)
}

func (s *historyService) Get(userID, id uint) (*models.RequestHistory, error) {
	h, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if h == nil {
		return nil, ErrNotFound
	}
	if h.UserID != userID {
		return nil, ErrForbidden
	}
	return h, nil
}

func (s *historyService) Delete(userID, id uint) error {
	h, err := s.Get(userID, id)
	if err != nil {
		return err
	}
	return s.repo.Delete(h.ID)
}