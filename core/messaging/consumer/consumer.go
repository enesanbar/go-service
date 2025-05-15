package consumer

type Consumer interface {
	Start() error
	Stop() error
}
