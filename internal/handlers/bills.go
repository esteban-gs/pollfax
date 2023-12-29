package handlers

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"pollfax/db"
	"pollfax/internal/dto"
)

func GetAll(c echo.Context) error {
	_db := db.Instance()
	bills := []dto.BillRes{}
	_db.Select(&bills, `SELECT title,
                              type,
													    bill_number,
													    origin_chamber,
													    url,
													    latest_action_text,
													    update_including_text
									           FROM bills`)

	return c.JSON(http.StatusOK, bills)
}

func GetAllCache(c echo.Context) error {
	return c.File("../bills.json")
}
