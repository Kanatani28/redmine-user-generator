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
	log.Println("# START MAIN")
	conf := loadConfig()

	users := readUserData()

	requestURL := `http://` + conf.Host + `/users.json?key=` + conf.ApiKey
	log.Println("## REQUEST URL : " + requestURL)

	for _, user := range users {

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

	log.Println("# FINISH MAIN")
}

func loadConfig() Conf {
	log.Println("## START LOADING CONFIG")
	c, err := ioutil.ReadFile(CONFIG_FILE)
	if err != nil {
		log.Fatal(err)
	}
	conf := Conf{}
	yaml.Unmarshal([]byte(c), &conf)

	log.Println("## FINISH LOADING CONFIG")

	return conf
}

func readUserData() [][]string {

	log.Println("## START READ USER DATA")

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
	validateUsers(users[1:])

	log.Println("## FINISH READ USER DATA")
	return users[1:]
}

func getContent(resp *http.Response) []byte {
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	return b
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

func validateUsers(users [][]string) {
	for i, user := range users {
		if user[LOGIN] == "" {
			log.Fatalf("loginが入力されていません。%d番目のユーザー", i+1)
		}
		if user[PASSWORD] == "" {
			log.Fatalf("passwordが入力されていません。%d番目のユーザー", i+1)
		}
		if len(user[PASSWORD]) < 8 {
			log.Fatalf("passwordは8桁以上である必要があります。%d番目のユーザー", i+1)
		}
		if user[FIRSTNAME] == "" {
			log.Fatalf("firstnameが入力されていません。%d番目のユーザー", i+1)
		}
		if user[LASTNAME] == "" {
			log.Fatalf("lastnameが入力されていません。%d番目のユーザー", i+1)
		}
		if user[MAIL] == "" {
			log.Fatalf("mailが入力されていません。%d番目のユーザー", i+1)
		}
		if !(user[ADMIN] == "true" || user[ADMIN] == "false") {
			log.Fatalf("adminはtrueかfalseで入力してください。%d番目のユーザー", i+1)
		}
	}
}
