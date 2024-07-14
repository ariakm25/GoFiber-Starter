package internal_casbin

import (
	internal_log "GoFiber-API/internal/log"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
)

var CasbinAdapter *gormadapter.Adapter
var CasbinEnforcer *casbin.Enforcer

func InitAdapter(pathModel string, host string, port string, user string, password string, db_name string, ssl_mode string) {
	build_dsn := "host=" + host + " port=" + port + " user=" + user + " password=" + password + " dbname=" + db_name + " sslmode=" + ssl_mode
	adapter, err := gormadapter.NewAdapter("postgres", build_dsn, true)

	if err != nil {
		internal_log.Logger.Sugar().Errorf("Error creating Casbin adapter with error: %s", err.Error())
		panic(err)
	}

	CasbinAdapter = adapter

	enforcer, err := casbin.NewEnforcer(pathModel, CasbinAdapter)

	if err != nil {
		internal_log.Logger.Sugar().Errorf("Error creating Casbin enforcer with error: %s", err.Error())
		panic(err)
	}

	CasbinEnforcer = enforcer

	err = CasbinEnforcer.LoadPolicy()

	if err != nil {
		internal_log.Logger.Sugar().Errorf("Error loading Casbin policy with error: %s", err.Error())
		panic(err)
	}
}
