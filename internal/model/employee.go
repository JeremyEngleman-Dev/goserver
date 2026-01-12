package model

type Employee struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}
