package service

type ApplicationService interface {
	Start() error
	Stop() error
}
