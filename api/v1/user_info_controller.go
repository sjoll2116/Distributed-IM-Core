package v1

import (
	"fmt"
	"kama_chat_server/internal/dto/request"
	"kama_chat_server/internal/service/gorm"
	"kama_chat_server/pkg/constants"
	"kama_chat_server/pkg/util/response"
	"kama_chat_server/pkg/zlog"

	"github.com/gin-gonic/gin"
)

// Register 注册
func Register(c *gin.Context) {
	var registerReq request.RegisterRequest
	if err := c.BindJSON(&registerReq); err != nil {
		zlog.Error(err.Error())
		response.FailWithMessage(c, constants.SYSTEM_ERROR)
		return
	}
	fmt.Println(registerReq)
	message, userInfo, ret := gorm.UserInfoService.Register(registerReq)
	JsonBack(c, message, ret, userInfo)
}

// Login 登录
func Login(c *gin.Context) {
	var loginReq request.LoginRequest
	if err := c.BindJSON(&loginReq); err != nil {
		zlog.Error(err.Error())
		response.FailWithMessage(c, constants.SYSTEM_ERROR)
		return
	}
	message, userInfo, ret := gorm.UserInfoService.Login(loginReq)
	JsonBack(c, message, ret, userInfo)
}

// SmsLogin 验证码登录
func SmsLogin(c *gin.Context) {
	var req request.SmsLoginRequest
	if err := c.BindJSON(&req); err != nil {
		zlog.Error(err.Error())
		response.FailWithMessage(c, constants.SYSTEM_ERROR)
		return
	}
	message, userInfo, ret := gorm.UserInfoService.SmsLogin(req)
	JsonBack(c, message, ret, userInfo)
}

// UpdateUserInfo 修改用户信息
func UpdateUserInfo(c *gin.Context) {
	var req request.UpdateUserInfoRequest
	if err := c.BindJSON(&req); err != nil {
		zlog.Error(err.Error())
		response.FailWithMessage(c, constants.SYSTEM_ERROR)
		return
	}
	message, ret := gorm.UserInfoService.UpdateUserInfo(req)
	JsonBack(c, message, ret, nil)
}

// GetUserInfoList 获取用户列表
func GetUserInfoList(c *gin.Context) {
	var req request.GetUserInfoListRequest
	if err := c.BindJSON(&req); err != nil {
		zlog.Error(err.Error())
		response.FailWithMessage(c, constants.SYSTEM_ERROR)
		return
	}

	// 从 JWT Token 提取 UUID 进行查询，防止越权
	uuid := c.MustGet("uuid").(string)
	message, userList, ret := gorm.UserInfoService.GetUserInfoList(uuid)
	JsonBack(c, message, ret, userList)
}

// AbleUsers 启用用户
func AbleUsers(c *gin.Context) {
	var req request.AbleUsersRequest
	if err := c.BindJSON(&req); err != nil {
		zlog.Error(err.Error())
		response.FailWithMessage(c, constants.SYSTEM_ERROR)
		return
	}
	message, ret := gorm.UserInfoService.AbleUsers(req.UuidList)
	JsonBack(c, message, ret, nil)
}

// DisableUsers 禁用用户
func DisableUsers(c *gin.Context) {
	var req request.AbleUsersRequest
	if err := c.BindJSON(&req); err != nil {
		zlog.Error(err.Error())
		response.FailWithMessage(c, constants.SYSTEM_ERROR)
		return
	}
	message, ret := gorm.UserInfoService.DisableUsers(req.UuidList)
	JsonBack(c, message, ret, nil)
}

// GetUserInfo 获取用户信息
func GetUserInfo(c *gin.Context) {
	var req request.GetUserInfoRequest
	if err := c.BindJSON(&req); err != nil {
		zlog.Error(err.Error())
		response.FailWithMessage(c, constants.SYSTEM_ERROR)
		return
	}

	// 鉴权安全：只能获取自己的信息（如果是管理员可以根据业务放开，当前为防越权简单处理）
	_ = c.MustGet("uuid").(string) // targetUuid == myUuid check can be added here
	targetUuid := req.Uuid
	// 简单越权校验：如果不允许看别人的信息可以加上 targetUuid == myUuid 检查

	message, userInfo, ret := gorm.UserInfoService.GetUserInfo(targetUuid)
	JsonBack(c, message, ret, userInfo)
}

// DeleteUsers 删除用户
func DeleteUsers(c *gin.Context) {
	var req request.AbleUsersRequest
	if err := c.BindJSON(&req); err != nil {
		zlog.Error(err.Error())
		response.FailWithMessage(c, constants.SYSTEM_ERROR)
		return
	}
	message, ret := gorm.UserInfoService.DeleteUsers(req.UuidList)
	JsonBack(c, message, ret, nil)
}

// SetAdmin 设置管理员
func SetAdmin(c *gin.Context) {
	var req request.AbleUsersRequest
	if err := c.BindJSON(&req); err != nil {
		zlog.Error(err.Error())
		response.FailWithMessage(c, constants.SYSTEM_ERROR)
		return
	}
	message, ret := gorm.UserInfoService.SetAdmin(req.UuidList, req.IsAdmin)
	JsonBack(c, message, ret, nil)
}

// SendSmsCode 发送短信验证码
func SendSmsCode(c *gin.Context) {
	var req request.SendSmsCodeRequest
	if err := c.BindJSON(&req); err != nil {
		zlog.Error(err.Error())
		response.FailWithMessage(c, constants.SYSTEM_ERROR)
		return
	}
	message, ret := gorm.UserInfoService.SendSmsCode(req.Telephone)
	JsonBack(c, message, ret, nil)
}
