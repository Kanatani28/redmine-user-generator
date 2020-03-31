package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"gopkg.in/yaml.v2"
)

type Conf struct {
	ApiKey   string `yaml:"api_key"`
	Host     string
	AuthUser string `yaml:"auth_user"`
	AuthPass string `yaml:"auth_pass"`
}

const JSON_TEMPLATE = `{"user": {"login": "${login}", "password": "${password}", "firstname": "${firstname}", "lastname": "${lastname}", "mail": "${mail}", "admin": "${admin}"}}`

const CONFIG_FILE = "config.yml"

// CSV file name
const USER_DATA_CSV = "users.csv"

// CSV columns
const (
	LOGIN = iota
	PASSWORD
	FIRSTNAME
	LASTNAME
	MAIL
	ADMIN
)

func main() {

	c, err := ioutil.ReadFile(CONFIG_FILE)
	if err != nil {
		log.Fatal(err)
	}
	conf := Conf{}
	yaml.Unmarshal([]byte(c), &conf)

	users := readUserData()

	requestURL := `http://` + conf.Host + `/users.json?key=` + conf.ApiKey

	for _, user := range users {
		fmt.Println(requestURL)
		fmt.Println(conf.AuthUser + ":" + conf.AuthPass)
		fmt.Println(getJSON(user))

		client := &http.Client{Timeout: time.Duration(10) * time.Second}
		jsonData := bytes.NewBuffer([]byte(getJSON(user)))

		req, err := http.NewRequest("POST", requestURL, jsonData)
		if err != nil {
			log.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		req.SetBasicAuth(conf.AuthUser, conf.AuthPass)
		resp, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()

		fmt.Println(string(getContent(resp)))

	}
}

func getContent(resp *http.Response) []byte {
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	return b
}

func readUserData() [][]string {

	file, err := os.Open(USER_DATA_CSV)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	users, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	return users[1:]
}

func getJSON(user []string) string {
	jsonStr := strings.Replace(JSON_TEMPLATE, "${login}", user[LOGIN], -1)
	jsonStr = strings.Replace(jsonStr, "${password}", user[PASSWORD], -1)
	jsonStr = strings.Replace(jsonStr, "${firstname}", user[FIRSTNAME], -1)
	jsonStr = strings.Replace(jsonStr, "${lastname}", user[LASTNAME], -1)
	jsonStr = strings.Replace(jsonStr, "${mail}", user[MAIL], -1)
	jsonStr = strings.Replace(jsonStr, "${admin}", user[ADMIN], -1)
	return jsonStr
}
