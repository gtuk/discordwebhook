package discordwebhook

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"net/http"
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
			err := resp.Body.Close()
			if err != nil {
				return err
			}
			return nil
		case http.StatusTooManyRequests:
			// Rate limit exceeded, retry after backoff duration
			var response DiscordResponse
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}
			err = json.Unmarshal(body, &response)
			if err != nil {
				return err
			}

			/*
				Calculate the time until reset and add it to the current local time.
				Some extra time of 750ms is added because without it I still encountered 429s.
			*/

			if response.RetryAfter != 0 {

				whole, frac := math.Modf(response.RetryAfter)
				resetAt := time.Now().Add(time.Duration(whole) * time.Second).Add(time.Duration(frac*1000) * time.Millisecond).Add(750 * time.Millisecond)
				time.Sleep(time.Until(resetAt))
			} else {
				time.Sleep(5 * time.Second)
			}

			err = resp.Body.Close()
			if err != nil {
				return err
			}
		default:
			// Handle other HTTP status codes
			err := resp.Body.Close()
			if err != nil {
				return err
			}
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
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {

			}
		}(resp.Body)

		responseBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		return fmt.Errorf(string(responseBody))
	}

	return nil
}
