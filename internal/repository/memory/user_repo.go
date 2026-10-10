package memory

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"aarambha-agency/internal/domain"
)

type UserRepo struct {
	mu      sync.RWMutex
	users   map[string]*domain.User
	byEmail map[string]string // lowercase email -> user ID
}

func NewUserRepo() *UserRepo {
	repo := &UserRepo{
		users:   make(map[string]*domain.User),
		byEmail: make(map[string]string),
	}
	return repo
}

func (r *UserRepo) Create(user *domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	cleanedEmail := strings.ToLower(strings.TrimSpace(user.Email))
	if _, exists := r.byEmail[cleanedEmail]; exists {
		return domain.ErrEmailAlreadyExists
	}

	if user.ID == "" {
		user.ID = fmt.Sprintf("USER-%d", time.Now().UnixNano())
	}
	now := time.Now()
	if user.CreatedAt.IsZero() {
		user.CreatedAt = now
	}
	if user.UpdatedAt.IsZero() {
		user.UpdatedAt = now
	}
	user.Email = cleanedEmail

	cloned := *user
	r.users[user.ID] = &cloned
	r.byEmail[cleanedEmail] = user.ID
	return nil
}

func (r *UserRepo) FindByEmail(email string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	cleanedEmail := strings.ToLower(strings.TrimSpace(email))
	userID, exists := r.byEmail[cleanedEmail]
	if !exists {
		return nil, domain.ErrUserNotFound
	}

	user, ok := r.users[userID]
	if !ok {
		return nil, domain.ErrUserNotFound
	}

	cloned := *user
	return &cloned, nil
}

func (r *UserRepo) FindByID(id string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, ok := r.users[id]
	if !ok {
		return nil, domain.ErrUserNotFound
	}

	cloned := *user
	return &cloned, nil
}

func (r *UserRepo) Update(user *domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	existing, ok := r.users[user.ID]
	if !ok {
		return domain.ErrUserNotFound
	}

	cleanedEmail := strings.ToLower(strings.TrimSpace(user.Email))
	if cleanedEmail != strings.ToLower(existing.Email) {
		if _, exists := r.byEmail[cleanedEmail]; exists {
			return domain.ErrEmailAlreadyExists
		}
		delete(r.byEmail, strings.ToLower(existing.Email))
		r.byEmail[cleanedEmail] = user.ID
	}

	user.UpdatedAt = time.Now()
	cloned := *user
	r.users[user.ID] = &cloned
	return nil
}

func (r *UserRepo) UpdateLastLogin(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	user, ok := r.users[id]
	if !ok {
		return domain.ErrUserNotFound
	}

	now := time.Now()
	user.LastLoginAt = &now
	user.UpdatedAt = now
	return nil
}

func (r *UserRepo) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.users)
}

func (r *UserRepo) FindAll() []domain.User {
	r.mu.RLock()
	defer r.mu.RUnlock()

	list := make([]domain.User, 0, len(r.users))
	for _, u := range r.users {
		list = append(list, *u)
	}
	return list
}
