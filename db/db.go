package db

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"gopkg.in/redis.v3"
)

const (
	//Dockyard Data Index
	GLOBAL_REPOSITORY_INDEX = "GLOBAL_REPOSITORY_INDEX"
	GLOBAL_IMAGE_INDEX      = "GLOBAL_IMAGE_INDEX"
	GLOBAL_TARSUM_INDEX     = "GLOBAL_TARSUM_INDEX"
	GLOBAL_TAG_INDEX        = "GLOBAL_TAG_INDEX"
	GLOBAL_COMPOSE_INDEX    = "GLOBAL_COMPOSE_INDEX"
	//Sail Data Index
	GLOBAL_USER_INDEX         = "GLOBAL_USER_INDEX"
	GLOBAL_ORGANIZATION_INDEX = "GLOBAL_ORGANIZATION_INDEX"
	GLOBAL_TEAM_INDEX         = "GLOBAL_TEAM_INDEX"
	GLOBAL_PRIVILEGE_INDEX    = "GLOBAL_PRIVILEGE_INDEX"

	GLOBAL_ADMIN_INDEX = "GLOBAL_ADMIN_INDEX"
	GLOBAL_LOG_INDEX   = "GLOBAL_LOG_INDEX"
)

var (
	Client *redis.Client
)

func InitDB(addr, passwd string, db int64) error {
	Client = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: passwd,
		DB:       db,
	})

	if _, err := Client.Ping().Result(); err != nil {
		return err
	} else {
		return nil
	}
}

func GeneralDBKey(key string) string {
	md5String := fmt.Sprintf("%s%d", key, time.Now().Unix())
	h := md5.New()
	h.Write([]byte(md5String))
	return hex.EncodeToString(h.Sum(nil))
}

func SaveObject(object interface{}, key string) error {
	if json, err := json.Marshal(object); err != nil {
		return err
	} else {
		if err := Client.Set(string(json), key, 0).Err(); err != nil {
			return err
		} else {
			return nil
		}
	}
}

func LoadObject(object interface{}, key string) error {
	if value, err := Client.Get(key).Result(); err != nil {
		return err
	} else {
		if err := json.Unmarshal([]byte(value), object); err != nil {
			return err
		} else {
			return nil
		}
	}
}
