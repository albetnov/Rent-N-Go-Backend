package utils

import (
	"github.com/gofiber/fiber/v2"
	"math"
	"strconv"
)

// Wrap
// Directly access Wrapper.Wrap method
func Wrap(data fiber.Map) *Wrapper {
	wrapper := Wrapper{}
	return wrapper.Wrap(data)
}

type Wrapper struct {
	data fiber.Map
	c    *fiber.Ctx
}

// Wrap
// Wrap the given payload
func (w *Wrapper) Wrap(data fiber.Map) *Wrapper {
	w.data = data
	return w
}

// Validation
// Wrap your fiber map with validation errors
func (w *Wrapper) Validation(store SessionStore) *Wrapper {
	err, validation := GetFailedValidation(store)
	w.data["Validation"] = validation
	w.data["Error"] = err

	return w
}

// Ctx
// Add context to wrapper (Pagination and Search need this)
func (w *Wrapper) Ctx(c *fiber.Ctx) *Wrapper {
	w.c = c
	return w
}

// Pagination
// Wrap your fiber map with data necessary to use pagination component
func (w *Wrapper) Pagination(totalRecord int64) *Wrapper {
	page, err := strconv.Atoi(w.c.Query("page", PAGE_DEFAULT))

	if err != nil {
		page = 1
	}

	pageSize, err := strconv.Atoi(w.c.Query("page_size", PAGING_SIZE))

	if err != nil {
		pageSize = 5
	}

	w.data["pagingTotal"] = int(math.Ceil(float64(totalRecord) / float64(pageSize)))
	w.data["pagingCurrent"] = page

	return w
}

// Search
// Wrap your fiber data with search value
func (w *Wrapper) Search(search string) *Wrapper {
	w.data["search"] = search
	return w
}

func (w *Wrapper) Csrf() *Wrapper {
	w.data["csrf"] = w.c.Locals("token")
	return w
}

// Get
// Return the wrapper data
func (w Wrapper) Get() fiber.Map {
	return w.data
}
