package models

import (
	"gorm.io/gorm"
)

type Operation string

const (
	OperationAdd      Operation = "+"
	OperationSubtract Operation = "-"
	OperationMultiply Operation = "*"
	OperationDivide   Operation = "/"
)

type Task struct {
	gorm.Model
	ExpressionID uint       `gorm:"not null"`
	Expression   Expression `gorm:"foreignKey:ExpressionID;references:ID"`
	Arg1         float64    `gorm:"not null"`
	Arg2         float64    `gorm:"not null"`
	Operation    Operation  `gorm:"not null"`
	Result       float64
	Status       ExpressionStatus `gorm:"not null;default:'pending'"`
	OrderNum     int              `gorm:"not null;column:order_num"`
}
