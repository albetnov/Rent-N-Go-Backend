package utils

import (
	"github.com/gofiber/fiber/v2"
	"math"
	"reflect"
	"rent-n-go-backend/models"
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

		if len(deps) > 1 && deps[1] != nil {
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

func getPagination(ctx *fiber.Ctx, totalRecord int64) (int, int) {
	page, err := strconv.Atoi(ctx.Query("page", PAGE_DEFAULT))

	if err != nil {
		page = 1
	}

	pageSize, err := strconv.Atoi(ctx.Query("page_size", PAGING_SIZE))

	if err != nil {
		pageSize = 5
	}

	return int(math.Ceil(float64(totalRecord) / float64(pageSize))), page
}

// Pagination
// Wrap your fiber map with data necessary to use pagination component need ctx, deps[0].
func (w *Wrapper) Pagination(totalRecord int64) *Wrapper {
	totalPage, page := getPagination(w.c, totalRecord)

	w.data["_pagingTotal"] = totalPage
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

func (w *Wrapper) WithMeta(totalRecord int64) *Wrapper {
	totalPage, currentPage := getPagination(w.c, totalRecord)

	w.data["meta"] = fiber.Map{
		"total_page":   totalPage,
		"current_page": currentPage,
		"has_previous": currentPage > 1,
		"has_next":     currentPage < totalPage,
	}

	return w
}

func processItem(c *fiber.Ctx, item reflect.Value) ([]fiber.Map, []fiber.Map) {
	var features []fiber.Map
	var pictures []fiber.Map

	for _, p := range item.FieldByName("Pictures").Interface().([]models.Pictures) {
		pictures = append(pictures, fiber.Map{
			"file_name": FormatUrl(c, p.FileName, p.Associate),
		})
	}

	return features, pictures
}

func MapToServiceableSingle[T comparable](
	c *fiber.Ctx,
	data T,
	callback func(data T, features, pictures []fiber.Map) fiber.Map) fiber.Map {
	arrOfData := reflect.Indirect(reflect.ValueOf(data))

	features, pictures := processItem(c, arrOfData)

	return callback(data, features, pictures)
}

func MapToServiceable[T comparable](
	c *fiber.Ctx,
	data []T,
	callback func(data T, features, pictures []fiber.Map) fiber.Map) []fiber.Map {
	var result []fiber.Map

	arrOfData := reflect.ValueOf(data)

	for i, v := range data {
		features, pictures := processItem(c, reflect.Indirect(arrOfData.Index(i)))

		result = append(result, callback(v, features, pictures))
	}

	return result
}
