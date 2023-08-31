package history

import "time"

type Service interface {
	WriteToFile(userID int, segmentTitle string, operationName string, date time.Time)
}

type ServiceStub struct {
}

func NewServiceStub() *ServiceStub {
	return &ServiceStub{}
}

func (s *ServiceStub) WriteToFile(userID int, segmentTitle string, operationName string, date time.Time) {
}
