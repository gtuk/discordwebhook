package discordwebhook

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"strconv"
	"time"
)

/*
SendMessageRateLimitAware is designed for scenarios where there is a need to dispatch numerous webhooks in rapid succession.
Its purpose is to prevent potential bans from Discord by ensuring that the requests are rate-limited,
thus maintaining a responsible and compliant approach to webhook communication.
*/
func SendMessageRateLimitAware(url string, message Message) error {
	// Validate parameters
	if url == "" {
		return errors.New("empty URL")
	}

	for {
		payload := new(bytes.Buffer)

		err := json.NewEncoder(payload).Encode(message)
		if err != nil {
			return err
		}

		// Make the HTTP request
		resp, err := http.Post(url, "application/json", payload)

		if err != nil {
			log.Printf("HTTP request failed: %v", err)
			return err
		}

		switch resp.StatusCode {
		case http.StatusOK, http.StatusNoContent:
			// Success
			resp.Body.Close()
			return nil
		case http.StatusTooManyRequests:
			// Rate limit exceeded, retry after backoff duration
			resetAfter := resp.Header.Get("X-RateLimit-Reset-After")
			parsedAfter, err := strconv.ParseFloat(resetAfter, 64)
			if err != nil {
				return err
			}

			/*
				Calculate the time until reset and add it to the current local time.
				Some extra time of 250ms is added because without it I still encountered 429s.
			*/
			whole, frac := math.Modf(parsedAfter)
			resetAt := time.Now().Add(time.Duration(whole) * time.Second).Add(time.Duration(frac*1000) * time.Millisecond).Add(250 * time.Millisecond)

			time.Sleep(time.Until(resetAt))
			resp.Body.Close()
		default:
			// Handle other HTTP status codes
			resp.Body.Close()
			responseBody, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			return fmt.Errorf("HTTP request failed with status %d, body: \n %s", resp.StatusCode, responseBody)
		}
	}
}

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

		responseBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		return fmt.Errorf(string(responseBody))
	}

	return nil
}
