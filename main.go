package main

import (
	"embed"
	"log/slog"
	"mouji/commons/auth"
	"mouji/commons/sqlite"
	"mouji/commons/templates"
	"mouji/features/home"
	"mouji/features/login"
	"mouji/features/pageviews"
	"mouji/features/projects"
	"mouji/features/settings"
	"mouji/features/users"
	"net/http"
	"os"
)

//go:embed all:commons all:features
var resources embed.FS

//go:embed assets/*
var assets embed.FS

//go:embed migrations/*.sql
var migrations embed.FS

func main() {
	defer func() {
		if r := recover(); r != nil {
			slog.Error("killing server", "error", r)
			os.Exit(1)
		}
	}()

	sqlite.NewDB()
	defer sqlite.DB.Close()

	sqlite.Migrate(migrations)

	templates.NewTemplates(resources)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	port = ":" + port

	slog.Info("starting server", "port", port)
	err := http.ListenAndServe(port, newRouter())
	if err != nil {
		panic(err)
	}
}

func newRouter() *http.ServeMux {
	mux := http.NewServeMux()

	// public
	mux.Handle("GET /assets/", http.FileServer(http.FS(assets)))
	mux.HandleFunc("GET /collect", pageviews.HandleCollect)
	mux.HandleFunc("GET /login", login.HandleLoginPage)
	mux.HandleFunc("POST /login", login.HandleLoginSubmit)

	// private
	addPrivateRoute(mux, "GET /", home.HandleHomePage)
	addPrivateRoute(mux, "GET /settings", settings.HandleSettingsPage)
	addPrivateRoute(mux, "GET /users/new", users.HandleNewUserPage)
	addPrivateRoute(mux, "POST /users/new", users.HandleNewUserSubmit)
	addPrivateRoute(mux, "GET /projects/new", projects.HandleNewProjectPage)
	addPrivateRoute(mux, "GET /projects/{project_id}", projects.HandleEditProjectPage)
	addPrivateRoute(mux, "POST /projects/", projects.HandleProjectDetailSubmit)
	addPrivateRoute(mux, "POST /projects/{project_id}", projects.HandleProjectDetailSubmit)

	return mux
}

func addPrivateRoute(mux *http.ServeMux, pattern string, handlerFunc func(w http.ResponseWriter, r *http.Request)) {
	handler := http.HandlerFunc(handlerFunc)
	mux.HandleFunc(pattern, auth.EnsureAuthenticated(handler))
}
