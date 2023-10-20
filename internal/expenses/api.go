package expenses

import (
	"net/http"
	"strconv"

	"github.com/boskuv/finance-control_expenses-service/internal/errors"
	"github.com/boskuv/finance-control_expenses-service/pkg/log"
	"github.com/boskuv/finance-control_expenses-service/pkg/pagination"
	routing "github.com/go-ozzo/ozzo-routing/v2"
)

// RegisterHandlers sets up the routing of the HTTP handlers.
func RegisterHandlers(r *routing.RouteGroup, service Service, logger log.Logger) { // authHandler routing.Handler
	res := resource{service, logger}

	r.Get("/expenses/<id>", res.get)
	r.Get("/expenses", res.query)

	//r.Use(authHandler)

	// the following endpoints require a valid JWT
	r.Post("/expenses", res.create)
	r.Put("/expenses/<id>", res.update)
	r.Delete("/expenses/<id>", res.delete)
}

type resource struct {
	service Service
	logger  log.Logger
}

func (r resource) get(c *routing.Context) error {
	id, _ := strconv.ParseUint(string(c.Param("id")), 10, 64) // TODO
	expense, err := r.service.Get(c.Request.Context(), id)
	if err != nil {
		return err
	}

	return c.Write(expense)
}

func (r resource) query(c *routing.Context) error {
	ctx := c.Request.Context()
	//count, err := r.service.Count(ctx) // TODO
	// if err != nil {
	// 	return err
	// }
	count := 3
	pages := pagination.NewFromRequest(c.Request, count)
	expenses, err := r.service.Query(ctx, pages.Offset(), pages.Limit())
	if err != nil {
		return err
	}
	pages.Items = expenses
	return c.Write(pages)
}

func (r resource) create(c *routing.Context) error {
	var input CreateExpenseRequest
	if err := c.Read(&input); err != nil {
		r.logger.With(c.Request.Context()).Info(err)
		return errors.BadRequest("")
	}
	expense, err := r.service.Create(c.Request.Context(), input)
	if err != nil {
		return err
	}

	return c.WriteWithStatus(expense, http.StatusCreated)
}

func (r resource) update(c *routing.Context) error {
	var input UpdateExpenseRequest
	if err := c.Read(&input); err != nil {
		r.logger.With(c.Request.Context()).Info(err)
		return errors.BadRequest("")
	}

	id, _ := strconv.ParseUint(string(c.Param("id")), 10, 64) // TODO
	expense, err := r.service.Update(c.Request.Context(), id, input)
	if err != nil {
		return err
	}

	return c.Write(expense)
}

func (r resource) delete(c *routing.Context) error {
	id, _ := strconv.ParseUint(string(c.Param("id")), 10, 64) // TODO
	expense, err := r.service.Delete(c.Request.Context(), id)
	if err != nil {
		return err
	}

	return c.Write(expense)
}
