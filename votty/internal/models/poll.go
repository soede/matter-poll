package models

type Poll struct {
	ID       string   `json:"id"`
	OwnerID  string   `json:"owner_id"`
	Question string   `json:"question"`
	Options  []string `json:"options"`
	IsActive bool     `json:"is_active"`
}
