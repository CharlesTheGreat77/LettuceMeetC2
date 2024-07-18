package parser

import (
	"time"
)

type Lettuce struct {
	Data struct {
		Event struct {
			ID                       string      `json:"id"`
			Title                    string      `json:"title"`
			Description              string      `json:"description"`
			Type                     int         `json:"type"`
			PollStartTime            string      `json:"pollStartTime"`
			PollEndTime              string      `json:"pollEndTime"`
			MaxScheduledDurationMins interface{} `json:"maxScheduledDurationMins"`
			TimeZone                 string      `json:"timeZone"`
			PollDates                []string    `json:"pollDates"`
			Start                    interface{} `json:"start"`
			End                      interface{} `json:"end"`
			IsScheduled              bool        `json:"isScheduled"`
			CreatedAt                time.Time   `json:"createdAt"`
			UpdatedAt                time.Time   `json:"updatedAt"`
			User                     interface{} `json:"user"`
			GoogleEvents             interface{} `json:"googleEvents"`
			PollResponses            []struct {
				ID   string `json:"id"`
				User struct {
					Typename string `json:"__typename"`
					Name     string `json:"name"`
					Email    string `json:"email"`
				} `json:"user"`
				Availabilities []struct {
					Start time.Time `json:"start"`
					End   time.Time `json:"end"`
				} `json:"availabilities"`
				Event struct {
					ID string `json:"id"`
				} `json:"event"`
			} `json:"pollResponses"`
		} `json:"event"`
	} `json:"data"`
}

type LettuceSchedule struct {
	Data struct {
		CreatePollResponse struct {
			PollResponse struct {
				ID   string `json:"id"`
				User struct {
					Typename string `json:"__typename"`
					Name     string `json:"name"`
					Email    string `json:"email"`
				} `json:"user"`
				Availabilities []struct {
					Start time.Time `json:"start"`
					End   time.Time `json:"end"`
				} `json:"availabilities"`
				Event struct {
					ID            string      `json:"id"`
					UpdatedAt     time.Time   `json:"updatedAt"`
					User          interface{} `json:"user"`
					PollResponses []struct {
						ID    string `json:"id"`
						Event struct {
							ID string `json:"id"`
						} `json:"event"`
						User struct {
							Typename string `json:"__typename"`
							Name     string `json:"name"`
							Email    string `json:"email"`
						} `json:"user"`
						Availabilities []struct {
							Start time.Time `json:"start"`
							End   time.Time `json:"end"`
						} `json:"availabilities"`
					} `json:"pollResponses"`
				} `json:"event"`
			} `json:"pollResponse"`
		} `json:"createPollResponse"`
	} `json:"data"`
}

type LettuceEvent struct {
	Data struct {
		CreateEvent struct {
			Event struct {
				ID string `json:"id"`
			} `json:"event"`
		} `json:"createEvent"`
	} `json:"data"`
}