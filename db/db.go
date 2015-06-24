package db

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/containerops/wrench/utils"
	"gopkg.in/redis.v3"
	"reflect"
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
	s := reflect.TypeOf(obj).Elem()

	for i := 0; i < s.NumField(); i++ {

		value := reflect.ValueOf(obj).Elem().Field(s.Field(i).Index[0])

		switch value.Kind() {

		case reflect.String:
			if _, err := Client.HSet(key, s.Field(i).Name, value.String()).Result(); err != nil {
				return err
			}

		case reflect.Bool:
			if _, err := Client.HSet(key, s.Field(i).Name, utils.BoolToString(value.Bool())).Result(); err != nil {
				return err
			}

		case reflect.Int64:
			if _, err := Client.HSet(key, s.Field(i).Name, utils.Int64ToString(value.Int())).Result(); err != nil {
				return err
			}

		case reflect.Slice:
			if "[]string" == value.Type().String() && !value.IsNil() {
				strJson := "["

				for i := 0; i < value.Len(); i++ {
					nowUUID := value.Index(i).String()
					if i != 0 {
						strJson += fmt.Sprintf(`,"%s"`, nowUUID)
					} else {
						strJson += fmt.Sprintf(`"%s"`, nowUUID)
					}
				}

				strJson += "]"

				if Client.HSet(key, s.Field(i).Name, strJson).Result(); err != nil {
					return err
				}
			}

		default:
		}

	}
	return nil
}

func Get(obj interface{}, UUID string) (err error) {

	nowTypeElem := reflect.ValueOf(obj).Elem()
	types := nowTypeElem.Type()

	for i := 0; i < nowTypeElem.NumField(); i++ {
		nowField := nowTypeElem.Field(i)
		nowFieldName := types.Field(i).Name

		switch nowField.Kind() {

		case reflect.String:
			nowValue, err := Client.HGet(UUID, nowFieldName).Result()
			nowField.SetString(nowValue)
			if err != nil {
				return err
			}

		case reflect.Bool:
			nowValue, err := Client.HGet(UUID, nowFieldName).Result()
			nowField.SetBool(utils.StringToBool(nowValue))
			if err != nil {
				return err
			}

		case reflect.Int64:
			nowValue, err := Client.HGet(UUID, nowFieldName).Result()
			nowField.SetInt(utils.StringToInt64(nowValue))
			if err != nil {
				return err
			}

		case reflect.Slice:
			if "[]string" == nowField.Type().String() {
				nowValue, err := Client.HGet(UUID, nowFieldName).Result()

				var stringSlice []string
				err = json.Unmarshal([]byte(nowValue), &stringSlice)

				//TBD : code as below just for testing,it will be updated later
				if err != nil && (len(nowValue) > 0) {
					//return err
				} else {
					sliceValue := reflect.ValueOf(stringSlice)
					nowField.Set(sliceValue)
				}
			}

		default:
		}
	}

	return nil
}
