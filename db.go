package main

import (
	"errors"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/driver/clickhouse"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

const CtxDBKey = "bark_serverless_token_db"

type (
	TokenStore interface {
		Name() string
	}
	LoadToken interface {
		LoadToken(key string) (string, error)
	}
	SaveToken interface {
		SaveToken(key, token string) error
	}
)

// ReadTokenStoreFromCtx 如果没有Store对象会直接返回nil
func ReadTokenStoreFromCtx(c *gin.Context) []TokenStore {
	value, isExist := c.Get(CtxDBKey)
	if !isExist {
		return nil
	}
	dbs, isExist := value.([]TokenStore)
	if !isExist {
		return nil
	}
	return dbs
}

// WriteTokenStoreToCtx 上下文中写入Store对象
func WriteTokenStoreToCtx(dbs ...TokenStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(CtxDBKey, dbs)
	}
}

type env struct{} // 从 Env 环境读取 Token - 不支持写入 token

var _ interface {
	TokenStore
	LoadToken
} = (*env)(nil)

func (*env) Name() string {
	return "env"
}

func (*env) LoadToken(key string) (string, error) {
	// 环境变量中设备Key的前缀
	const deviceKeyPrefix = "device_"

	token, isExist := os.LookupEnv(deviceKeyPrefix + key)
	if !isExist {
		return "", errors.New("failed to get token from env")
	}
	return token, nil
}

type gormDB struct {
	name string
	*gorm.DB
} // 从 gorm.DB 对象读取 Token

// 基于 GORM 支持 MysQL
func newDB() *gormDB {
	var (
		name      string
		dialector gorm.Dialector
	)

	for _, v := range os.Environ() {
		keyAndValue := strings.SplitN(v, "=", 2)
		key, dsn := keyAndValue[0], keyAndValue[1]
		if !strings.HasPrefix(key, "dsn") {
			continue
		}
		switch strings.TrimPrefix(key, "dsn") {
		case "postgresql": // eg. dsn_postgresql
			name = "postgresql"
			dialector = postgres.Open(dsn)
		case "sqlserver": // eg. dsn_sqlserver
			name = "sqlserver"
			dialector = sqlserver.Open(dsn)
		case "clickhouse": // eg. dsn_clickhouse
			name = "clickhouse"
			dialector = clickhouse.Open(dsn)
		default: // eg. dsn_*
			name = "mysql"
			dialector = mysql.Open(dsn)
		}
	}

	if dialector == nil || name == "" {
		return nil
	}

	db, err := gorm.Open(dialector, &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		zap.L().Error("connect to db failed", zap.Error(err))
		return nil
	}

	// 迁移数据库 - 迁移失败不会影响数据库实例初始化
	_ = db.AutoMigrate(new(Token))

	zap.L().With(zap.String("db", name), zap.String("dialector", dialector.Name())).
		Info("load db successfully")
	return &gormDB{
		name: name,
		DB:   db,
	}
}

var _ interface {
	TokenStore
	LoadToken
	SaveToken
} = (*gormDB)(nil)

func (s *gormDB) Name() string { return s.name }

func (s *gormDB) LoadToken(key string) (string, error) {
	if key == "" {
		return "", errors.New("key is empty")
	}
	var t = &Token{Key: key}
	if err := s.DB.First(t).Error; err != nil {
		return "", err
	}
	return t.Token, nil
}

func (s *gormDB) SaveToken(key, token string) error {
	if key == "" {
		return errors.New("key is empty")
	}
	if token == "" {
		return errors.New("token is empty")
	}
	t := &Token{Key: key}

	if errors.Is(s.DB.Where(t).First(t).Error, gorm.ErrRecordNotFound) {
		t.Token = token
		db := s.DB.Create(t)
		if err := db.Error; err != nil {
			zap.L().Error("create token failed", zap.Error(err))
			return err
		}
		zap.L().Info("create token success", zap.Uint("id", t.ID))
		return nil
	}

	t.Token = token
	if err := s.DB.Updates(t).Error; err != nil {
		return err
	}
	return nil
}
