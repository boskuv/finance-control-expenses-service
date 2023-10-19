package expenses_service

import "context"

var nextID uint64 = 1

type repository struct{}

func newRepo() IRepository {
	return repository{}
}

func (r repository) SaveExpense(ctx context.Context, expense *Expense) error {
	expense.ID = nextID

	nextID++

	return nil
}
