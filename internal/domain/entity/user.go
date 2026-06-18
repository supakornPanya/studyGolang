package entity

type User struct {
	ID           uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	Username     string `gorm:"uniqueIndex;not null" json:"username"`
	PasswordHash string `gorm:"column:password;not null" json:"-"`
	// Role: {admin, seller, customer}
	Role         string `gorm:"not null" json:"role"`
}
