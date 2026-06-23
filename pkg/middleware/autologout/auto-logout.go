package middleware

import (
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	"ideyanale-be/pkg/modules/users/script"

	"github.com/gofiber/fiber/v3"
)

type activityTracker struct {
	mu       sync.RWMutex
	lastSeen map[int]time.Time
}

var (
	tracker     = &activityTracker{lastSeen: make(map[int]time.Time)}
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

func (t *activityTracker) touch(id int) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.lastSeen[id] = time.Now()
}

func (t *activityTracker) lastSeenAt(id int) (time.Time, bool) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	ts, ok := t.lastSeen[id]
	return ts, ok
}

func (t *activityTracker) remove(id int) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.lastSeen, id)
}

func (t *activityTracker) snapshot() map[int]time.Time {
	t.mu.RLock()
	defer t.mu.RUnlock()
	out := make(map[int]time.Time, len(t.lastSeen))
	for k, v := range t.lastSeen {
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
				for userID, last := range tracker.snapshot() {
					if now.Sub(last) > timeout {
						tracker.remove(userID)
						if err := script.LogoutUser(userID); err != nil {
							log.Printf("auto-logout: failed for user %d: %v", userID, err)
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

		if last, seen := tracker.lastSeenAt(id); seen && time.Since(last) > timeout {
			tracker.remove(id)
			_ = script.LogoutUser(id)
			return fiber.NewError(fiber.StatusUnauthorized, "session expired due to inactivity")
		}

		tracker.touch(id)
		return c.Next()
	}
}

// in middleware/auto_logout.go
func TouchActivity(id int) {
	tracker.touch(id)
}