package app

import (
	"go-chat/controller"

	"github.com/julienschmidt/httprouter"
)

func SetupRoutes(authController controller.AuthController, chatController controller.ChatController, userController controller.UserController) *httprouter.Router {

	router := httprouter.New()

	router.POST("/users/register", authController.Register)
	router.POST("/users/login", authController.Login)

	router.GET("/messages/:roomId", chatController.GetMessages)
	router.POST("/messages/chatRoom", chatController.GetorCreateChatRoom)
	
	router.POST("/friends/add", userController.AddFriend)
	router.GET("/friends/list/:userID", userController.GetFriendLists)
	router.GET("/friend-request/:userID", userController.GetFriendRequests)
	router.POST("/friend-request/respond", userController.UpdateFriendRequest)

	return router
}
