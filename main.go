package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type RestResponse struct {
	PageInfo struct {
		TotalResults int `json:"totalResults"`
	} `json:"pageInfo"`
}

func main() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Println(err)
	}
	for {
		OpenCsv()
		time.Sleep(viper.GetDuration("delay") * time.Minute)
	}
}

func OpenCsv() {
	file, err := os.Open("videos.csv")
	if err != nil {
		log.Println(err)
	} else {
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			lastComma := strings.LastIndex(line, ",")
			id := line[lastComma+1:]
			GetRequest(line, id)
		}
	}
	file.Close()
}
func GetRequest(line string, id string) {
	resp, err := http.Get(fmt.Sprintf("https://youtube.googleapis.com/youtube/v3/videos?part=status&id=%s&key=%s", id, viper.GetString("yt-api-key")))
	if err != nil {
		return
	}
	defer resp.Body.Close()

	var jsonResp RestResponse

	err = json.NewDecoder(resp.Body).Decode(&jsonResp)
	if err != nil {
		return
	} else if jsonResp.PageInfo.TotalResults == 0 {
		http.Get(fmt.Sprintf("https://api.telegram.org/%s/sendMessage?chat_id=%s&text=%s", viper.GetString("tg-api-key"), viper.GetString("tg-chat-id"), line))
	}
}
