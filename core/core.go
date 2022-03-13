package core

import (
	"authServer/models"
	"authServer/repository"
	"authServer/service"
	"authServer/utils"
	jwt2 "github.com/dgrijalva/jwt-go"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

type Server struct{}

var repo = repository.NewRepository()
var jwt = service.NewJWTService()

func (s *Server) Register(ctx context.Context, request *RegisterRequest) (*StatusResponse, error) {
	log.Debug("REQUEST REGISTER")
	var response = new(StatusResponse)
	response.Status = true

	var newUser = models.User{
		UserName: request.Username,
		Email:    request.Email,
		Password: utils.HashAndSalt(request.Password),
		Role:     request.Role,
	}
	err := repo.RegisterUser(newUser)
	if err != nil {
		response.Status = false
		log.Error("Error Register User ", request.Username)
	} else {
		log.Debug("User register with Success!!")
	}
	return response, err
}

func (s *Server) Login(ctx context.Context, request *LoginRequest) (*LoginResponse, error) {
	log.Debug("REQUEST LOGIN")
	userResponse, err := repo.ExistUser(&models.User{
		Email:    request.Email,
		Password: request.Password,
	})
	if err != nil {
		return &LoginResponse{}, err
	}
	var response = new(LoginResponse)
	var user = new(User)
	user.Username = userResponse.UserName

	response.User = user
	response.Token = jwt.GenerateToken(userResponse.Email, userResponse.Role)

	return response, nil
}

func (s *Server) ValidateToken(ctx context.Context, request *ValidateRequest) (*User, error) {
	log.Debug("REQUEST VALIDATE TOKEN")
	token, err := jwt.ValidateToken(request.Token)
	if err != nil {
		return &User{}, err
	}
	claims := token.Claims.(jwt2.MapClaims)
	return &User{Username: claims["email"].(string),
		Role: claims["role"].(string)}, nil
}

func (s *Server) ChangePassword(ctx context.Context, request *ChangePasswordRequest) (*StatusResponse, error) {
	log.Debug("REQUEST CHANGE PASSWORD")
	var pass = new(models.ChangePass)
	pass.Email = request.Email
	pass.OldPassword = request.OldPassword
	pass.NewPassword = request.NewPassword
	if err := repo.ChangePassword(pass); err != nil {
		return &StatusResponse{Status: false}, err
	}

	return &StatusResponse{Status: true}, nil
}
func (s *Server) DeleteUser(ctx context.Context, request *DeleteUserRequest) (*StatusResponse, error) {
	log.Debug("REQUEST DELETE USER")
	var delUser = new(models.User)
	delUser.Email = request.Email
	delUser.Password = request.Password

	if err := repo.DeleteUser(delUser); err != nil {
		return &StatusResponse{Status: false}, err
	}
	return &StatusResponse{Status: true}, nil
}
