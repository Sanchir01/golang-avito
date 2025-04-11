package pvz

type ServicePVZ interface{}

type Service struct {
	repository ServicePVZ
}

func NewService(repository ServicePVZ) *Service {
	return &Service{repository: repository}
}
