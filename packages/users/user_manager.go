package users

import (
	"fmt"
	"sync"
)

// UserManager handles the list of users that are stored in the user config.
// It calls a function if the list changed.
type UserManager struct {
	storeCallback func([]*User) error
	storeOnChange bool
	usersLock     sync.RWMutex
	users         map[string]*User
}

// NewUserManager creates a new user manager.
func NewUserManager(storeCallback func([]*User) error) *UserManager {
	return &UserManager{
		storeCallback: storeCallback,
		storeOnChange: false,
		users:         make(map[string]*User),
	}
}

// StoreOnChange sets whether storing changes to the config is active or not.
func (pm *UserManager) StoreOnChange(store bool) {
	pm.storeOnChange = store
}

// Store calls the storeCallback if storeOnChange is active.
func (pm *UserManager) Store() error {
	if !pm.storeOnChange {
		return nil
	}

	if pm.storeCallback == nil {
		return nil
	}

	users := make([]*User, 0, len(pm.users))
	for k := range pm.users {
		users = append(users, pm.users[k])
	}

	if err := pm.storeCallback(users); err != nil {
		return fmt.Errorf("failed to store user config: %w", err)
	}

	return nil
}

// Users returns a copy of all known users.
func (pm *UserManager) Users() map[string]*User {
	pm.usersLock.RLock()
	defer pm.usersLock.RUnlock()

	usersCopy := make(map[string]*User, len(pm.users))
	for k := range pm.users {
		usersCopy[k] = pm.users[k].Clone()
	}

	return usersCopy
}

// User returns a copy of a user.
func (pm *UserManager) User(name string) (*User, error) {
	pm.usersLock.RLock()
	defer pm.usersLock.RUnlock()

	if _, exists := pm.users[name]; !exists {
		return nil, fmt.Errorf("unable to get user: user \"%s\" does not exist", name)
	}

	return pm.users[name].Clone(), nil
}

// AddUser adds a user to the user manager.
func (pm *UserManager) AddUser(user *User) error {
	pm.usersLock.Lock()
	defer pm.usersLock.Unlock()

	if _, exists := pm.users[user.Name]; exists {
		return fmt.Errorf("unable to add user: user \"%s\" already exists", user.Name)
	}

	pm.users[user.Name] = user

	return pm.Store()
}

// ModifyUser modifies a user in the user manager.
func (pm *UserManager) ModifyUser(user *User) error {
	pm.usersLock.Lock()
	defer pm.usersLock.Unlock()

	if _, exists := pm.users[user.Name]; !exists {
		return fmt.Errorf("unable to modify user: user \"%s\" does not exist", user.Name)
	}

	pm.users[user.Name] = user

	return pm.Store()
}

// RemoveUser removes a user from the user manager.
func (pm *UserManager) RemoveUser(name string) error {
	pm.usersLock.Lock()
	defer pm.usersLock.Unlock()

	if _, exists := pm.users[name]; !exists {
		return fmt.Errorf("unable to remove user: user \"%s\" does not exist", name)
	}

	delete(pm.users, name)

	return pm.Store()
}
