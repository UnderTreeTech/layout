package service

import (
	"context"

	"github.com/UnderTreeTech/waterdrop/pkg/log"
)

// initCrons 初始化并注册所有的定时任务
func (s *Service) initCrons() {
	// 示例：每分钟执行一次
	// s.cron.AddFunc("0 */1 * * * *", s.sampleCronTask)

	// 启动定时任务管理器
	s.cron.Start()
}

// sampleCronTask 示例定时任务
func (s *Service) sampleCronTask() {
	ctx := context.Background()
	log.Info(ctx, "sample cron task running")
}
