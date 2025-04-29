package initialize

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/config"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/example"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"

	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

var policies = map[string]dbresolver.Policy{
	"random": dbresolver.RandomPolicy{},
}

func Gorm() (db *gorm.DB) {
	var open func(string) gorm.Dialector
	var s config.GeneralDB
	var dsn string
	switch global.GVA_CONFIG.System.DbType {
	case "mssql":
		global.GVA_ACTIVE_DBNAME = &global.GVA_CONFIG.Mssql.Dbname
		db, open, s = GormMssql()
	case "mysql":
		global.GVA_ACTIVE_DBNAME = &global.GVA_CONFIG.Mysql.Dbname
		db, open, s = GormMysql()
	case "pgsql":
		global.GVA_ACTIVE_DBNAME = &global.GVA_CONFIG.Pgsql.Dbname
		db, open, s = GormPgSql()
	case "oracle":
		global.GVA_ACTIVE_DBNAME = &global.GVA_CONFIG.Oracle.Dbname
		db, open, s = GormOracle()
	case "sqlite":
		global.GVA_ACTIVE_DBNAME = &global.GVA_CONFIG.Sqlite.Dbname
		db, open, s = GormSqlite()
	default:
		global.GVA_ACTIVE_DBNAME = &global.GVA_CONFIG.Mysql.Dbname
		db, open, s = GormMysql()
	}

	// 多个数据库支持
	c := global.GVA_CONFIG.Registers
	if db != nil && len(c) > 0 {
		registers := make([]config.Registers, len(c))
		for i := range c {
			sources := make([]string, len(c[i].Sources))

			if global.GVA_CONFIG.System.DbType == "sqlite" {
				dsn = filepath.Join(s.Path, s.Dbname+".db")
				for idx, source := range c[i].Sources {
					sources[idx] = filepath.Join(s.Path, time.Now().Format(source))
					if _, err := os.Stat(sources[idx]); os.IsNotExist(err) {
						if _, err := os.Stat(dsn); err == nil {
							if source, err := os.Open(dsn); err == nil {
								if destination, err := os.Create(sources[idx]); err == nil {
									defer destination.Close()
									if _, err := io.Copy(destination, source); err != nil {
										fmt.Println("初始化分库失败")
									}
								}
							}
						}
					}
				}
			}

			registers[i] = config.Registers{
				Sources:  sources,
				Replicas: c[i].Replicas,
				Policy:   c[i].Policy,
				Tables:   c[i].Tables,
			}
		}
		var register *dbresolver.DBResolver
		for _, e := range registers {
			if len(e.Tables) != 0 || len(e.Sources) != 0 || len(e.Replicas) != 0 {
				var config dbresolver.Config
				if len(e.Sources) > 0 {
					config.Sources = make([]gorm.Dialector, len(e.Sources))
					for i := range e.Sources {
						config.Sources[i] = open(e.Sources[i])
					}
				}
				if len(e.Replicas) > 0 {
					config.Replicas = make([]gorm.Dialector, len(e.Replicas))
					for i := range e.Replicas {
						config.Replicas[i] = open(e.Replicas[i])
					}
				}
				if e.Policy != "" {
					policy, ok := policies[e.Policy]
					if ok {
						config.Policy = policy
					}
				}
				tables := make([]any, len(e.Tables))
				for i, table := range e.Tables {
					tables[i] = table
				}
				if register == nil {
					register = dbresolver.Register(config, tables...)
				} else {
					register = register.Register(config, tables...)
				}
			}
		}
		if register == nil {
			register = dbresolver.Register(dbresolver.Config{})
		}
		if s.ConnMaxIdleTime > 0 {
			register = register.SetConnMaxIdleTime(time.Duration(s.ConnMaxIdleTime) * time.Second)
		}
		if s.ConnMaxLifetime > 0 {
			register = register.SetConnMaxLifetime(time.Duration(s.ConnMaxLifetime) * time.Second)
		}
		if s.MaxOpenConns > 0 {
			register = register.SetMaxOpenConns(s.MaxOpenConns)
		}
		if s.MaxIdleConns > 0 {
			register = register.SetMaxIdleConns(s.MaxIdleConns)
		}
		if register != nil {
			_ = db.Use(register)
		}
	}

	return db
}

func RegisterTables() {
	db := global.GVA_DB
	err := db.AutoMigrate(

		system.SysApi{},
		system.SysIgnoreApi{},
		system.SysUser{},
		system.SysBaseMenu{},
		system.JwtBlacklist{},
		system.SysAuthority{},
		system.SysDictionary{},
		system.SysOperationRecord{},
		system.SysAutoCodeHistory{},
		system.SysDictionaryDetail{},
		system.SysBaseMenuParameter{},
		system.SysBaseMenuBtn{},
		system.SysAuthorityBtn{},
		system.SysAutoCodePackage{},
		system.SysExportTemplate{},
		system.Condition{},
		system.JoinTemplate{},
		system.SysParams{},

		example.ExaFile{},
		example.ExaCustomer{},
		example.ExaFileChunk{},
		example.ExaFileUploadAndDownload{},
		example.ExaAttachmentCategory{},
	)
	if err != nil {
		global.GVA_LOG.Error("register table failed", zap.Error(err))
		os.Exit(0)
	}

	err = bizModel()

	if err != nil {
		global.GVA_LOG.Error("register biz_table failed", zap.Error(err))
		os.Exit(0)
	}
	global.GVA_LOG.Info("register table success")
}
