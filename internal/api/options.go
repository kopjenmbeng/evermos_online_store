package api

import(
	"github.com/kopjenmbeng/evermos_online_store/internal/middleware/jwe_auth"
)

func JWE(jw *jwe_auth.JWE) Option {
	return func(s *Server) {
		s.jwe = jw
	}
}