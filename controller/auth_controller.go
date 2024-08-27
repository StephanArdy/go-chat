package controller

import (
	"encoding/json"
	"go-chat/dto"
	"go-chat/service"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type AuthController interface {
	Register(w http.ResponseWriter, r *http.Request, _ httprouter.Params)
	Login(w http.ResponseWriter, r *http.Request, _ httprouter.Params)
}

type AuthControllerImpl struct {
	authService service.AuthService
}

func NewAuthController(authService service.AuthService) AuthController {
	return &AuthControllerImpl{authService: authService}
}

// @Summary Register a new user
// @Description Create a new user account
// @Tags users
// @Accept json
// @Produce json
// @Param user body dto.RegisterDataRequest true "User Data"
// @Success 200 {object} dto.RegisterDataResponse
// @Failure 400 {object} error
// @Router /users/register [post]
func (a *AuthControllerImpl) Register(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	registerRequest := dto.RegisterDataRequest{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&registerRequest); err != nil {
		log.Println(err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	data, err := a.authService.RegisterUser(ctx, registerRequest)
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to register new user", http.StatusBadRequest)
		return
	}

	resp := dto.Response{
		Code:   200,
		Status: "OK",
		Data:   data,
	}

	w.Header().Add("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(resp); err != nil {
		log.Println(err)
		http.Error(w, "Failed to encode response", http.StatusBadRequest)
		return
	}

}

// @Summary Login a user
// @Description Authenticate a user and return a token
// @Tags users
// @Accept json
// @Produce json
// @Param login body dto.LoginDataRequest true "Login Credentials"
// @Success 200 {object} dto.LoginDataResponse
// @Failure 400 {object} error
// @Router /users/login [post]
func (a *AuthControllerImpl) Login(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	loginRequest := dto.LoginDataRequest{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&loginRequest); err != nil {
		log.Println(err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	data, err := a.authService.CheckLogin(ctx, loginRequest)
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to verify login, try again later!", http.StatusBadRequest)
		return
	}

	resp := dto.Response{
		Code:   200,
		Status: "OK",
		Data:   data,
	}

	w.Header().Add("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	if err = encoder.Encode(resp); err != nil {
		log.Println(err)
		http.Error(w, "Failed to encode response", http.StatusBadRequest)
		return
	}
}
