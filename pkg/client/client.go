package client

import "context"

type Client interface {
	Connect() (*Client, error)
	Do(ctx context.Context, r *request.Request) ([]Data, error)
}

type REST struct {
}

func (r *REST) Connect() (*Client, error) {

}

