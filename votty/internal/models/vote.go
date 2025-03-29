package models

type Vote struct {
	PollID string `json:"poll_id"`
	UserID string `json:"user_id"`
	Choice uint64 `json:"choice"`
}
