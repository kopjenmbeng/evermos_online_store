package respond

import (
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

const ErrMaxStack = 5

type (
	Response struct {
		RequestId string         `json:"request_id"`
		Content   interface{}    `json:"content,omitempty"`
		Error     *ErrorWithCode `json:"error,omitempty"`
		Status    int            `json:"status"`
	}
	ErrorWithCode struct {
		Code    int               `json:"code,omitempty"`
		Message string            `json:"message"`
		Reasons map[string]string `json:"reasons,omitempty"`
	}
)

func (res *Response) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func Yay(w http.ResponseWriter, r *http.Request, status int, content interface{}) {
	render.Status(r, status)
	_ = render.Render(w, r, &Response{
		RequestId: middleware.GetReqID(r.Context()),
		Content:   content,
		Status:    status,
	})
}

func Nay(w http.ResponseWriter, r *http.Request, status int, err error) {
	render.Status(r, status)
	_ = render.Render(w, r, &Response{
		RequestId: middleware.GetReqID(r.Context()),
		Status:    status,
		Error: &ErrorWithCode{
			Code:    status,
			Message: err.Error(),
		},
	})
}
