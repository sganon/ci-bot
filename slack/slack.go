package slack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Attachment struct {
	Color      string  `json:"color"`
	Pretext    string  `json:"pretext,omitempty"`
	AuthorName string  `json:"author_name,omitempty"`
	Title      string  `json:"title"`
	TitleLink  string  `json:"title_link,omitempty"`
	Text       string  `json:"text"`
	Fields     []Field `json:"fields"`

	Fallback string `json:"fallback"`
}

type Field struct {
	Title string `json:"title"`
	Value string `json:"value"`
	Short bool   `json:"short"`
}

func (a Attachment) Send(hookURL string) error {
	client := http.Client{}

	data := struct {
		Attachments []Attachment `json:"attachments"`
	}{Attachments: []Attachment{a}}

	b, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("slack send error: marshal attachement: %v", err)
	}
	req, err := http.NewRequest("POST", hookURL, bytes.NewBuffer(b))
	if err != nil {
		return fmt.Errorf("slack send error: creating request: %v", err)
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("slack send error: sending request: %v", err)
	}

	fmt.Println(res.StatusCode)
	b, _ = ioutil.ReadAll(res.Body)
	fmt.Println(string(b))
	return nil
}

func ErrorMessage(msg string) Attachment {
	attch := Attachment{
		Color: "#C93C20",
		Title: "Error!",
		Fields: []Field{
			Field{Title: "❌", Value: msg},
		},
		Fallback: "Error: " + msg,
	}
	return attch
}
