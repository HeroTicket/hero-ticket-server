package ipfs

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
)

type Service interface {
	PinFile(ctx context.Context, file io.Reader, filename string) (*PinFileResponse, error)
}

type IpfsServiceConfig struct {
	ApiKey string
	Secret string
	Client *http.Client
}

type IpfsService struct {
	apiKey string
	secret string
	client *http.Client
}

func New(cfg IpfsServiceConfig) Service {
	svc := &IpfsService{
		apiKey: cfg.ApiKey,
		secret: cfg.Secret,
		client: http.DefaultClient,
	}

	if cfg.Client != nil {
		svc.client = cfg.Client
	}

	return svc
}

func (svc *IpfsService) PinFile(ctx context.Context, file io.Reader, filename string) (*PinFileResponse, error) {
	body := &bytes.Buffer{}

	m := multipart.NewWriter(body)

	part, err := m.CreateFormFile("file", filename)
	if err != nil {
		return nil, err
	}

	if _, err := io.Copy(part, file); err != nil {
		return nil, err
	}

	m.Close()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, PinFileToIPFSUrl, body)
	if err != nil {
		return nil, err
	}

	fmt.Println("Content-Type:", m.FormDataContentType())

	req.Header.Set("Content-Type", m.FormDataContentType())
	req.Header.Set("pinata_api_key", svc.apiKey)
	req.Header.Set("pinata_secret_api_key", svc.secret)

	resp, err := svc.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errorFromResponse(resp)
	}

	var data PinFileResponse

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	return &data, nil
}

func errorFromResponse(resp *http.Response) error {
	var data map[string]interface{}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return err
	}

	if msg, ok := data["error"].(string); ok {
		return fmt.Errorf("pinata error: %s", msg)
	}

	if m, ok := data["error"].(map[string]interface{}); ok {
		if msg, ok := m["details"].(string); ok {
			return fmt.Errorf("pinata error: %s", msg)
		}
	}

	return ErrUnknown
}
