package user

import "errors"

type User struct {
	ID           string `json:"id"`
	Username     string `json:"username"`
	PasswordHash string `json:"-"`
}

type Repository interface {
	Create(user *User) error
	GetByUsername(username string) (*User, error)
}

type MemoryRepository struct {
	users map[string]*User
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		users: make(map[string]*User),
	}
}

func (r *MemoryRepository) Create(u *User) error {
	for _, existing := range r.users {
		if existing.Username == u.Username {
			return errors.New("username already exists")
		}
	}
	r.users[u.ID] = u
	return nil
}

func (r *MemoryRepository) GetByUsername(username string) (*User, error) {
	for _, u := range r.users {
		if u.Username == username {
			return u, nil
		}
	}
	return nil, errors.New("user not found")
}
