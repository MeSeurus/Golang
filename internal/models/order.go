package models

type Order struct {
	ID       int
	Customer string
	Products string
	Total    float64
	Status   string
}
