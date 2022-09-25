package service

import (
	"time"

	"github.com/nmluci/go-backend/pkg/dto"
)

func (s *service) Ping() (pingResponse dto.PublicPingResponse) {
	return dto.PublicPingResponse{
		Message:   "KyaaNakaWaZettaiDame",
		Timestamp: time.Now().Unix(),
	}
}
