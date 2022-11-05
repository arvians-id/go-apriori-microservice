package model

import "time"

type Category struct {
	IdCategory int       `json:"id_category"`
	Name       string    `json:"name"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type CreateCategoryRequest struct {
	Name string `json:"name"`
}

type UpdateCategoryRequest struct {
	IdCategory int    `json:"id_category"`
	Name       string `json:"name"`
}
