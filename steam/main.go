package steam

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"
)

type SteamResponse struct {
	Success bool `json:"success"`
	Data    struct {
		Name string `json:"name"`
	} `json:"data"`
}

type SteamHelper struct {
	AppNames map[string]string
	Client   *http.Client
}

func NewHelper(proxy string) *SteamHelper {
	transport := &http.Transport{}

	// set proxy
	if proxy != "" {
		if parsed_proxy, err := url.Parse(proxy); err == nil {
			transport.Proxy = http.ProxyURL(parsed_proxy)
		}
	}

	// create client
	client := &http.Client{Transport: transport}

	// load existing app names
	games_file, error := os.OpenFile("games.json", os.O_RDONLY, 0666)
	appnames := make(map[string]string)
	if error == nil {
		defer games_file.Close()
		json.NewDecoder(games_file).Decode(&appnames)
	}

	return &SteamHelper{
		AppNames: appnames,
		Client:   client,
	}
}

func (self *SteamHelper) GetAppName(appid string) string {
	// check if app name is already loaded
	if appname, ok := self.AppNames[appid]; ok {
		return appname
	}

	// get app name from steam
	resp, err := self.Client.Get("https://store.steampowered.com/api/appdetails?appids=" + appid)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	// parse response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ""
	}
	var data interface{}
	json.Unmarshal(body, &data)
	success := data.(map[string]interface{})[appid].(map[string]interface{})["success"].(bool)
	if success {
		name := data.(map[string]interface{})[appid].(map[string]interface{})["data"].(map[string]interface{})["name"].(string)
		self.AppNames[appid] = name
		return name
	} else {
		return ""
	}
}

func (self *SteamHelper) Save() {
	games_file, error := os.OpenFile("games.json", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if error == nil {
		defer games_file.Close()
		encoder := json.NewEncoder(games_file)
		encoder.SetEscapeHTML(false)
		encoder.Encode(self.AppNames)
	}
}
