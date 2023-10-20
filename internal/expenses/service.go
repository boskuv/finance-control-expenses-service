package expenses

import (
	"context"

	"github.com/boskuv/finance-control_expenses-service/internal/entity"
	"github.com/boskuv/finance-control_expenses-service/pkg/log"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Service interface {
	Get(ctx context.Context, id uint64) (Expense, error)
	Query(ctx context.Context, offset, limit int) ([]Expense, error)
	//Count(ctx context.Context) (int, error)
	Create(ctx context.Context, input CreateExpenseRequest) (Expense, error)
	Update(ctx context.Context, id uint64, input UpdateExpenseRequest) (Expense, error)
	Delete(ctx context.Context, id uint64) (Expense, error)
}

type Expense struct {
	entity.Expense
}

// CreateExpenseRequest represents an Expense creation request.
type CreateExpenseRequest struct {
	Name string `json:"name"`
}

// Validate validates the CreateExpenseRequest fields.
func (m CreateExpenseRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Name, validation.Required, validation.Length(0, 128)),
	)
}

// UpdateExpenseRequest represents an Expense update request.
type UpdateExpenseRequest struct {
	Name string `json:"name"`
}

// Validate validates the CreateExpenseRequest fields.
func (m UpdateExpenseRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Name, validation.Required, validation.Length(0, 128)),
	)
}

type service struct {
	repo   Repository
	logger log.Logger
}

// NewService creates a new Expense service.
func NewService(repo Repository, logger log.Logger) Service {
	return service{repo, logger}
}

// Get returns the Expense with the specified the Expense ID.
func (s service) Get(ctx context.Context, id uint64) (Expense, error) {
	expense, err := s.repo.Get(ctx, id)
	if err != nil {
		return Expense{}, err
	}
	return Expense{expense}, nil
}

// Create creates a new Expense.
func (s service) Create(ctx context.Context, req CreateExpenseRequest) (Expense, error) {
	if err := req.Validate(); err != nil {
		return Expense{}, err
	}
	id := GenerateID()
	//now := time.Now()
	// err := s.repo.Create(ctx, entity.Expense{
	// 	ID:        id,
	// 	Name:      req.Name,
	// 	CreatedAt: now,
	// 	UpdatedAt: now,
	// })
	// if err != nil {
	// 	return Expense{}, err
	// }
	return s.Get(ctx, id)
}

// Update updates the Expense with the specified ID.
func (s service) Update(ctx context.Context, id uint64, req UpdateExpenseRequest) (Expense, error) {
	if err := req.Validate(); err != nil {
		return Expense{}, err
	}

	Expense, err := s.Get(ctx, id)
	if err != nil {
		return Expense, err
	}
	// Expense.Name = req.Name
	// Expense.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, Expense.Expense); err != nil {
		return Expense, err
	}
	return Expense, nil
}

// Delete deletes the Expense with the specified ID.
func (s service) Delete(ctx context.Context, id uint64) (Expense, error) {
	expense, err := s.Get(ctx, id)
	if err != nil {
		return Expense{}, err
	}
	if err = s.repo.Delete(ctx, id); err != nil {
		return Expense{}, err
	}
	return expense, nil
}

// Count returns the number of Expenses.
// func (s service) Count(ctx context.Context) (int, error) {
// 	return s.repo.Count(ctx)
// }

// Query returns the Expenses with the specified offset and limit.
func (s service) Query(ctx context.Context, offset, limit int) ([]Expense, error) {
	items, err := s.repo.Query(ctx, offset, limit)
	if err != nil {
		return nil, err
	}
	result := []Expense{}
	for _, item := range items {
		result = append(result, Expense{item})
	}
	return result, nil
}

func GenerateID() uint64 {
	return 1 //uuid.New().String()
}
