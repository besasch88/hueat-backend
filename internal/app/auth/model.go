package auth

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type authSessionModel struct {
	ID           uuid.UUID `gorm:"primaryKey;column:id;type:varchar(36)"`
	UserID       string    `gorm:"column:user_id;type::varchar(36)"`
	CreatedAt    time.Time `gorm:"column:created_at;type:timestamp"`
	ExpiresAt    time.Time `gorm:"column:expires_at;type:timestamp"`
	RefreshToken string    `gorm:"column:refresh_token;type:text"`
}

func (m authSessionModel) TableName() string {
	return "hueat_auth_session"
}

func (m authSessionModel) toEntity() authSessionEntity {
	return authSessionEntity(m)
}

type authUserModel struct {
	ID          uuid.UUID      `gorm:"primaryKey;column:id;type:varchar(36)"`
	Username    string         `gorm:"column:username;type::varchar(255)"`
	Password    string         `gorm:"column:password;type::varchar(255)"`
	Permissions pq.StringArray `gorm:"column:permissions;type::text[]"`
	CreatedAt   time.Time      `gorm:"column:created_at;type:timestamp"`
}

func (m authUserModel) TableName() string {
	return "hueat_user"
}

func (m authUserModel) toEntity() authUserEntity {
	return authUserEntity{
		ID:          m.ID,
		Username:    m.Username,
		Password:    m.Password,
		Permissions: m.Permissions,
		CreatedAt:   m.CreatedAt,
	}
}
