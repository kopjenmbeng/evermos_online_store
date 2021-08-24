package authentication

import (
	"net/http"
	"github.com/RoseRocket/xerrs"
	"github.com/kopjenmbeng/evermos_online_store/internal/utility/respond"
	"github.com/kopjenmbeng/evermos_online_store/internal/middleware"
)

func GetTokenHandler(w http.ResponseWriter, r *http.Request) {
	var (
		// err error
		rc = r.Context()
	)

	// app_code:=r.Header.Get("X-Client-id")
	phoneNumber:=r.URL.Query().Get("phone_number")
	var password string = r.URL.Query().Get("password")
	useCase := UseCaseFromContext(rc)
	data, code, err := useCase.GetToken(rc, phoneNumber, password)
	if err != nil {
		middleware.GetLogEntry(r).Error(xerrs.Details(err, respond.ErrMaxStack))
		respond.Nay(w, r, code, err)
		return
	}
	respond.Yay(w, r, http.StatusOK, data)
	return

}
