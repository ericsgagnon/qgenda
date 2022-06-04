package qgenda

import "context"

type Pipeline interface {

	Exec(ctx context.Context) error
}

// rqf -> pipeline
// pipeline get,process,load
