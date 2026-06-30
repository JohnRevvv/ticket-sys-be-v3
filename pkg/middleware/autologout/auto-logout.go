package middleware

import (
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	SupAdScript "ideyanale-be/pkg/modules/super-admin/script"
	UserScript "ideyanale-be/pkg/modules/users/script"

	"github.com/gofiber/fiber/v3"
)

type sessionInfo struct {
	LastSeen time.Time
	Role     string
}

type activityTracker struct {
	mu       sync.RWMutex
	sessions map[int]sessionInfo
}

var (
	tracker     = &activityTracker{sessions: make(map[int]sessionInfo)}
	autoTimeout time.Duration
	scannerOnce sync.Once
)

func getAutoLogoutTimeout() time.Duration {
	if autoTimeout != 0 {
		return autoTimeout
	}
	minutes := 30
	if v := os.Getenv("AUTO_LOGOUT_MINUTES"); v != "" {
		if parsed, err := strconv.Atoi(v); err == nil {
			minutes = parsed
		}
	}
	autoTimeout = time.Duration(minutes) * time.Minute
	return autoTimeout
}

func (t *activityTracker) touch(id int, role string) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.sessions[id] = sessionInfo{
		LastSeen: time.Now(),
		Role:     role,
	}
}
func (t *activityTracker) get(id int) (sessionInfo, bool) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	s, ok := t.sessions[id]
	return s, ok
}

// func (t *activityTracker) lastSeenAt(id int) (time.Time, bool) {
// 	t.mu.RLock()
// 	defer t.mu.RUnlock()
// 	ts, ok := t.lastSeen[id]
// 	return ts, ok
// }

func (t *activityTracker) remove(id int) {
	t.mu.Lock()
	defer t.mu.Unlock()

	delete(t.sessions, id)
}

func (t *activityTracker) snapshot() map[int]sessionInfo {
	t.mu.RLock()
	defer t.mu.RUnlock()

	out := make(map[int]sessionInfo)

	for k, v := range t.sessions {
		out[k] = v
	}

	return out
}

// AutoLogout mirrors JWTProtected()'s style: no args, drop-in middleware.
func AutoLogout() fiber.Handler {
	timeout := getAutoLogoutTimeout()

	// start the background scanner exactly once, lazily
	scannerOnce.Do(func() {
		go func() {
			ticker := time.NewTicker(time.Minute)
			defer ticker.Stop()
			for range ticker.C {
				now := time.Now()
				for id, session := range tracker.snapshot() {
					if now.Sub(session.LastSeen) > timeout {
						tracker.remove(id)

						var err error

						if session.Role == "Super-Admin" {
							err = SupAdScript.LogoutSuperAdmin(id)
						} else {
							err = UserScript.LogoutUser(id)
						}

						if err != nil {
							log.Printf("auto-logout: failed for user %d: %v", id, err)
						}
					}
				}
			}
		}()
	})

	return func(c fiber.Ctx) error {
		id, ok := c.Locals("id").(int)
		if !ok {
			return c.Next()
		}

		role, _ := c.Locals("role").(string)

		if session, seen := tracker.get(id); seen &&
			time.Since(session.LastSeen) > timeout {

			tracker.remove(id)

			if role == "Super-Admin" {
				_ = SupAdScript.LogoutSuperAdmin(id)
			} else {
				_ = UserScript.LogoutUser(id)
			}

			return fiber.NewError(
				fiber.StatusUnauthorized,
				"session expired due to inactivity",
			)
		}

		// even if not in tracker (e.g. evicted by the scanner, or server restarted),
		// refuse if the DB says this user is already logged out
		if role == "Super-Admin" {
			loggedIn, err := SupAdScript.IsSuperAdminLoggedIn(id)
			if err == nil && !loggedIn {
				return fiber.NewError(
					fiber.StatusUnauthorized,
					"session expired due to inactivity",
				)
			}
		} else {
			loggedIn, err := UserScript.IsLoggedIn(id)
			if err == nil && !loggedIn {
				return fiber.NewError(
					fiber.StatusUnauthorized,
					"session expired due to inactivity",
				)
			}
		}

		tracker.touch(id, role)
		return c.Next()
	}
}

// in middleware/auto_logout.go
func TouchActivity(id int, role string) {
	tracker.touch(id, role)
}
