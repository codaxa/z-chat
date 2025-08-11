package models

import "time"

// Room represents a chat room in the application
type Room struct {
	ID        string    `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name      string    `json:"name" gorm:"type:varchar(50);uniqueIndex;not null" validate:"required,min=3,max=50"`
	CreatedBy string    `json:"created_by" gorm:"type:uuid;not null" validate:"required"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	Messages  []Message `json:"messages,omitempty" gorm:"foreignKey:RoomID"`
}

// RoomMember represents the junction table for room memberships
type RoomMember struct {
	ID       string    `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	RoomID   string    `json:"room_id" gorm:"type:uuid;not null" validate:"required"`
	UserID   string    `json:"user_id" gorm:"type:uuid;not null" validate:"required"`
	JoinedAt time.Time `json:"joined_at" gorm:"autoCreateTime"`
	Role     string    `json:"role" gorm:"type:varchar(20);default:'member'" validate:"oneof=admin member"`
	Room     Room      `json:"room" gorm:"foreignKey:RoomID"`
	User     User      `json:"user" gorm:"foreignKey:UserID"`
}

// TableName returns the table name for the Room model
func (Room) TableName() string {
	return "rooms"
}

// TableName returns the table name for the RoomMember model
func (RoomMember) TableName() string {
	return "room_members"
}

// Validate checks the Room fields for validity
func (r *Room) Validate() error {
	return validate.Struct(r)
}

// Validate checks the RoomMember fields for validity
func (rm *RoomMember) Validate() error {
	return validate.Struct(rm)
}
