package server

import (
	"fmt"

	"github.com/fabianpoels/ai-bubble-tracker/cache"
	"github.com/fabianpoels/ai-bubble-tracker/config"
	"github.com/fabianpoels/ai-bubble-tracker/db"
)

func Init() {
	config := config.GetConfig()
	r := NewRouter()
	db.DbConnect()
	cache.CacheConnect()
	r.Run(fmt.Sprintf("%s:%s", config.GetString("server.host"), config.GetString("server.port")))
}
