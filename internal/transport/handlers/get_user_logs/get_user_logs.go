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
)

type logsGetter interface {
	GetUserLogs(context.Context, uint64, time.Time) (string, error)
}

func New(service logsGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := mux.Vars(r)["userID"]
		userID, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("invalid user id"))
			return
		}
		filterDate, err := getFilterDate(r.URL.Query())
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()
		filename, err := service.GetUserLogs(ctx, userID, filterDate)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
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
