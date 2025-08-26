package usecase

import (
	"context"
	"database/sql"
	"time"

	"github.com/MikhailKopeikin/my-go-api/domain"
	"github.com/MikhailKopeikin/my-go-api/transport/request"
	"github.com/MikhailKopeikin/my-go-api/utils"
	"github.com/MikhailKopeikin/my-go-api/utils/crypto"
	"github.com/MikhailKopeikin/my-go-api/utils/jwt"
)

type authUsecase struct {
	userRepo       domain.UserRepository
	cryptoSvc      crypto.CryptoService
	jwtSvc         jwt.JWTService
	contextTimeout time.Duration
}

// NewAuthUsecase will create new an authUsecase object representation of AuthUsecase interface
func NewAuthUsecase(userRepo domain.UserRepository, cryptoSvc crypto.CryptoService, jwtSvc jwt.JWTService, contextTimeout time.Duration) *authUsecase {
	return &authUsecase{
		userRepo:       userRepo,
		cryptoSvc:      cryptoSvc,
		jwtSvc:         jwtSvc,
		contextTimeout: contextTimeout,
	}
}

func (u *authUsecase) SignUp(c context.Context, request *request.SignUpReq) (err error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	user, err := u.userRepo.GetByEmail(ctx, request.Email)
	if err != nil && err != sql.ErrNoRows {
		return
	}

	if user.ID != 0 {
		err = utils.NewBadRequestError("email already registered")
		return
	}

	passwordHash, err := u.cryptoSvc.CreatePasswordHash(ctx, request.Password)
	if err != nil {
		return
	}

	err = u.userRepo.Create(ctx, &domain.User{
		Email:     request.Email,
		Password:  passwordHash,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	return
}

func (u *authUsecase) SignIn(c context.Context, request *request.SignInReq) (accessToken string, err error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	user, err := u.userRepo.GetByEmail(ctx, request.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			err = utils.NewBadRequestError("email and password not match")
			return
		}
		return
	}

	if !u.cryptoSvc.ValidatePassword(ctx, user.Password, request.Password) {
		err = utils.NewBadRequestError("email and password not match")
		return
	}

	accessToken, err = u.jwtSvc.GenerateToken(ctx, user.ID)
	return
}
