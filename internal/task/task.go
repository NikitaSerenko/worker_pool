package task

import "context"

type Service interface {
	Name() string
	Process(ctx context.Context) error
}
