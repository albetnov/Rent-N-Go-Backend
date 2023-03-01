package user

import "gorm.io/gorm"

// Sok sok normalisasi, no null value
// moga tak menyesal hehe.
type Nik struct {
	gorm.Model
	UserID     uint
	Nik        string `gorm:"unique"`
	IsVerified bool
}
