package service

import "github.com/longjoy/micro-service/domain/repository"

type PermitService struct {
	UserRepository repository.UserRepository
}

func NewPermitService(ury repository.UserRepository) PermitService {
	return PermitService{UserRepository: ury}
}
