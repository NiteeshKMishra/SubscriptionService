package session

import (
	"encoding/gob"
	"os"
	"time"

	"github.com/NiteeshKMishra/SubscriptionService/cmd/database"
	"github.com/alexedwards/scs/redisstore"
	"github.com/alexedwards/scs/v2"
	"github.com/gomodule/redigo/redis"
)

func InitSession() *scs.SessionManager {
	gob.Register(database.User{})

	session := scs.New()
	session.Store = redisstore.New(initRedis())
	session.Lifetime = 24 * time.Hour

	return session
}

func initRedis() *redis.Pool {
	connString := os.Getenv("REDIS_HOST")
	redisPool := &redis.Pool{
		MaxIdle: 10,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", connString)
		},
	}

	return redisPool
}
