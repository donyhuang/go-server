package pkg

import (
	"crypto/md5"
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

// GetMd5Str 小写的16进制表示方法
func GetMd5Str(data []byte) string {
	v := md5.Sum(data)
	return fmt.Sprintf("%x", v)
}

func HandleUA(ua string) string {
	escapeUA, err := url.QueryUnescape(ua)
	if err == nil {
		ua = escapeUA
	}
	if index := strings.Index(ua, "Android"); index == -1 {
		return handleIosUA(ua)
	}
	return handleAndroidUA(ua)
}
func handleIosUA(ua string) string {

	exp, err := regexp.Compile(`(\([^\)]+\))`)
	if err == nil {
		matches := exp.FindStringSubmatch(ua)
		if len(matches) > 1 {
			ua = matches[1]
		}
	}
	ua = strings.Replace(ua, " U;", "", -1)
	ua = strings.Replace(ua, "; wv", "", -1)
	ua = strings.Replace(ua, "\\s\\w\\w-\\w\\w;", "", -1)
	ua = strings.Replace(ua, "0.0;", "0;", -1)
	ua = strings.Replace(ua, "1.0;", "1;", -1)
	ua = strings.Replace(ua, "zh-cn;", "", -1)
	ua = strings.Replace(ua, "zh-CN;", "", -1)
	ua = strings.Replace(ua, " ", "", -1)
	ua = strings.Replace(ua, "%20", "", -1)

	return ua

}
func handleAndroidUA(ua string) string {
	ua = strings.Replace(ua, " U;", "", -1)
	ua = strings.Replace(ua, "; wv", "", -1)
	re := regexp.MustCompile(`\s\w\w-\w\w;`)
	ua = re.ReplaceAllString(ua, "")
	re = regexp.MustCompile(`Android\s\d+.*?\/`)
	uaList := re.FindStringSubmatch(ua)
	if len(uaList) == 0 {
		return ""
	}
	uas := strings.Split(uaList[0], ";")
	var system, model string
	for i := 0; i <= 1; i++ {
		if strings.Index(uas[i], "Android") != -1 {
			re = regexp.MustCompile(`Android\s\d+`)
			systemList := re.FindStringSubmatch(uas[i])
			system = systemList[0]
		} else if strings.Index(uas[i], "/") != -1 {
			model = strings.TrimSpace(uas[i])
			modeList := strings.Split(model, " ")
			model = modeList[0]
		}
	}
	return strings.ReplaceAll(system+model, " ", "")
}
