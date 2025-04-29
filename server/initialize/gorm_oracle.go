package initialize

import (
	//"github.com/dzwvip/oracle"
	// "github.com/godror/godror"
	"github.com/flipped-aurora/gin-vue-admin/server/config"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/initialize/internal"

	oracle "github.com/seelly/gorm-oracle"
	"gorm.io/gorm"
)

// GormOracle 初始化oracle数据库
// 如果需要Oracle库 放开import里的注释 把下方 mysql.Config 改为 oracle.Config ;  mysql.New 改为 oracle.New
func GormOracle() (*gorm.DB, func(string) gorm.Dialector, config.GeneralDB) {
	o := global.GVA_CONFIG.Oracle
	return GormOracleByConfig(o), oracle.Open, o.GeneralDB
}

// GormOracleByConfig 初始化Oracle数据库用过传入配置
func GormOracleByConfig(o config.Oracle) *gorm.DB {
	return initOracleDatabase(o)
}

// initOracleDatabase 初始化Oracle数据库的辅助函数
func initOracleDatabase(o config.Oracle) *gorm.DB {
	if o.Dbname == "" {
		return nil
	}

	oracleConfig := oracle.Config{
		DSN:               o.Dsn(), // DSN data source name
		DefaultStringSize: 191,     // string 类型字段的默认长度
	}

	if db, err := gorm.Open(oracle.New(oracleConfig), internal.Gorm.Config(o.Prefix, o.Singular)); err != nil {
		panic(err)
	} else {
		sqlDB, _ := db.DB()
		sqlDB.SetMaxIdleConns(o.MaxIdleConns)
		sqlDB.SetMaxOpenConns(o.MaxOpenConns)
		return db
	}
}
