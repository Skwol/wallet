package common

import "context"

type Storage interface {
	GenerateFakeData(context.Context) error
}
