package controller

import (
	"encoding/json"
	"go-chat/dto"
	"go-chat/service"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type UserController interface {
	AddFriend(w http.ResponseWriter, r *http.Request, param httprouter.Params)
	UpdateFriendRequest(w http.ResponseWriter, r *http.Request, param httprouter.Params)
	GetFriendRequests(w http.ResponseWriter, r *http.Request, param httprouter.Params)
	GetFriendLists(w http.ResponseWriter, r *http.Request, param httprouter.Params)
}

type UserControllerImpl struct {
	userService service.UserService
}

func NewUserController(userService service.UserService) UserController {
	return &UserControllerImpl{userService: userService}
}

// @Summary Add a friend
// @Description Send a friend request to another user
// @Tags friends
// @Accept json
// @Produce json
// @Param friend body dto.FriendRequestParameter true "Friend Request Data"
// @Success 200 {object} dto.FriendRequestResponse
// @Failure 400 {object} error
// @Router /friends/add [post]
func (u *UserControllerImpl) AddFriend(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	friendRequest := dto.FriendRequestParameter{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&friendRequest); err != nil {
		log.Println(err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	data, err := u.userService.AddFriend(ctx, friendRequest)
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to add new friend", http.StatusBadRequest)
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

// @Summary Respond to a friend request
// @Description Accept or reject a friend request
// @Tags friends
// @Accept json
// @Produce json
// @Param response body dto.UpdateRequestParameter true "Friend Request Response"
// @Success 200 {object} dto.UpdateFriendRequestResponse
// @Failure 400 {object} error
// @Router /friend-request/respond [post]
func (u *UserControllerImpl) UpdateFriendRequest(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	updateRequest := dto.UpdateRequestParameter{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&updateRequest); err != nil {
		log.Println(err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	data, err := u.userService.UpdateFriendRequest(ctx, updateRequest)
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to update friend request", http.StatusBadRequest)
		return
	}

	resp := dto.Response{
		Code:   200,
		Status: "OK",
		Data:   data,
	}

	w.Header().Add("Content_Type", "application/json")
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(resp); err != nil {
		log.Println(err)
		http.Error(w, "Failed to encode response", http.StatusBadRequest)
		return
	}
}

// @Summary Get friend requests
// @Description Retrieve a list of friend requests for a specific user
// @Tags friends
// @Accept json
// @Produce json
// @Param userID path string true "User ID"
// @Success 200 {object} dto.FriendRequestResponse
// @Failure 400 {object} error
// @Router /friend-request/{userID} [get]
func (u *UserControllerImpl) GetFriendRequests(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	userID := param.ByName("userID")

	ctx := r.Context()

	data, err := u.userService.GetFriendRequests(ctx, userID)
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to get friend requests", http.StatusBadRequest)
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

// @Summary Get friend lists
// @Description Retrieve a list of friends for a specific user
// @Tags friends
// @Accept json
// @Produce json
// @Param userID path string true "User ID"
// @Success 200 {object} dto.GetFriendListsResponse
// @Failure 400 {object} error
// @Router /friends/list/{userID} [get]
func (u *UserControllerImpl) GetFriendLists(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	userID := param.ByName("userID")

	ctx := r.Context()

	data, err := u.userService.GetFriendLists(ctx, userID)
	if err != nil {
		log.Println(err)
		http.Error(w, "Failed to get friend lists", http.StatusBadRequest)
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
