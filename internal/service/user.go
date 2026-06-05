package service

import (
	"context"

	m "github.com/UnderTreeTech/layout/internal/server/http/model"
)

// GetUserInfo 获取用户信息
func (s *Service) GetUserInfo(ctx context.Context, uid string) (reply *m.GetUserInfoReply, err error) {
	return &m.GetUserInfoReply{
		UserId:   "1",
		UserName: "johnsun",
	}, nil
}
