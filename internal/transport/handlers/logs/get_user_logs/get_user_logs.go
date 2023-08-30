package get_user_logs

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/kiryu-dev/segments-api/internal/transport/handlers"
)

type logsGetter interface {
	GetUserLogs(context.Context, uint64, time.Time) (string, error)
}

// GetUserLogs godoc
//
//	@Summary		Получить историю изменения сегментов пользователя
//	@Description	Получение истории добавления и удаления сегментов указанного пользователя за определенный год и месяц в формате CSV.
//	@Tags			logs
//	@Produce		octet-stream
//	@Param			userID	path	int		true	"user id"
//	@Param			date	query	string	false	"filter date"	Format(year-month)
//	@Success		200
//	@Failure		400		{object}	handlers.responseError	"error"
//	@Failure		500		{object}	handlers.responseError	"error"
//	@Failure		default	{object}	handlers.responseError	"error"
//	@Router			/log/{userID} [get]
func New(service logsGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		idStr := mux.Vars(r)["userID"]
		userID, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			handlers.WriteJSONError(w, http.StatusBadRequest, "invalid user id")
			return
		}
		filterDate, err := getFilterDate(r.URL.Query())
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			handlers.WriteJSONError(w, http.StatusBadRequest, err.Error())
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()
		filename, err := service.GetUserLogs(ctx, userID, filterDate)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			handlers.WriteServerError(w, http.StatusInternalServerError)
			return
		}
		defer os.Remove(filename)
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s",
			strconv.Quote(filename)))
		w.Header().Set("Content-Type", "application/octet-stream")
		http.ServeFile(w, r, filename)
	}
}

func getFilterDate(queries url.Values) (time.Time, error) {
	if !queries.Has("date") {
		return time.Time{}, fmt.Errorf("enter the date in year-month format")
	}
	filterDate, err := time.Parse("2006-1", queries.Get("date"))
	if err != nil {
		return time.Time{}, fmt.Errorf("enter the date in year-month format")
	}
	return filterDate, nil
}
