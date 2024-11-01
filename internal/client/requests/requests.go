package requests

import (
	"time"

	"github.com/MihailSergeenkov/GophKeeper/internal/client/config"
	"github.com/go-resty/resty/v2"
)

// Request определяет тип HTTP запроса.
type Request struct {
	r *resty.Request
}

// RequestOptionFunc определяет тип функции для опций.
type RequestOptionFunc func(*Request)

// NewRequests инициализатор для HTTP запросов.
func NewRequests(cfg *config.Config) *Request {
	r := resty.New().
		SetTimeout(time.Duration(cfg.GetRequestTimeout()) * time.Second).
		SetRetryCount(cfg.GetRequestRetry()).R()

	return &Request{r: r}
}

// Get функция для выполнения get HTTP запросов.
func (o *Request) Get(url string, opts ...RequestOptionFunc) (*resty.Response, error) {
	for _, opt := range opts {
		opt(o)
	}

	resp, err := o.r.Get(url)

	return resp, err //nolint:wrapcheck // Нужно обернуть, но возврат должен остаться оригинальным
}

// Post функция для выполнения post HTTP запросов.
func (o *Request) Post(url string, opts ...RequestOptionFunc) (*resty.Response, error) {
	for _, opt := range opts {
		opt(o)
	}

	resp, err := o.r.Post(url)

	return resp, err //nolint:wrapcheck // Нужно обернуть, но возврат должен остаться оригинальным
}

// WithHeader добавляет header к запросу.
func WithHeader(key, value string) RequestOptionFunc {
	return func(o *Request) {
		o.r.SetHeader(key, value)
	}
}

// WithResult добавляет возможность сохранения данных ответа.
func WithResult(resultObject any) RequestOptionFunc {
	return func(o *Request) {
		o.r.SetResult(resultObject)
	}
}

// WithBody добавляет возможность отправки тела запроса.
func WithBody(body any) RequestOptionFunc {
	return func(o *Request) {
		o.r.SetBody(body)
	}
}

// WithPathParams добавляет возможность использовать параметры в запросе.
func WithPathParams(params map[string]string) RequestOptionFunc {
	return func(o *Request) {
		o.r.SetPathParams(params)
	}
}

// WithOutput добавляет возможность сохранять ответ в файл.
func WithOutput(file string) RequestOptionFunc {
	return func(o *Request) {
		o.r.SetOutput(file)
	}
}

// WithFile добавляет возможность отправлять файл в запросе.
func WithFile(filePath string) RequestOptionFunc {
	return func(o *Request) {
		o.r.SetFile("file", filePath)
	}
}

// WithFormData добавляет возможность отправлять данные формы в запросе.
func WithFormData(data map[string]string) RequestOptionFunc {
	return func(o *Request) {
		o.r.SetFormData(data)
	}
}
