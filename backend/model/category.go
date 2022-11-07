package model

import "time"

type Category struct {
	IdCategory int64     `json:"id_category"`
	Name       string    `json:"name"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
