package web

import (
	"net"
	"net/http"
	"os"
	"sort"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/location/v2"
	"github.com/gin-gonic/gin"

	"server/log"
	"server/settings"
	"server/torr"
	"server/version"
	"server/web/api"
	"server/web/auth"
	"server/web/blocker"
)

var (
	BTS        = torr.NewBTS()
	waitChan   = make(chan error)
	httpServer *http.Server
)

func Start() {
	log.TLogln("Start TorrServer " + version.Version + " torrent " + version.GetTorrentVersion())
	ips := GetLocalIps()
	if len(ips) > 0 {
		log.TLogln("Local IPs:", ips)
	}
	err := BTS.Connect()
	if err != nil {
		log.TLogln("BTS.Connect() error!", err)
		os.Exit(1)
	}

	gin.SetMode(gin.ReleaseMode)

	corsCfg := cors.DefaultConfig()
	corsCfg.AllowAllOrigins = true
	corsCfg.AllowPrivateNetwork = true
	corsCfg.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "X-Requested-With", "Accept", "Authorization"}

	route := gin.New()
	route.Use(log.WebLogger(), blocker.Blocker(), gin.Recovery(), cors.New(corsCfg), location.Default())
	auth.SetupAuth(route)

	route.GET("/echo", echo)
	api.SetupRoute(route)

	httpServer = &http.Server{
		Addr:    ":" + settings.Port,
		Handler: route,
	}

	go func() {
		log.TLogln("Start http server at port", settings.Port)
		waitChan <- httpServer.ListenAndServe()
	}()
}

func Wait() error {
	return <-waitChan
}

func Stop() {
	if httpServer != nil {
		httpServer.Close()
	}
	BTS.Disconnect()
	waitChan <- nil
}

// echo godoc
//
//	@Summary		Tests server status
//	@Description	Tests whether server is alive or not
//
//	@Tags			API
//
//	@Produce		plain
//	@Success		200	{string}	string	"Server version"
//	@Router			/echo [get]
func echo(c *gin.Context) {
	c.String(200, "%v", version.Version)
}

func GetLocalIps() []string {
	ifaces, err := net.Interfaces()
	if err != nil {
		log.TLogln("Error get local IPs")
		return nil
	}
	var list []string
	for _, i := range ifaces {
		addrs, _ := i.Addrs()
		if i.Flags&net.FlagUp == net.FlagUp {
			for _, addr := range addrs {
				var ip net.IP
				switch v := addr.(type) {
				case *net.IPNet:
					ip = v.IP
				case *net.IPAddr:
					ip = v.IP
				}
				if !ip.IsLoopback() && !ip.IsLinkLocalUnicast() && !ip.IsLinkLocalMulticast() {
					list = append(list, ip.String())
				}
			}
		}
	}
	sort.Strings(list)
	return list
}

