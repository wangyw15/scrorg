package organizer

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/wangyw15/scrorg/steam"

	"github.com/spf13/viper"
)

func SanitizeAppName(appname string) string {
	invalid_chars := `<>:"/\|?*`

	// replace invalid characters with space
	for _, char := range invalid_chars {
		appname = strings.ReplaceAll(appname, string(char), " ")
	}

	// replace multiple spaces with single space
	space_pattern := regexp.MustCompile(`\s+`)
	appname = space_pattern.ReplaceAllString(appname, " ")
	return appname
}

func RunOrganize() {
	// get target directory
	target := viper.GetString("target")
	if target == "" {
		fmt.Println("\"target\" is not set")
	}

	// get proxy
	proxy := viper.GetString("proxy")
	if proxy == "" {
		if proxy_url, err := http.ProxyFromEnvironment(&http.Request{}); err == nil {
			proxy = proxy_url.String()
		}
	}

	// print settings
	fmt.Println("Using proxy", proxy)
	fmt.Println("Working in", target)

	// create helper
	helper := steam.NewHelper(proxy)
	defer helper.Save()

	// read directory
	files, err := os.ReadDir(target)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		fmt.Print(file.Name())
		appid := strings.Split(file.Name(), "_")[0]
		appname := helper.GetAppName(appid)
		appname = SanitizeAppName(appname)
		if appname != "" {
			fmt.Println(" ->", appname)
			os.MkdirAll(filepath.Join(target, appname), 0777)
			os.Rename(filepath.Join(target, file.Name()), filepath.Join(target, appname, file.Name()))
		}
	}

	fmt.Println("Completed")
	fmt.Scanln()
}
