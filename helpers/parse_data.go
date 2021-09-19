package helpers

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

var aid = "U01MENEF744"
var rid = "U01MM6PS3PB"

// data stores unmarshalled messages of all the files
var data = map[string]File{}

// result stores the all the data in our required format which is used for building the CSV file
var result []map[Reply]string

// File is a slice of all Messages in a file
type File []Message

// Message slice type is used for storing ALL messages in the file
type Message struct {
	ClientMsgID     string   `json:"client_msg_id"`
	Type            string   `json:"type"`
	Text            string   `json:"text"`
	User            string   `json:"user"`
	ReplyUsersCount int      `json:"reply_users_count"`
	ReplyUsers      []string `json:"reply_users"`
	Replies         []Reply  `json:"replies"`
	ParentUserId    string   `json:"parent_user_id"`
	ThreadTs        string   `json:"thread_ts"`
}

type Reply struct {
	User string `json:"user"`
	Ts   string `json:"ts"`
}

// walk parses each file and stores the bytes in rawData
func walk(path string, d fs.DirEntry, e error) error {
	if e != nil {
		return e
	}
	if !d.IsDir() {
		var file File
		byteVal, _ := os.ReadFile(path)
		json.Unmarshal(byteVal, &file)
		data[path] = file

	}
	return nil
}

func buildCSV() {
	// 1. Store relevant information in result
	for _, file := range data {
		for _, msg := range file {
			// msg shouldn't be from alex or rahul and should've a reply
			if msg.User != aid || msg.User != rid && msg.ReplyUsersCount > 0 {
				for _, reply := range msg.Replies {
					if reply.User == aid || reply.User == rid {
						mp := make(map[Reply]string)
						mp[reply] = msg.Text
						result = append(result, mp)
					}
				}
			}
		}
	}
	fmt.Println(result)

	// 2. Convert Reply struct to actual replies and build CSV

}

func Start() {
	filepath.WalkDir("./data/side-projects/2021-09-16.json", walk)
	buildCSV()
}

// Helpers
// contains returns true if val is present in arr, else false
func contains(arr []string, val string) bool {
	for _, v := range arr {
		if v == val {
			return true
		}
	}
	return false
}
