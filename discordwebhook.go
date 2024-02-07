package discordwebhook

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
)

func SendMessage(url string, message Message) error {
	payload := new(bytes.Buffer)

	err := json.NewEncoder(payload).Encode(message)
	if err != nil {
		return err
	}

	resp, err := http.Post(url, "application/json", payload)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 && resp.StatusCode != 204 {
		defer resp.Body.Close()

		responseBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		return fmt.Errorf(string(responseBody))
	}

	return nil
}

func SendFile(url string, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Err: Could not open file")
		return err
	}
	defer file.Close()

	body := new(bytes.Buffer)

	writer := multipart.NewWriter(body)

	pathSplit := strings.Split(filePath, "/")
	part, err := writer.CreateFormFile("file", pathSplit[len(pathSplit)-1])
	if err != nil {
		log.Fatal("Err: Could not create form file")
		return err
	}

	_, err = io.Copy(part, file)
	if err != nil {
		log.Fatal("Err: Could not copy file content to form file field")
		return err
	}

	resp, err := http.Post(url, writer.FormDataContentType(), body)
	if err != nil {
		log.Fatal("Err: Post request failed")
		return err
	}
	defer resp.Body.Close()
	return nil
}
