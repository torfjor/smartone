package main

import (
	"context"
	"log"
	"net/http"

	"github.com/torfjor/smartone"
)

type AccessController interface {
	Valid(ctx context.Context, code string) (bool, error)
}

type accessController struct{}

func (a *accessController) Valid(ctx context.Context, code string) (bool, error) {
	return true, nil
}

func main() {
	if err := run(&accessController{}); err != nil {
		log.Fatal(err)
	}
}

func run(ac AccessController) error {
	http.ListenAndServe(":8080", smartone.OnCodeHandler(func(ctx context.Context, r smartone.OnCodeRequest) (smartone.OnCodeResponse, error) {
		valid, err := ac.Valid(ctx, r.Code)
		if err != nil {
			return smartone.OnCodeResponse{}, err
		}
		return smartone.OnCodeResponse{
			Result: valid,
		}, nil
	}))

	return nil
}
