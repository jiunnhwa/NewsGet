package url

import (
	"log"
	"net/url"
	mylogger "newsget/service/logger"
	"strings"
)

// www.google.com -> google.com
func GetHostName(urlstr string) string {
	u, err := url.Parse(urlstr)
	if err != nil {
		log.Fatal(err)
	}
	s := strings.Split(strings.Replace(u.Hostname(), "www.", "", -1), ".")
	return (s[0])
}

//x.com -> com
func GetHostSuffix(urlstr string) string {
	//2021/05/17 14:57:32 parse "https://www.irishtimes.com/\t\t\t\t\t\t\t/life-and-style/health-family/parenting/q-a-is-vaccinating-children-the-next-step-to-combat-covid-19-1.4561951\t": net/url: invalid control character in URL
	u, err := url.Parse(urlstr)
	if err != nil {
		log.Println(err)
		mylogger.Trace.Println("GetHostSuffix:", err)
		return "err"
	}
	s := strings.Split(u.Hostname(), ".")
	return (s[len(s)-1])
}
