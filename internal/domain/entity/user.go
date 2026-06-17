package entity

type User struct {
	ID           uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	Username     string `gorm:"uniqueIndex;not null" json:"username"`
	PasswordHash string `gorm:"column:password;not null" json:"-"`
	Email        string `gorm:"column:email" json:"email"`
	CanRead      bool   `gorm:"default:true" json:"can_read"`
	CanWrite     bool   `gorm:"default:false" json:"can_write"`
	CanUpdate    bool   `gorm:"default:false" json:"can_update"`
	CanDelete    bool   `gorm:"default:false" json:"can_delete"`
}
