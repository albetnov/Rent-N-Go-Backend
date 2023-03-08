package utils

import (
	"github.com/gofiber/fiber/v2"
	"math"
	"strconv"
)

// Wrap
// Directly access Wrapper.Wrap method
func Wrap(data fiber.Map, deps ...interface{}) *Wrapper {
	wrapper := Wrapper{}

	if len(deps) > 0 {
		if deps[0] != nil {
			wrapper.c = deps[0].(*fiber.Ctx)
		}

		if deps[1] != nil {
			wrapper.sess = deps[1].(SessionStore)
		}
	}

	return wrapper.Wrap(data)
}

type Wrapper struct {
	data fiber.Map
	c    *fiber.Ctx
	sess SessionStore
}

// Wrap
// Wrap the given payload
func (w *Wrapper) Wrap(data fiber.Map) *Wrapper {
	w.data = data
	return w
}

// Validation
// Wrap your fiber map with validation errors. Need sess, deps[1].
func (w *Wrapper) Validation() *Wrapper {
	err, validation := GetFailedValidation(w.sess)
	w.data["_Validation"] = validation
	w.data["_Error"] = err

	return w
}

// Message
// Need sess, deps[1] parameter. Wrap your response with message component compliance.
func (w *Wrapper) Message() *Wrapper {
	w.data["_Message"] = w.sess.GetFlash("message")
	return w
}

// Error
// Need sess, deps[1] parameter. Wrap your response with error component compliance.
func (w *Wrapper) Error() *Wrapper {
	w.data["_Error"] = w.sess.GetFlash("error")
	return w
}

// Pagination
// Wrap your fiber map with data necessary to use pagination component need ctx, deps[0].
func (w *Wrapper) Pagination(totalRecord int64) *Wrapper {
	page, err := strconv.Atoi(w.c.Query("page", PAGE_DEFAULT))

	if err != nil {
		page = 1
	}

	pageSize, err := strconv.Atoi(w.c.Query("page_size", PAGING_SIZE))

	if err != nil {
		pageSize = 5
	}

	w.data["_pagingTotal"] = int(math.Ceil(float64(totalRecord) / float64(pageSize)))
	w.data["_pagingCurrent"] = page

	return w
}

// Search
// Wrap your fiber data with search value
func (w *Wrapper) Search(search string) *Wrapper {
	w.data["_search"] = search
	return w
}

// Csrf
// Wrap your fiber data with csrf, need Ctx, deps[0].
func (w *Wrapper) Csrf() *Wrapper {
	w.data["_csrf"] = w.c.Locals("token")
	return w
}

// Get
// Return the wrapper data
func (w Wrapper) Get() fiber.Map {
	return w.data
}
