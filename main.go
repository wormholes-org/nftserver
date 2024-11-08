package main

import (
	"fmt"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/nftexchange/nftserver/common/sync"
	"github.com/nftexchange/nftserver/models"
	_ "github.com/nftexchange/nftserver/routers"
	"log"
	"net/http"
	_ "net/http/pprof"
)

// @title           Swagger Nftserver API
// @version         2.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/

// @host      localhost:9051
// @BasePath  /api/v2

// @securityDefinitions.basic  BasicAuth
func main() {
	//8c995fd78bddf528bd548cce025f62d4c3c0658362dbfd31b23414cf7ce2e8ed
	//verify := signature.VerifyAppconf("./conf/app.conf", "0x2b0aD05ADDa21BA4E5b94C4f9aE3BCeA15A380c5")
	//if verify != true {
	//	fmt.Println("app.conf verify error.")
	//	return
	//}
	/*err :=  os.Remove("./conf/app.conf")
	if err != nil {
		fmt.Println("delete app.conf err=", err)
		return
	}*/
	if models.DebugPort != "" {
		go func() {
			log.Println(http.ListenAndServe("0.0.0.0:"+models.DebugPort, nil))
		}()
	}

	err := models.InitSysParams(models.Sqldsndb)
	if err != nil {
		fmt.Println("InitSysParams err=", err)
		return
	}
	fmt.Println(models.NFTUploadAuditRequired)
	//err = models.InitSyncBlockTs(models.Sqldsndb)
	err = sync.InitSyncBlockTs(models.Sqldsndb)
	if err != nil {
		fmt.Println("init err exit")
		return
	}
	if !models.LimitWritesDatabase {
		go TimeProc(models.Sqldsndb)
	}
	//beego.BConfig.MaxMemory = nftexchangev2.UpLoadSize
	beego.Run()
}
