package server

import "fmt"

// 部署svc的服务器

type Server struct {
	ID        int64  `json:"id"`
	Summary   string `json:"summary"`
	Host      string `json:"host"`
	Port      int    `json:"port"`
	Status    int    `json:"status"`
	CreatedAt int64  `json:"created_at"`
	DeletedAt *int64 `json:"deleted_at"`
}

func (s *Server) URL() string {
	return fmt.Sprintf("http://%s:%d", s.Host, s.Port)
}
