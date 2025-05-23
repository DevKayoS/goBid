package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (api *Api) BindRoutes() {
	api.Router.Use(
		middleware.RequestID,
		middleware.Recoverer,
		middleware.DefaultLogger,
		api.Sessions.LoadAndSave,
	)

	// csrfMiddleware := csrf.Protect(
	// 	[]byte(os.Getenv("GOBID_CSRF_KEY")),
	// 	csrf.Secure(false), // DEV ONLY
	// )

	// api.Router.Use(csrfMiddleware)

	api.Router.Route("/api", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			// r.Get("/csrf-token", api.HandleGetCSRFtoken)
			r.Route("/users", func(r chi.Router) {
				r.Post("/sign-up", api.handleSignupUser)
				r.Post("/login", api.handleLoginUser)

				r.Group(func(r chi.Router) {
					r.Use(api.AuthMiddleware)
					r.Post("/logout", api.handleLogoutUser)
				})
			})

			r.Route("/products", func(r chi.Router) {
				r.Group(func(r chi.Router) {
					r.Use(api.AuthMiddleware)
					r.Post("/", api.handleCreateProduct)
					r.Get("/list", api.handleListAvailableProduct)

					r.Get("/ws/subscriber/{product_Id}", api.handleSubscribeUserToAuction)
				})
			})

		})
	})
}
