package seed

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/itzpushan/todo/ent"
	"golang.org/x/crypto/bcrypt"
)

var sampleTitles = []string{
	"Learn Go", "Build REST API", "Read Ent Docs", "Design Schema", "Implement Auth",
	"Test Endpoints", "Refactor Code", "Write Unit Tests", "Document Code", "Review PRs",
	"Deploy to Server", "Fix Bugs", "Setup CI/CD", "Write Seeder", "Create Middleware",
	"Secure Routes", "Write Specs", "Run Benchmarks", "Optimize Queries", "Manage State",
	"Handle Errors", "Add Pagination", "Hash Passwords", "Generate JWT", "Configure Logger",
	"Implement Caching", "Handle Edge Cases", "Improve UX", "Write Blog Post", "Push to Git",
	"Use Go Modules", "Understand Context", "Implement Rate Limiting", "Paginate Results",
	"Create CLI Tool", "Use Goroutines", "Integrate Postgres", "Create Dockerfile", "Run Tests",
	"Analyze Logs", "Refactor Router", "Generate Ent Code", "Design API Contract", "Track Metrics",
}

var sampleDescriptions = []string{
	"This is a short task.",
	"Complete this by EOD.",
	"Optional but useful.",
	"Needs team review.",
	"High priority item.",
	"Low hanging fruit.",
	"Refer to docs.",
	"Ask mentor for input.",
	"Double-check output.",
	"Push to production soon.",
}

type UserSeed struct {
	Name     string
	Email    string
	Password string
}

var usersToSeed = []UserSeed{
	{"Alice", "alice@example.com", "alice123"},
	{"Bob", "bob@example.com", "bob123"},
	{"Charlie", "charlie@example.com", "charlie123"},
	{"Diana", "diana@example.com", "diana123"},
	{"Eve", "eve@example.com", "eve123"},
}

func Run(ctx context.Context, client *ent.Client) error {
	log.Println("Seeding database...")

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	userIDs := []uuid.UUID{}

	for _, u := range usersToSeed {
		hashedPwd, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("failed to hash password for user %s: %w", u.Name, err)
		}

		user, err := client.User.Create().
			SetID(uuid.New()).
			SetName(u.Name).
			SetEmail(u.Email).
			SetPassword(string(hashedPwd)).
			Save(ctx)
		if err != nil {
			return fmt.Errorf("failed to create user %s: %w", u.Name, err)
		}
		userIDs = append(userIDs, user.ID)
	}

	todoBulk := make([]*ent.TodoCreate, 0, 40)

	for range [40]int{} {
		todoBulk = append(todoBulk, client.Todo.Create().
			SetID(uuid.New()).
			SetTitle(sampleTitles[r.Intn(len(sampleTitles))]).
			SetDescription(sampleDescriptions[r.Intn(len(sampleDescriptions))]).
			SetAuthorID(userIDs[r.Intn(len(userIDs))]),
		)
	}

	if _, err := client.Todo.CreateBulk(todoBulk...).Save(ctx); err != nil {
		return fmt.Errorf("failed to create todos: %w", err)
	}

	log.Println("Done seeding database.")
	return nil
}
