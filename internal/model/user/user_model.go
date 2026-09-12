package user

import "time"

type UserModel struct {
	ID        int64     `gorm:"primaryKey;autoIncrement;column:id"`
	Name      string    `gorm:"column:name"`
	Email     string    `gorm:"column:email"`
	Password  string    `gorm:"column:password"`
	Role      string    `gorm:"column:role"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (UserModel) TableName() string {
	return "users"
}
