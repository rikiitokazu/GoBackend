package models

import (
	"time"
)

type User struct {
	CustomerID        int       `json:"id"`
	Name              string    `json:"name"`
	Email             string    `json:"email"`
	Password          string    `json:"password"`
	RegisteredCourses []string  `json:"registered_courses"`
	DateCreated       time.Time `json:"date_created"`
	LastActive        time.Time `json:"last_active"`
}

type LoginAuth struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
