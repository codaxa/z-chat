// Package repository defines interfaces for data access operations
package repository

import (
	"context"
	"z-chat/internal/domain/models"
)

// RoomRepository defines the interface for room-related operations.
type RoomRepository interface {
	CreateRoom(ctx context.Context, room *models.Room) error
	GetRoomByID(ctx context.Context, id string) (*models.Room, error)
	GetRooms(ctx context.Context, limit, offset int) ([]*models.Room, error)
	GetRoomAdmins(ctx context.Context, roomID string) ([]*models.User, error)
	AddRoomAdmin(ctx context.Context, roomID, userID string) error
	RemoveRoomAdmin(ctx context.Context, roomID, userID string) error
	DeleteRoom(ctx context.Context, id string) error
	GetRoomMembers(ctx context.Context, roomID string) ([]*models.User, error)
	AddRoomMember(ctx context.Context, roomID, userID string) error
	RemoveRoomMember(ctx context.Context, roomID, userID string) error
	IsRoomMember(ctx context.Context, roomID, userID string) (bool, error)
	IsRoomAdmin(ctx context.Context, roomID, userID string) (bool, error)
	GetUserRooms(ctx context.Context, userID string) ([]*models.Room, error)
}
