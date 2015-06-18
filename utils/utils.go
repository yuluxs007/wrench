package utils

import (
	"encoding/base64"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func Exist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

func EncodeBasicAuth(username string, password string) string {
	auth := username + ":" + password
	msg := []byte(auth)
	authorization := make([]byte, base64.StdEncoding.EncodedLen(len(msg)))
	base64.StdEncoding.Encode(authorization, msg)
	return string(authorization)
}

func DecodeBasicAuth(authorization string) (username string, password string, err error) {
	basic := strings.Split(strings.TrimSpace(authorization), " ")
	if len(basic) <= 1 {
		return "", "", err
	}

	decLen := base64.StdEncoding.DecodedLen(len(basic[1]))
	decoded := make([]byte, decLen)
	authByte := []byte(basic[1])
	n, err := base64.StdEncoding.Decode(decoded, authByte)

	if err != nil {
		return "", "", err
	}
	if n > decLen {
		return "", "", fmt.Errorf("Something went wrong decoding auth config")
	}

	arr := strings.SplitN(string(decoded), ":", 2)
	if len(arr) != 2 {
		return "", "", fmt.Errorf("Invalid auth configuration file")
	}

	username = arr[0]
	password = strings.Trim(arr[1], "\x00")

	return username, password, nil
}

func StringToBool(value string) bool {
	if boolean, _ := strconv.ParseBool(value); boolean == true {
		return true
	} else if boolean == false {
		return false
	} else {
		return false
	}
}

func StringToInt64(value string) int64 {
	retval, err := strconv.ParseInt(value, 0, 64)
	if err != nil {
		fmt.Println("err:=", err)
		return 0
	}
	return retval
}

func BoolToString(boolean bool) string {
	if boolean == true {
		return "true"
	} else {
		return "false"
	}
}

func Int64ToString(i int64) string {
	return strconv.FormatInt(i, 10)
}

func IsDirExists(path string) bool {
	fi, err := os.Stat(path)

	if err != nil {
		return os.IsExist(err)
	} else {
		return fi.IsDir()
	}

	panic("not reached")
}
