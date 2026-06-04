package api

import (
	"ci-jarvis/internal/store/postgres"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type RunHandler struct {
	db *postgres.DB
}

func NewRunHandler(db *postgres.DB) *RunHandler {
	return &RunHandler{db: db}
}

func (h *RunHandler) ListRuns(c *fiber.Ctx) error {
	limitStr := c.Query("limit", "10")
	offsetStr := c.Query("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 10
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		offset = 0
	}

	runs, err := h.db.ListRuns(c.Context(), limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to fetch runs",
		})
	}

	return c.JSON(fiber.Map{
		"runs":   runs,
		"limit":  limit,
		"offset": offset,
	})
}

func (h *RunHandler) GetRun(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "missing run id",
		})
	}

	run, err := h.db.GetRunById(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "run not found",
		})
	}

	return c.JSON(run)
}
