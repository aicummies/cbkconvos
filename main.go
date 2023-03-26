package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/akamensky/argparse"
)

type CBKConvoMessage struct {
	ID        string `json:"id"`
	Type      string `json:"type"`
	Text      string `json:"text"`
	CreatedAt int    `json:"createdAt"`
	UpdatedAt int    `json:"updatedAt"`
}

type CBKConvoResponse struct {
	Items []CBKConvoMessage `json:"items"`
}

type TavernMessageHeader struct {
	UserName      string `json:"user_name"`
	CharacterName string `json:"character_name"`
	CreateDate    int    `json:"create_date"`
}

type TavernMessage struct {
	Name     string `json:"name"`
	IsUser   bool   `json:"is_user"`
	IsName   bool   `json:"is_name"` // what is this even for? hello? TAI devs?
	SendDate int    `json:"send_date"`
	Message  string `json:"mes"`
}

func main() {
	bearerToken := readBearerToken()

	parser := argparse.NewParser("cbk", "Get your convos in a tavern compatible format. Set up bearer token as a file called \"bearer\" on the folder")
	conversationId := parser.String("c", "conversation", &argparse.Options{Required: true, Help: "Conversation id, taken from the URL"})
	name1 := parser.String("1", "username", &argparse.Options{Required: false, Default: "You", Help: "The user's name, default is \"You\" unless set"})
	name2 := parser.String("2", "charactername", &argparse.Options{Required: true, Help: "The bot's \"roleplaying a\" name"})

	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
	}
	// DO shit
	convos := getConversationMessages(bearerToken, *conversationId)
	saveFormattedConversation(convos, *name1, *name2)
}

func readBearerToken() string {
	filename := "bearer"
	file, err := os.Open(filename)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(file)
	var result string
	for scanner.Scan() {
		result = scanner.Text()
		break
	}
	return result
}

func getConversationMessages(bearerToken, conversationId string) []CBKConvoMessage {
	client := &http.Client{}
	request, err := http.NewRequest("GET", fmt.Sprintf("https://api.chatbotkit.com/v1/conversation/%s/message/list", conversationId), nil)
	if err != nil {
		log.Fatal(err)
	}
	bearerToken = fmt.Sprintf("Bearer %s", bearerToken)
	request.Header.Add("Authorization", bearerToken)
	request.Header.Add("accept", "application/json")
	response, _ := client.Do(request)
	if response.StatusCode != 200 {
		err_details, _ := ioutil.ReadAll(response.Body)
		panic(fmt.Sprintf("Error, got code %d, full details: %s", response.StatusCode, string(err_details)))
	}
	responseData, _ := ioutil.ReadAll(response.Body)
	var responseObject CBKConvoResponse
	json.Unmarshal(responseData, &responseObject)
	return responseObject.Items
}

func saveFormattedConversation(convos []CBKConvoMessage, name1, name2 string) {
	folder := fmt.Sprintf("./chats/%s", name2)
	filename := fmt.Sprint(convos[0].CreatedAt)

	err := os.MkdirAll(folder, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Create(folder + "/" + filename + ".jsonl")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	backstory_file, err := os.Create(folder + "/backstory_" + filename + ".txt")
	if err != nil {
		panic(err)
	}
	defer backstory_file.Close()

	backstory_file.WriteString(convos[0].Text)

	first_message := TavernMessageHeader{
		UserName:      name1,
		CharacterName: name2,
		CreateDate:    convos[0].CreatedAt,
	}
	first_line, _ := json.Marshal(first_message)
	file.WriteString(string(first_line) + "\n")

	convos = convos[1:]
	length := len(convos)
	for index, message := range convos {
		is_user := false
		name := name2
		if message.Type == "user" {
			is_user = true
			name = name1
		}
		data, _ := json.Marshal(TavernMessage{Name: name, IsUser: is_user, IsName: is_user, SendDate: message.CreatedAt, Message: message.Text})
		new_line := "\n"
		if index == length-1 {
			new_line = ""
		}
		file.WriteString(string(data) + new_line)
	}
}
