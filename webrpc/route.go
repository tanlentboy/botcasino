package webrpc

import (
	"net/http"
	"sync"

	"github.com/zhangpanyi/botcasino/webrpc/handlers"
)

var once sync.Once

// InitRoute 初始路由
func InitRoute(router *http.ServeMux) {
	once.Do(func() {
		router.HandleFunc("/webrpc/addad", handlers.AddAd)
		router.HandleFunc("/webrpc/delad", handlers.DelAd)
		router.HandleFunc("/webrpc/getads", handlers.GetAds)
		router.HandleFunc("/webrpc/addasset", handlers.AddAsset)
		router.HandleFunc("/webrpc/backup", handlers.Backup)
		router.HandleFunc("/webrpc/broadcast", handlers.Broadcast)
		router.HandleFunc("/webrpc/deductasset", handlers.DeductAsset)
		router.HandleFunc("/webrpc/frozen", handlers.Frozen)
		router.HandleFunc("/webrpc/unfrozen", handlers.Unfrozen)
		router.HandleFunc("/webrpc/get_assets", handlers.GetAssets)
		router.HandleFunc("/webrpc/restore", handlers.RestoreOrder)
		router.HandleFunc("/webrpc/subscribers", handlers.GetSubscribers)
	})
}
