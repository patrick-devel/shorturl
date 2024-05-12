package models

import (
	"encoding/json"
	"net/url"
)

type Request struct {
	URL url.URL `json:"url"`
}

type RequestBulk struct {
	OriginalURL   url.URL `json:"original_url"`
	CorrelationID string  `json:"correlation_id"`
}

func (r *RequestBulk) UnmarshalJSON(data []byte) error {
	type ReqAlias RequestBulk

	aliasValue := struct {
		ReqAlias
		OriginalURL   string `json:"original_url"`
		CorrelationID string `json:"correlation_id"`
	}{
		ReqAlias: ReqAlias(*r),
	}

	if err := json.Unmarshal(data, &aliasValue); err != nil {
		return err
	}

	uri, err := url.ParseRequestURI(aliasValue.OriginalURL)
	if err != nil {
		return err
	}

	r.OriginalURL = *uri
	r.CorrelationID = aliasValue.CorrelationID

	return nil
}

type ListRequestBulk []RequestBulk

type ResponseBulk struct {
	ShortURL      string `json:"short_url"`
	CorrelationID string `json:"correlation_id"`
}

type ListResponseBulk []ResponseBulk

func (r *Request) UnmarshalJSON(data []byte) error {
	type ReqAlias Request

	aliasValue := struct {
		ReqAlias
		URL string `json:"url"`
	}{
		ReqAlias: ReqAlias(*r),
	}

	if err := json.Unmarshal(data, &aliasValue); err != nil {
		return err
	}

	uri, err := url.ParseRequestURI(aliasValue.URL)
	if err != nil {
		return err
	}

	r.URL = *uri

	return nil
}

type Response struct {
	Result string `json:"result"`
}

type Event struct {
	UUID        string `json:"uuid"`
	CreatorID   string `json:"creator_id"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type ResponseGetURLs struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}
