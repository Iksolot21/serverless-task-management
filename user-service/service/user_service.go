package service

import (
	"database/sql"
	"fmt"
	"time"
	"your-repo/user-service/pb"
)

func GetUserById(db *sql.DB, req *pb.GetUserByIdRequest) (*pb.User, error) {
	var user pb.User
	query := `SELECT id, username, email, created_at FROM users WHERE id = $1`
	row := db.QueryRow(query, req.Id)
	var created time.Time
	err := row.Scan(&user.Id, &user.Username, &user.Email, &created)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	user.CreatedAt = created.String()

	return &user, nil
}

func GetUsers(db *sql.DB, req *pb.GetUsersRequest) (*pb.GetUsersResponse, error) {
	var users []*pb.User
	query := `SELECT id, username, email, created_at FROM users`
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var user pb.User
		var created time.Time
		err = rows.Scan(&user.Id, &user.Username, &user.Email, &created)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		user.CreatedAt = created.String()
		users = append(users, &user)
	}
	return &pb.GetUsersResponse{Users: users}, nil
}
