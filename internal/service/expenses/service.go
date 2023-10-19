package expenses_service

import (
	"context"
	"errors"
	"time"
	//category_service "github.com/ozonmp/week-3-workshop/category-service/pkg/category-service"
)

var ErrWrongCategory = errors.New("category does not exist")

type IRepository interface {
	SaveExpense(ctx context.Context, expense *Expense) error
}

// type ICategoryClient interface {
// 	IsCategoryExists(ctx context.Context, categoryID int64) (ok bool, err error)
// }

type Service struct {
	repo IRepository
	//client ICategoryClient
}

func NewService() *Service {
	return &Service{
		repo: newRepo(),
		//client: newClient(grpcClient),
	}
}

func (s *Service) AddExpense(
	ctx context.Context,
	amount float64,
	comment string,
	categoryID uint8,
	userID uint64,
) (*Expense, error) {
	//exists, err := s.client.IsCategoryExists(ctx, categoryID)
	// if err != nil {
	// 	return nil, err
	// }

	// if !exists {
	// 	return nil, ErrWrongCategory
	// }

	expense := &Expense{
		Amount:     amount,
		Timestamp:  time.Now(),
		Comment:    comment, // TODO: optional
		CategoryID: categoryID,
		UserID:     userID,
	}

	if err := s.repo.SaveExpense(ctx, expense); err != nil {
		return nil, err
	}

	return expense, nil
}
