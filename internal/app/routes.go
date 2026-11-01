package app

import (
	"net/http"

	"github.com/puppe1990/amarra-cais/pkg/amarra/live"
	"github.com/puppe1990/amarra-cais/pkg/cais"
	"github.com/puppe1990/amarra-cais/pkg/cais/middleware"

	"github.com/puppe1990/aura-ciclos/internal/handlers"
)

func registerRoutes(r *cais.Router, deps Deps, cfg cais.Config) {
	home := handlers.NewHomeHandler(deps.Views, deps.Site, deps.Catalog, cfg)
	dashboard := handlers.NewDashboardHandler(deps.Views, deps.Store, deps.Site, deps.Catalog, cfg)
	auth := handlers.NewAuthHandler(deps.Views, deps.Store, deps.Site, deps.Store.Sessions(), cfg, deps.Catalog)
	tracker := handlers.NewTrackerHandler(deps.Views, deps.Store, deps.Site, deps.Catalog, cfg)

	loginLimit := middleware.NewRateLimiter(10, cfg)
	resetLimit := middleware.NewRateLimiter(10, cfg)

	r.Get("/", home.ServeHTTP)
	r.Get("/login", auth.Login)
	r.Post("/login", loginLimit.Middleware(http.HandlerFunc(auth.LoginPost)).ServeHTTP)
	r.Get("/signup", auth.SignUp)
	r.Post("/signup", loginLimit.Middleware(http.HandlerFunc(auth.SignUpPost)).ServeHTTP)
	r.Get("/anonimo", auth.Anonimo)
	r.Post("/anonimo", loginLimit.Middleware(http.HandlerFunc(auth.AnonimoPost)).ServeHTTP)
	r.Get("/forgot-password", auth.ForgotPassword)
	r.Post("/forgot-password", resetLimit.Middleware(http.HandlerFunc(auth.ForgotPasswordPost)).ServeHTTP)
	r.Get("/reset-password", auth.ResetPassword)
	r.Post("/reset-password", resetLimit.Middleware(http.HandlerFunc(auth.ResetPasswordPost)).ServeHTTP)
	r.Post("/logout", auth.LogoutPost)
	r.Post("/locale", handlers.PostLocale(cfg))
	r.Get("/dashboard", middleware.RequireAuthFunc("/login", dashboard.ServeHTTP))
	r.Get("/settings", middleware.RequireAuthFunc("/login", auth.Settings))
	r.Post("/settings", resetLimit.Middleware(http.HandlerFunc(auth.SettingsPost)).ServeHTTP)

	r.Get("/ciclo", middleware.RequireAuthFunc("/login", tracker.Ciclo))
	r.Post("/ciclo/periodo", middleware.RequireAuthFunc("/login", tracker.PeriodoPost))
	r.Get("/registro", middleware.RequireAuthFunc("/login", tracker.Registro))
	r.Post("/registro", middleware.RequireAuthFunc("/login", tracker.RegistroPost))
	r.Get("/registro/vida-sexual", middleware.RequireAuthFunc("/login", tracker.VidaSexual))
	r.Post("/registro/vida-sexual", middleware.RequireAuthFunc("/login", tracker.VidaSexualPost))
	r.Get("/relatorios", middleware.RequireAuthFunc("/login", tracker.Relatorios))
	r.Post("/modo", middleware.RequireAuthFunc("/login", tracker.ModoPost))
	r.Get("/apoio", middleware.RequireAuthFunc("/login", tracker.Apoio))
	r.Post("/apoio/maya", middleware.RequireAuthFunc("/login", tracker.MayaPost))
	r.Get("/apoio/{slug}", middleware.RequireAuthFunc("/login", handlers.ArtigoRoute(tracker)))
	r.Post("/lembretes", middleware.RequireAuthFunc("/login", tracker.LembretePost))
}

func registerLiveViews(hub *live.Hub, deps Deps) {
	_ = hub
	_ = deps
	// cais:live-views
}
