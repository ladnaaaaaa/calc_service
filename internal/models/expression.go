package models

import (
	"gorm.io/gorm"
)

type ExpressionStatus string

const (
	StatusPending   ExpressionStatus = "pending"
	StatusComputing ExpressionStatus = "computing"
	StatusCompleted ExpressionStatus = "completed"
	StatusError     ExpressionStatus = "error"
)

type Expression struct {
	gorm.Model
	Expression string           `gorm:"not null"`
	Status     ExpressionStatus `gorm:"not null;default:'pending'"`
	Result     float64
	UserID     uint   `gorm:"not null"`
	User       User   `gorm:"foreignKey:UserID"`
	Tasks      []Task `gorm:"foreignKey:ExpressionID"`
}
