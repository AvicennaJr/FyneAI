package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type myEntry struct {
	widget.Entry
	output *widget.Label
	app    fyne.App
}

func (m *myEntry) TypedShortcut(s fyne.Shortcut) {
	if _, ok := s.(*desktop.CustomShortcut); !ok {
		m.Entry.TypedShortcut(s)
		return
	}

	t := s.(*desktop.CustomShortcut)
	if t.Modifier == fyne.KeyModifierControl {
		switch t.KeyName {
		case fyne.KeyBackslash:
			sendToLLM(m.Entry.Text, m.output)
			m.Entry.SetText("")
		case fyne.KeyBackspace:
			m.app.Quit()
		}
	}
}

const oneMB = 1024 * 1024

// ResponseStruct represents the structure of the JSON response
type ResponseStruct struct {
	Model         string `json:"model"`
	CreatedAt     string `json:"created_at"`
	Response      string `json:"response"`
	Done          bool   `json:"done"`
	Context       []int  `json:"context"`
	TotalDuration int64  `json:"total_duration"`
	LoadDuration  int    `json:"load_duration"`
	// Add other fields as needed
}

func sendToLLM(s string, widget *widget.Label) {
	widget.SetText("")
	data := map[string]interface{}{
		"model":  "zephyr",
		"prompt": s,
	}

	// Marshal the map to JSON
	body, err := json.Marshal(data)

	// HTTP client
	req, err := http.NewRequest("POST", "http://localhost:11434/api/generate", bytes.NewBuffer(body))
	if err != nil {
		log.Fatal("Error creating HTTP request: ", err.Error())
	}

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Error making HTTP request: ", err.Error())
	}
	defer resp.Body.Close() // Close the response body when done

	var responseBuffer bytes.Buffer

	bytesRead := 0
	buf := make([]byte, oneMB)

	// Read the response body
	for {
		n, err := resp.Body.Read(buf)
		bytesRead += n

		if n > 0 {
			responseBuffer.Write(buf[:n])
		}

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatal("Error reading HTTP response: ", err.Error())
		}

		decoder := json.NewDecoder(&responseBuffer)
		var responseObject ResponseStruct
		if err := decoder.Decode(&responseObject); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal("Error decoding JSON: ", err.Error())
		}

		widget.SetText(widget.Text + responseObject.Response)
	}

}

func main() {

	myApp := app.NewWithID("Avi AI")
	myWindow := myApp.NewWindow("Avi AI")
	myWindow.Resize(fyne.NewSize(500, 500))
	input := &myEntry{}
	input.output = widget.NewLabel("")
	input.output.Wrapping = fyne.TextWrapBreak
	input.Scroll = container.ScrollVerticalOnly
	input.app = myApp

	input.MultiLine = true
	input.SetPlaceHolder("Enter text...")

	d := container.NewVScroll(input.output)
	d.SetMinSize(fyne.NewSize(500, 500))
	bottomBox := container.NewHBox(
		widget.NewButtonWithIcon("copy content", theme.ContentCopyIcon(), func() {
			myWindow.Clipboard().SetContent(input.output.Text)
		}),
	)
	content := container.NewVBox(input, d, bottomBox)

	myWindow.SetContent(content)

	myWindow.Canvas().Focus(input)
	myWindow.ShowAndRun()
}
