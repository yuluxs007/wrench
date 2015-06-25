package db

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"gopkg.in/redis.v3"
	"strings"
	"time"
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

func GetUUID(ObjectType, Object string) (UUID string, err error) {

	index := ""

	switch strings.TrimSpace(ObjectType) {

	case "user":
		index = GLOBAL_USER_INDEX
	case "repository":
		index = GLOBAL_REPOSITORY_INDEX
	case "organization":
		index = GLOBAL_ORGANIZATION_INDEX
	case "team":
		index = GLOBAL_TEAM_INDEX
	case "image":
		index = GLOBAL_IMAGE_INDEX
	case "tarsum":
		index = GLOBAL_TARSUM_INDEX
	case "tag":
		index = GLOBAL_TAG_INDEX
	case "compose":
		index = GLOBAL_COMPOSE_INDEX
	case "admin":
		index = GLOBAL_ADMIN_INDEX
	case "log":
		index = GLOBAL_LOG_INDEX
	default:

	}

	if UUID, err = Client.HGet(index, Object).Result(); err != nil {
		return "", err
	} else {
		return UUID, nil
	}
}

func Save(obj interface{}, key string) (err error) {

	result, err := json.Marshal(&obj)
	if err != nil {
		return err
	}

	if _, err := Client.Set(key, string(result), 0).Result(); err != nil {
		return err
	}

	return nil
}

func Get(obj interface{}, key string) (err error) {
	result, err := Client.Get(key).Result()
	if err != nil {
		return err
	}

	if err = json.Unmarshal([]byte(result), &obj); err != nil {
		return err
	}

	return nil
}
