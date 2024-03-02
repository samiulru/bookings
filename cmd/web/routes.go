package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/samiulru/bookings/internal/config"
	"github.com/samiulru/bookings/internal/handlers"
)

func routes(app *config.AppConfig) http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	mux.Use(NoSurf)
	mux.Use(SessionLoad)

	mux.Get("/", handlers.Repo.Home)
	mux.Get("/about", handlers.Repo.About)
	mux.Get("/rooms/{id}", handlers.Repo.RoomsHandler)
	mux.Get("/contact", handlers.Repo.Contact)
	mux.Get("/user/login", handlers.Repo.UserLogin)
	mux.Post("/user/login", handlers.Repo.PostUserLogin)
	mux.Get("/user/logout", handlers.Repo.AdminLogout)

	mux.Get("/search-availability", handlers.Repo.SearchAvailability)
	mux.Post("/search-availability", handlers.Repo.PostSearchAvailability)
	mux.Post("/search-availability-json", handlers.Repo.SearchAvailabilityJSON)

	mux.Get("/choose-room/{id}", handlers.Repo.ChooseRoom)
	mux.Get("/book-now", handlers.Repo.BookNow)

	mux.Get("/make-reservation", handlers.Repo.Reservation)
	mux.Post("/make-reservation", handlers.Repo.PostReservation)
	mux.Get("/reservation-summary", handlers.Repo.ReservationSummary)

	//FileServer for serving files
	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	//secure routes that available only for specific user
	mux.Route("/admin", func(mux chi.Router) {
		mux.Use(Auth)

		mux.Get("/dashboard", handlers.Repo.AdminDashboard)
		mux.Get("/new-reservations", handlers.Repo.AdminNewReservations)
		mux.Get("/all-reservations", handlers.Repo.AdminAllReservations)
		mux.Get("/reservations-calender", handlers.Repo.AdminReservationsCalender)
		mux.Post("/reservations-calender", handlers.Repo.AdminPostReservationsCalender)

		mux.Get("/reservations/{src}/{id}/show", handlers.Repo.AdminShowReservation)
		mux.Post("/reservations/{src}/{id}", handlers.Repo.AdminPostShowReservation)
		mux.Get("/process-reservation/{src}/{id}/do", handlers.Repo.AdminProcessReservation)
		mux.Get("/delete-reservation/{src}/{id}/do", handlers.Repo.AdminDeleteReservation)

		mux.Get("/rooms/show-all-room", handlers.Repo.AdminShowAllRooms)
		mux.Get("/rooms/{id}/delete", handlers.Repo.AdminDeleteRoom)
		mux.Get("/rooms/{id}/edit", handlers.Repo.AdminEditRoom)
		mux.Post("/rooms/{id}/edit", handlers.Repo.AdminPostEditRoom)
		mux.Get("/rooms/add-new-room", handlers.Repo.AdminAddNewRoom)
		mux.Post("/rooms/add-new-room", handlers.Repo.AdminPostAddNewRoom)

	})
	return mux
}
