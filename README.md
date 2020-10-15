# Tidomat smartONE handlers

Exposes a set of `http.HandlerFunc`s for decoding requests and encoding responses 
to [Tidomat smartONE](https://www.tidomat.se/smartone/en/hw.asp?id=SO-3303&s=2&l=us).

## OnCodeHandler
Calls the given `OnCodeFunc` with the decoded request and expects `(OnCodeResponse, error)` in return.

### Example
`cmd/entry/main.go`

````go
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

type accessController struct {}

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

````

## OnCardHandler
Not implemented yet.