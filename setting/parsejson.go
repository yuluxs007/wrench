package setting

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

type EndpointsConf struct {
	Endpoints []EndpointDesc `json:"endpoints,omitempty"`
}

type EndpointDesc struct {
	Name      string        `json:"name"`
	URL       string        `json:"url"`
	Headers   http.Header   `json:"headers"`
	Timeout   time.Duration `json:"timeout"`
	Threshold int           `json:"threshold"`
	Backoff   time.Duration `json:"backoff"`
	EventDB   string        `json:"eventdb"`
	Disabled  bool          `json:"disabled"`
}

type AuthorConf struct {
	Authors []AuthorDesc `json:"auth,omitempty"`
}

type AuthorDesc struct {
	Name           string `json:"name"`
	Realm          string `json:"realm"`
	Hostenabled    int    `json:"hostenabled"`
	Service        string `json:"service"`
	Issuer         string `json:"issuer"`
	Rootcertbundle string `json:"rootcertbundle"`
}

type Desc struct {
	EndpointsConf
	AuthorConf
	Raw []byte `json:"-"`
}

var JSONConfCtx Desc

func GetConfFromJSON(path string) error {
	fp, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("Open JSON config failed,err: %v", err.Error())
	}

	buf, err := ioutil.ReadAll(fp)
	if err != nil {
		return fmt.Errorf("Read JSON config failed,err: %v", err.Error())
	}

	if err := json.Unmarshal(buf, &JSONConfCtx); err != nil {
		return fmt.Errorf("Unmarshal endpoints config error! err: %v", err.Error())
	}

	JSONConfCtx.Raw = make([]byte, len(buf), len(buf))
	copy(JSONConfCtx.Raw, buf)

	return nil
}
