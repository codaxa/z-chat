package postgres

import (
	"context"
	"z-chat/internal/domain/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// RoomRepository provides methods to interact with the rooms table in the database.
type RoomRepository struct {
	db interface {
		Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
		QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
		Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
		Close()
	}
}

// NewRoomRepository creates a new instance of RoomRepository with the provided database connection.
func NewRoomRepository(db interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	Close()
}) *RoomRepository {
	return &RoomRepository{
		db: db,
	}
}

// CreateRoom inserts a new room into the rooms table.
func (r *RoomRepository) CreateRoom(ctx context.Context, room *models.Room) error {
	query := `INSERT INTO rooms (id, name, created_by, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.Exec(ctx, query, room.ID, room.Name, room.CreatedBy, room.CreatedAt, room.UpdatedAt)
	return err
}

// GetRoomByID retrieves a room by its ID from the rooms table.
func (r *RoomRepository) GetRoomByID(ctx context.Context, id string) (*models.Room, error) {
	query := `SELECT id, name, created_by, created_at, updated_at FROM rooms WHERE id = $1`
	row := r.db.QueryRow(ctx, query, id)

	var rom models.Room
	err := row.Scan(&rom.ID, &rom.Name, &rom.CreatedBy, &rom.CreatedAt, &rom.UpdatedAt)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &rom, nil
}

// GetRooms retrieves rooms with pagination
func (r *RoomRepository) GetRooms(ctx context.Context, limit, offset int) ([]*models.Room, error) {
	query := `SELECT id, name, created_by, created_at, updated_at FROM rooms 
              ORDER BY created_at DESC LIMIT $1 OFFSET $2`

	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rooms []*models.Room
	for rows.Next() {
		var room models.Room
		err := rows.Scan(&room.ID, &room.Name, &room.CreatedBy, &room.CreatedAt, &room.UpdatedAt)
		if err != nil {
			return nil, err
		}
		rooms = append(rooms, &room)
	}

	return rooms, rows.Err()
}

// GetRoomAdmins retrieves all admins for a room
func (r *RoomRepository) GetRoomAdmins(ctx context.Context, roomID string) ([]*models.User, error) {
	query := `SELECT u.id, u.username, u.email, u.created_at, u.updated_at 
              FROM users u 
              JOIN room_members rm ON u.id = rm.user_id 
              WHERE rm.room_id = $1 AND rm.role = 'admin'`

	rows, err := r.db.Query(ctx, query, roomID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var admins []*models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, err
		}
		admins = append(admins, &user)
	}

	return admins, rows.Err()
}

// AddRoomAdmin adds a user as admin to a room
func (r *RoomRepository) AddRoomAdmin(ctx context.Context, roomID, userID string) error {
	query := `INSERT INTO room_members (room_id, user_id, role, joined_at) VALUES ($1, $2, 'admin', NOW())`
	_, err := r.db.Exec(ctx, query, roomID, userID)
	return err
}

// RemoveRoomAdmin removes a user as admin from a room
func (r *RoomRepository) RemoveRoomAdmin(ctx context.Context, roomID, userID string) error {
	query := `DELETE FROM room_members WHERE room_id = $1 AND user_id = $2`
	_, err := r.db.Exec(ctx, query, roomID, userID)
	return err
}

// DeleteRoom deletes a room and all associated data
func (r *RoomRepository) DeleteRoom(ctx context.Context, id string) error {
	// This will cascade delete room_members, and messages due to foreign keys
	query := `DELETE FROM rooms WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

// GetRoomMembers retrieves all members for a room
func (r *RoomRepository) GetRoomMembers(ctx context.Context, roomID string) ([]*models.User, error) {
	query := `SELECT u.id, u.username, u.email, u.created_at, u.updated_at 
              FROM users u 
              JOIN room_members rm ON u.id = rm.user_id 
              WHERE rm.room_id = $1`

	rows, err := r.db.Query(ctx, query, roomID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []*models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, err
		}
		members = append(members, &user)
	}

	return members, rows.Err()
}

// AddRoomMember adds a user as member to a room
func (r *RoomRepository) AddRoomMember(ctx context.Context, roomID, userID string) error {
	query := `INSERT INTO room_members (room_id, user_id, joined_at, role) 
              VALUES ($1, $2, NOW(), 'member')`
	_, err := r.db.Exec(ctx, query, roomID, userID)
	return err
}

// RemoveRoomMember removes a user from a room
func (r *RoomRepository) RemoveRoomMember(ctx context.Context, roomID, userID string) error {
	query := `DELETE FROM room_members WHERE room_id = $1 AND user_id = $2`
	_, err := r.db.Exec(ctx, query, roomID, userID)
	return err
}

// IsRoomMember checks if a user is a member of a room
func (r *RoomRepository) IsRoomMember(ctx context.Context, roomID, userID string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM room_members WHERE room_id = $1 AND user_id = $2)`
	var exists bool
	err := r.db.QueryRow(ctx, query, roomID, userID).Scan(&exists)
	return exists, err
}

// IsRoomAdmin checks if a user is an admin of a room
func (r *RoomRepository) IsRoomAdmin(ctx context.Context, roomID, userID string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM room_members WHERE room_id = $1 AND user_id = $2 AND role = 'admin')`
	var exists bool
	err := r.db.QueryRow(ctx, query, roomID, userID).Scan(&exists)
	return exists, err
}

// GetUserRooms retrieves all rooms a user is a member of
func (r *RoomRepository) GetUserRooms(ctx context.Context, userID string) ([]*models.Room, error) {
	query := `SELECT r.id, r.name, r.created_by, r.created_at, r.updated_at 
              FROM rooms r 
              JOIN room_members rm ON r.id = rm.room_id 
              WHERE rm.user_id = $1 
              ORDER BY r.created_at DESC`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rooms []*models.Room
	for rows.Next() {
		var room models.Room
		err := rows.Scan(&room.ID, &room.Name, &room.CreatedBy, &room.CreatedAt, &room.UpdatedAt)
		if err != nil {
			return nil, err
		}
		rooms = append(rooms, &room)
	}

	return rooms, rows.Err()
}
