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
)

const ctxKey = "bark_serverless_token_db"

type (
	Db interface {
		Name() string
	}
	LoadToken interface {
		LoadToken(key string) (string, error)
	}
	SaveToken interface {
		SaveToken(key, token string) error
	}
)

// 如果没有db对象会直接返回nil
func readDBFromCtx(c *gin.Context) []Db {
	value, isExist := c.Get(ctxKey)
	if !isExist {
		return nil
	}
	dbs, isExist := value.([]Db)
	if !isExist {
		return nil
	}
	return dbs
}

func writeDbToCtx(dbs ...Db) gin.HandlerFunc {
	tmp := make([]Db, 0, len(dbs))
	for _, db := range dbs {
		if db != nil {
			tmp = append(tmp, db)
		}
	}
	return func(c *gin.Context) {
		if len(tmp) > 0 {
			c.Set(ctxKey, tmp)
		}
	}
}

type envDB struct{} // 从 Env 环境读取 Token - 不支持写入 token

var _ interface {
	Db
	LoadToken
} = (*envDB)(nil)

func (*envDB) Name() string {
	return "env"
}

func (*envDB) LoadToken(key string) (string, error) {
	// 环境变量中设备Key的前缀
	const deviceKeyPrefix = "device_"

	token, isExist := os.LookupEnv(deviceKeyPrefix + key)
	if !isExist {
		return "", errors.New("failed to get token from env")
	}
	return token, nil
}

type sqlDB struct {
	name string
	*gorm.DB
} // 从 gorm.DB 对象读取 Token

// 基于 GORM 支持 MysQL
func newSQLDB() *sqlDB {
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
		case "postgresql":
			name = "postgresql"
			dialector = postgres.Open(dsn)
		case "sqlserver":
			name = "sqlserver"
			dialector = sqlserver.Open(dsn)
		case "clickhouse":
			name = "clickhouse"
			dialector = clickhouse.Open(dsn)
		default:
			name = "mysql"
			dialector = mysql.Open(dsn)
		}
	}

	if dialector == nil || name == "" {
		return nil
	}

	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		return nil
	}

	// 迁移数据库 - 迁移失败不会影响数据库实例初始化
	_ = db.AutoMigrate(new(Token))

	zap.S().Infof("load sql db successfully, db: %s, dialector: %s", name, dialector.Name())
	return &sqlDB{
		name: name,
		DB:   db,
	}
}

var _ interface {
	Db
	LoadToken
	SaveToken
} = (*sqlDB)(nil)

func (s *sqlDB) Name() string { return s.name }

func (s *sqlDB) LoadToken(key string) (string, error) {
	if key == "" {
		return "", errors.New("key is empty")
	}
	var t = &Token{Key: key}
	if err := s.DB.First(t).Error; err != nil {
		return "", err
	}
	return t.Token, nil
}

func (s *sqlDB) SaveToken(key, token string) error {
	if key == "" {
		return errors.New("key is empty")
	}
	if token == "" {
		return errors.New("token is empty")
	}
	t := &Token{Key: key}
	s.DB.First(t)
	if err := s.DB.Save(&Token{Model: gorm.Model{ID: t.ID}, Key: key, Token: token}).Error; err != nil {
		return err
	}
	return nil
}
