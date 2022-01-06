package common

import "context"

type Storage interface {
	GenerateFakeData(context.Context, int) error
}
