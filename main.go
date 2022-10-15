package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"time"
)

var (
	port = ":8080"

	status map[string]map[string]any
	nameFileJSON = "status.json"
)

func init()  {
	getDataFromFileJSON()

}

func main() {
	http.HandleFunc("/status", getStatus)
	http.HandleFunc("/", setStatus)
	fmt.Println("Start server on port", port)
	go autoReload(15)

	http.ListenAndServe(port, nil)
	
}

func autoReload(delayReload int) {
	time.Sleep(time.Duration(delayReload) * time.Second)
	
	waterRandomCode := getRandomNumber(2)
	windRandomCode := getRandomNumber(1)

	// fmt.Println(windRandomCode)
	
	changeValueFileJSON(waterRandomCode, windRandomCode)
	getDataFromFileJSON()
	autoReload(delayReload)
}

func getDataFromFileJSON()  {
	jsonFile, err := os.Open(nameFileJSON)

	if err != nil {
		fmt.Println(err)
		return
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)


	err = json.Unmarshal([]byte(byteValue), &status)

	if err != nil {
		fmt.Println(err)
		return
	}	

}

func changeValueFileJSON(waterCode, windCode int)  {
	status["status"]["water"] = waterCode
	status["status"]["wind"] = windCode

	byteValueNew, err := json.Marshal(status)

	if err != nil {
		fmt.Println(err)
		return
	}

	err = ioutil.WriteFile(nameFileJSON, byteValueNew, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func getRandomNumber(random int64) int {
	rand.Seed(time.Now().UnixMilli() + random)
	return rand.Intn(100)
}


func getStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	if r.Method == "GET" {
		json.NewEncoder(w).Encode(&status)
		return	
	}

	http.Error(w, "Invalid Method", http.StatusBadRequest)

}

func setStatus(w http.ResponseWriter, r *http.Request) {
	tpl, err := template.ParseFiles("./view/template.html")

	if err != nil {
		panic(err)
	}

	tpl.Execute(w, nil)
}