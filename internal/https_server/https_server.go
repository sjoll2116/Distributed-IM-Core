package https_server

import (
	v1 "kama_chat_server/api/v1"
	"kama_chat_server/internal/config"
	"kama_chat_server/internal/middleware"
	"kama_chat_server/pkg/ssl"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var GE *gin.Engine

func init() {
	GE = gin.Default()
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"*"}
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	GE.Use(cors.New(corsConfig))

	// 添加全局异常捕获中间件，防止服务器因为未捕获的 Panic 崩溃
	GE.Use(middleware.GlobalRecovery())

	GE.Use(ssl.TlsHandler(config.GetConfig().MainConfig.Host, config.GetConfig().MainConfig.Port))
	GE.Static("/static/avatars", config.GetConfig().StaticAvatarPath)
	GE.Static("/static/files", config.GetConfig().StaticFilePath)

	// 公开路由 (无需 JWT 保护)
	publicGroup := GE.Group("/")
	{
		publicGroup.POST("/login", v1.Login)
		publicGroup.POST("/register", v1.Register)
		publicGroup.POST("/user/smsLogin", v1.SmsLogin)
		publicGroup.POST("/user/sendSmsCode", v1.SendSmsCode)
	}

	// 保护路由 (需要 JWT Token 鉴权)
	authGroup := GE.Group("/")
	authGroup.Use(middleware.JWTAuthMiddleware())
	{
		authGroup.POST("/user/updateUserInfo", v1.UpdateUserInfo)
		authGroup.POST("/user/getUserInfoList", v1.GetUserInfoList)
		authGroup.POST("/user/ableUsers", v1.AbleUsers)
		authGroup.POST("/user/getUserInfo", v1.GetUserInfo)
		authGroup.POST("/user/disableUsers", v1.DisableUsers)
		authGroup.POST("/user/deleteUsers", v1.DeleteUsers)
		authGroup.POST("/user/setAdmin", v1.SetAdmin)
		authGroup.POST("/user/wsLogout", v1.WsLogout)

		authGroup.POST("/group/createGroup", v1.CreateGroup)
		authGroup.POST("/group/loadMyGroup", v1.LoadMyGroup)
		authGroup.POST("/group/checkGroupAddMode", v1.CheckGroupAddMode)
		authGroup.POST("/group/enterGroupDirectly", v1.EnterGroupDirectly)
		authGroup.POST("/group/leaveGroup", v1.LeaveGroup)
		authGroup.POST("/group/dismissGroup", v1.DismissGroup)
		authGroup.POST("/group/getGroupInfo", v1.GetGroupInfo)
		authGroup.POST("/group/getGroupInfoList", v1.GetGroupInfoList)
		authGroup.POST("/group/deleteGroups", v1.DeleteGroups)
		authGroup.POST("/group/setGroupsStatus", v1.SetGroupsStatus)
		authGroup.POST("/group/updateGroupInfo", v1.UpdateGroupInfo)
		authGroup.POST("/group/getGroupMemberList", v1.GetGroupMemberList)
		authGroup.POST("/group/removeGroupMembers", v1.RemoveGroupMembers)

		authGroup.POST("/session/openSession", v1.OpenSession)
		authGroup.POST("/session/getUserSessionList", v1.GetUserSessionList)
		authGroup.POST("/session/getGroupSessionList", v1.GetGroupSessionList)
		authGroup.POST("/session/deleteSession", v1.DeleteSession)
		authGroup.POST("/session/checkOpenSessionAllowed", v1.CheckOpenSessionAllowed)

		authGroup.POST("/contact/getUserList", v1.GetUserList)
		authGroup.POST("/contact/loadMyJoinedGroup", v1.LoadMyJoinedGroup)
		authGroup.POST("/contact/getContactInfo", v1.GetContactInfo)
		authGroup.POST("/contact/deleteContact", v1.DeleteContact)
		authGroup.POST("/contact/applyContact", v1.ApplyContact)
		authGroup.POST("/contact/getNewContactList", v1.GetNewContactList)
		authGroup.POST("/contact/passContactApply", v1.PassContactApply)
		authGroup.POST("/contact/blackContact", v1.BlackContact)
		authGroup.POST("/contact/cancelBlackContact", v1.CancelBlackContact)
		authGroup.POST("/contact/getAddGroupList", v1.GetAddGroupList)
		authGroup.POST("/contact/refuseContactApply", v1.RefuseContactApply)
		authGroup.POST("/contact/blackApply", v1.BlackApply)

		authGroup.POST("/message/getMessageList", v1.GetMessageList)
		authGroup.POST("/message/getGroupMessageList", v1.GetGroupMessageList)
		authGroup.POST("/message/uploadAvatar", v1.UploadAvatar)
		authGroup.POST("/message/uploadFile", v1.UploadFile)

		authGroup.POST("/chatroom/getCurContactListInChatRoom", v1.GetCurContactListInChatRoom)
	}

	GE.GET("/wss", v1.WsLogin)

}
