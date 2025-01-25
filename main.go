package main

import (
	"embed"
	"log/slog"
	"mouji/commons/auth"
	"mouji/commons/session"
	"mouji/commons/sqlite"
	"mouji/commons/templates"
	"mouji/features/home"
	"mouji/features/login"
	"mouji/features/projects"
	"mouji/features/settings"
	"mouji/features/users"
	"net/http"
	"os"
	"strings"
	"time"
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

	// seed.SeedUsers(resources)
	// seed.SeedProjects(resources)
	// seed.SeedPageViews(resources)

	templates.NewTemplates(resources)

	go runBackgroundTasks()

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
	mux.HandleFunc("GET /assets/", handleStaticAssets)
	// mux.HandleFunc("GET /collect", pageviews.HandleCollect)
	mux.HandleFunc("GET /login", login.HandleLoginPage)
	mux.HandleFunc("POST /login", login.HandleLoginSubmit)

	// private
	addPrivateRoute(mux, "GET /", home.HandleHomePage)
	addPrivateRoute(mux, "GET /settings", settings.HandleSettingsPage)
	addPrivateRoute(mux, "GET /settings/server_url", settings.HandleServerURLPage)
	addPrivateRoute(mux, "POST /settings/server_url", settings.HandleServerURLSubmit)
	// addPrivateRoute(mux, "GET /users/new", users.HandleNewUserPage)
	// addPrivateRoute(mux, "POST /users/new", users.HandleNewUserSubmit)
	addPrivateRoute(mux, "GET /users/me/password", users.HandleChangePasswordPage)
	// addPrivateRoute(mux, "POST /users/me/password", users.HandleChangePasswordSubmit)
	addPrivateRoute(mux, "GET /projects/new", projects.HandleNewProjectPage)
	addPrivateRoute(mux, "GET /projects/{project_id}", projects.HandleEditProjectPage)
	// addPrivateRoute(mux, "POST /projects/", projects.HandleProjectDetailSubmit)
	// addPrivateRoute(mux, "POST /projects/{project_id}", projects.HandleProjectDetailSubmit)

	return mux
}

func handleStaticAssets(w http.ResponseWriter, r *http.Request) {
	if strings.HasSuffix(r.URL.Path, ".woff2") {
		w.Header().Set("Cache-Control", "public, max-age=31536000") // 1 year
	}

	http.FileServer(http.FS(assets)).ServeHTTP(w, r)
}

func addPrivateRoute(mux *http.ServeMux, pattern string, handlerFunc func(w http.ResponseWriter, r *http.Request)) {
	handler := http.HandlerFunc(handlerFunc)
	mux.HandleFunc(pattern, auth.EnsureAuthenticated(handler))
}

func runBackgroundTasks() {
	for range time.Tick(24 * time.Hour) {
		session.DeleteExpiredSessions()
	}
}
