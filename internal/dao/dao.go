package dao

import (
	"context"
	"errors"
	"fmt"

	"github.com/UnderTreeTech/waterdrop/pkg/database/mongo"

	"github.com/UnderTreeTech/waterdrop/pkg/conf"
	"github.com/UnderTreeTech/waterdrop/pkg/database/redis"
	"github.com/UnderTreeTech/waterdrop/pkg/database/sql"
	"github.com/UnderTreeTech/waterdrop/pkg/log"
)

// interface Dao
type Dao interface {
	Close() error
	Ping(ctx context.Context) error

	GetCollection(name string) *mongo.Collection
}

// struct dao
type dao struct {
	db    *sql.DB
	redis *redis.Redis
	mongo *mongo.DB
}

// New return a dao that implements interface Dao
func New() Dao {
	db := NewMySQL()
	mongo := NewMongo()
	redis := NewRedis()
	return &dao{
		db:    db,
		redis: redis,
		mongo: mongo,
	}
}

// Close close backend base services
func (d *dao) Close() (err error) {
	d.db.Close()
	d.mongo.Close()
	d.redis.Close()
	return
}

// Ping ping backend base services, like db, redis, es etc.
func (d *dao) Ping(ctx context.Context) error {
	if err := d.db.Ping(ctx); err != nil {
		log.Error(ctx, "ping db fail", log.String("error", err.Error()))
		return err
	}

	if err := d.mongo.Ping(); err != nil {
		log.Error(ctx, "ping mongo fail", log.String("error", err.Error()))
		return err
	}

	if alive := d.redis.Ping(ctx); !alive {
		return errors.New("redis has gone")
	}

	return nil
}

func (d *dao) GetCollection(name string) *mongo.Collection {
	return d.mongo.GetCollection(name)
}

// attach transaction flag to context
type txKey struct{}

func (d *dao) Begin(ctx context.Context) (context.Context, error) {
	tx, err := d.db.Begin(ctx)
	if err != nil {
		return ctx, err
	}

	ctx = context.WithValue(ctx, txKey{}, tx)
	return ctx, err
}

func (d *dao) Commit(ctx context.Context) error {
	tx, err := d.GetTxFromCtx(ctx)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (d *dao) Rollback(ctx context.Context) error {
	tx, err := d.GetTxFromCtx(ctx)
	if err != nil {
		return err
	}

	return tx.Rollback()
}

func (d *dao) GetTxFromCtx(ctx context.Context) (*sql.Tx, error) {
	tx, ok := ctx.Value(txKey{}).(*sql.Tx)
	if !ok {
		return nil, errors.New("assert tx err")
	}

	return tx, nil
}

// NewMySQL returns mysql instance
func NewMySQL() *sql.DB {
	config := &sql.Config{}
	if err := conf.Unmarshal("mysql", config); err != nil {
		panic(fmt.Sprintf("unmarshal mysql config fail,err msg %s", err.Error()))
	}
	log.Infof("db config", log.Any("config", config))
	db := sql.NewMySQL(config)

	return db
}

// NewRedis returns redis instance
func NewRedis() *redis.Redis {
	config := &redis.Config{}
	if err := conf.Unmarshal("redis", config); err != nil {
		panic(fmt.Sprintf("unmarshal redis config fail,err msg %s", err.Error()))
	}
	log.Infof("redis config", log.Any("config", config))

	redis, err := redis.New(config)
	if err != nil {
		panic(fmt.Sprintf("new redis client fail,err msg %s", err.Error()))
	}
	return redis
}

// NewMongo return mongo instance
func NewMongo() *mongo.DB {
	cfg := &mongo.Config{}
	if err := conf.Unmarshal("mongo", cfg); err != nil {
		panic(fmt.Sprintf("unmarshal mongo config fail, err msg %s", err.Error()))
	}

	db := mongo.Open(cfg)
	return db
}
