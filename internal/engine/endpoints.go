package engine

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/alwindoss/ark/internal/vault"
	"github.com/go-chi/chi/v5"
	"github.com/go-kit/kit/endpoint"
)

type saveRequest struct {
	Key   string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
}

type saveResponse struct {
	Key string `json:"key,omitempty"`
	Err string `json:"err,omitempty"` // errors don't JSON-marshal, so we use a string
}

func makeSaveEndpoint(svc vault.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(saveRequest)
		val := strings.NewReader(req.Value)
		err := svc.Save([]byte(req.Key), val)
		if err != nil {
			return saveResponse{
				Key: req.Key,
				Err: err.Error(),
			}, nil
		}
		return saveResponse{
			Key: req.Key,
			Err: ""}, nil
	}
}

type retrieveRequest struct {
	Key string `json:"key,omitempty"`
}

type retrieveResponse struct {
	Key   string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
	Err   string `json:"err,omitempty"` // errors don't JSON-marshal, so we use a string
}

func makeRetrieveEndpoint(svc vault.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(retrieveRequest)
		respReader, err := svc.Retrieve([]byte(req.Key))
		if err != nil {
			return retrieveResponse{
				Key:   req.Key,
				Value: "",
				Err:   err.Error()}, nil
		}
		var buff strings.Builder
		n, err := io.Copy(&buff, respReader)
		if err != nil {
			return retrieveResponse{
				Key:   req.Key,
				Value: "",
				Err:   err.Error()}, nil
		}
		log.Printf("copied %d bytes", n)
		return retrieveResponse{
			Key:   req.Key,
			Value: buff.String(),
		}, nil
	}
}

func decodeSaveRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request saveRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func encodeSaveResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

func encodeRetrieveResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

func decodeRetrieveRequest(_ context.Context, r *http.Request) (interface{}, error) {
	key := chi.URLParam(r, "key")
	var request retrieveRequest
	request.Key = key
	return request, nil
}
