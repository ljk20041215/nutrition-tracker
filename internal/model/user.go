package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID            string         `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Email         string         `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	PasswordHash  string         `gorm:"not null" json:"-"`
	Nickname      string         `gorm:"type:varchar(50)" json:"nickname"`
	Gender        int            `gorm:"type:int;default:0" json:"gender"` // 0:未知,1:男,2:女
	Age           int            `gorm:"type:int" json:"age"`
	Height        float64        `gorm:"type:float" json:"height"`                 // cm
	Weight        float64        `gorm:"type:float" json:"weight"`                 // kg
	ActivityLevel int            `gorm:"type:int;default:3" json:"activity_level"` // 1-5
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"` // 软删除字段
}

// 注意：如果使用软删除，查询时会自动过滤已删除的记录
// 如果要查询所有记录（包括已删除的），需要使用 Unscoped()
