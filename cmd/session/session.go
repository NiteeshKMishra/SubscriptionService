package session

import (
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/redisstore"
	"github.com/alexedwards/scs/v2"
	"github.com/gomodule/redigo/redis"
)

func InitSession() *scs.SessionManager {
	session := scs.New()
	session.Store = redisstore.New(initRedis())
	session.Lifetime = 24 * time.Hour
	session.Cookie = scs.SessionCookie{
		Persist:  true,
		SameSite: http.SameSiteLaxMode,
		Secure:   true,
	}

	return session
}

func initRedis() *redis.Pool {
	connString := os.Getenv("REDIS")
	redisPool := &redis.Pool{
		MaxIdle: 10,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", connString)
		},
	}

	return redisPool
}
