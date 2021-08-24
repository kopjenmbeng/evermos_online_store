package authentication

import (
	"context"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/kopjenmbeng/evermos_online_store/internal/dto"
	"github.com/kopjenmbeng/evermos_online_store/internal/middleware"
	"github.com/kopjenmbeng/evermos_online_store/internal/middleware/jwe_auth"
)

type IAuthenticationUseCase interface {
	GetToken(ctx context.Context, phone_number string, password string) (token *TokenResponse, status int, Err error)
	Register(ctx context.Context,req RegisterRequest)(status int,err error)
}

type AuthenticationUseCase struct {
	repository IAuthenticationRepository
	r *http.Request
}

func NewAuthenticationUseCase(repository IAuthenticationRepository,r *http.Request) IAuthenticationUseCase {
	return &AuthenticationUseCase{repository: repository,r: r}
}

func(use_case *AuthenticationUseCase)Register(ctx context.Context,req RegisterRequest)(status int,err error){

	salt:=uuid.New().String()
	iteration:=middleware.GenerateRandomNumber(900,999)
	secLength:=middleware.GenerateRandomNumber(32,64)
	password:=middleware.HashPassword(req.Password,salt,iteration,secLength)
	cus:=dto.Customer{CustomerId: uuid.New().String(),PhoneNumber: req.PhoneNumber,FullName: req.FullName,Salt: salt,Password: password,Iteration: iteration,SecurityLength: secLength}
	status,err=use_case.repository.Register(ctx,cus)
	return
}
func (use_case *AuthenticationUseCase) GetToken(ctx context.Context, phone_number string, password string) (token *TokenResponse, status int, Err error) {
	data, status, err := use_case.repository.GetUser(ctx, phone_number)
	if err!=nil{
		return nil,status,err
	}
	var result TokenResponse
	if status == http.StatusOK {
		paramPassword := middleware.HashPassword(password, data.Salt, data.Iteration, data.SecurityLength)
		if data.Password != paramPassword {
			return nil, http.StatusUnauthorized, errors.New("password yang anda masukan salah.")
		}
		token,Exp,_:=jwe_auth.GenerateToken(use_case.r,data.CustomerId,false)
		result.Token=token
		result.Expiry=Exp
	
	}
	return &result, 0, nil
}
