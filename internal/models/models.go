package models

import (
	"encoding/json"
	"net/url"
)

type Request struct {
	URL url.URL `json:"url"`
}

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
	ShortUrl    string `json:"short_url"`
	OriginalUrl string `json:"original_url"`
}
