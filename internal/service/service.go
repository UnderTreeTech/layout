package service

import (
	"github.com/UnderTreeTech/layout/internal/dao"
	"github.com/robfig/cron/v3"
)

type Config struct {
	MaxIdPoolNum int
}

type Service struct {
	cfg  *Config
	dao  dao.Dao
	cron *cron.Cron
}

func New(d dao.Dao, cfg *Config) *Service {
	// 初始化定时任务管理器
	chains := []cron.JobWrapper{cron.Recover(cron.DefaultLogger), cron.SkipIfStillRunning(cron.DefaultLogger)}
	cronManager := cron.New(cron.WithSeconds(), cron.WithChain(chains...))

	svc := &Service{
		dao:  d,
		cfg:  cfg,
		cron: cronManager,
	}

	svc.initCrons()
	return svc
}

func (s *Service) Close() {
	if s.cron != nil {
		s.cron.Stop()
	}
	s.dao.Close()
}

func (s *Service) Ping() {

}
