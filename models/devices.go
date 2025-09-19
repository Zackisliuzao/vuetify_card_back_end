package models

type Device struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Status   int    `json:"status"`
	Position string `json:"position"`
}
