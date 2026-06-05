package model

// GetUserInfoReq 查询用户信息请求
type GetUserInfoReq struct {
	UserId string `json:"uid" form:"uid" validate:"required"`
}

// GetUserInfoReply 查询用户信息返回
type GetUserInfoReply struct {
	UserId   string `json:"uid"`
	UserName string `json:"user_name"`
}
