package app

import "context"

type WebApp interface { // TODO
	Start(port string) error
	Stop(ctx context.Context) error
}
