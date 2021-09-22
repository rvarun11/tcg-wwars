package helpers

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

const aid = "U01MENEF744"
const rid = "U01MM6PS3PB"

// data stores unmarshalled messages of all the files
var data []file

// file is a slice of all Messages in a file
type file []message

// message slice type is used for storing ALL messages in the file
type message struct {
	ClientMsgID     string   `json:"client_msg_id"`
	Type            string   `json:"type"`
	Text            string   `json:"text"`
	User            string   `json:"user"`
	Ts              string   `json:"ts"`
	ReplyUsersCount int      `json:"reply_users_count"`
	ReplyUsers      []string `json:"reply_users"`
	Replies         []reply  `json:"replies"`
	ThreadTs        string   `json:"thread_ts"`
	ParentUserId    string   `json:"parent_user_id"`
}

type reply struct {
	User string `json:"user"`
	Ts   string `json:"ts"`
}

// walk parses each file and stores the bytes in rawData
func walk(path string, d fs.DirEntry, e error) error {
	if e != nil {
		return e
	}
	if !d.IsDir() {
		var f file
		byteVal, _ := os.ReadFile(path)
		_ = json.Unmarshal(byteVal, &f)
		data = append(data, f)

	}
	return nil
}

// buildResult builds the result in required format and returns it
func buildResult() [][]string {
	result := [][]string{{"message", "alex_response", "rahul_response"}}
	for i, f := range data {
		for _, msg := range f {
			// msg shouldn't be from alex or rahul and should be a reply
			if msg.User != aid || msg.User != rid && msg.ReplyUsersCount > 0 {
				for _, r := range msg.Replies {
					if r.User == aid || r.User == rid {
						row := buildRow(i, msg.Text, r.Ts)
						result = append(result, row)
					}
				}
			}
		}
	}

	return result
}

// buildRow takes the index of current file with user and ts, and returns the row
// this works because data stores files in linear order so the reply must be in that file or a later file, but not before
func buildRow(index int, text string, ts string) []string {
	row := []string{text, "", ""}

	for i := index; i < len(data); i++ {
		f := data[i]
		for _, msg := range f {
			if msg.Ts == ts {
				switch msg.User {
				case aid:
					row[1] = msg.Text
				case rid:
					row[2] = msg.Text
				}
			}
		}
	}
	return row
}

// BuildCsv builds a CSV file in required format for all json files in given directory
func BuildCsv(path string) {
	_ = filepath.WalkDir(path, walk)

	result := buildResult()
	csvFile, _ := os.Create("./test.csv")
	defer csvFile.Close()

	fmt.Println("Building CSV.....")
	writer := csv.NewWriter(csvFile)
	for _, row := range result {
		_ = writer.Write(row)
	}
	writer.Flush()

}
