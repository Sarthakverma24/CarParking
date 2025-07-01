package main

import (
	"CarParking/db"
	"CarParking/wrapper"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	loggers := db.Logger.Sugar()

	postgresSQLWrapper := db.NewPostgresSQLWrapper()
	if postgresSQLWrapper == nil {
		loggers.Fatal("Failed to initialize PostgresSQL wrapper")
	}
	loggers.Infoln("✅ Successfully initialized PostgresSQL connection")

	userWrapper := wrapper.UserWrapper{Db: postgresSQLWrapper}
	lotWrapper := wrapper.LotWrapper{Db: postgresSQLWrapper}

	r := mux.NewRouter()

	// Register APIs
	RegisterUserRoutes(r, &userWrapper)
	RegisterLotRoutes(r, &lotWrapper)

	// Add CORS middleware
	handlerWithCORS := corsMiddleware(r)

	// Start server
	if err := http.ListenAndServe(":8080", handlerWithCORS); err != nil {
		log.Fatal("Server error:", err)
	}
}

func RegisterUserRoutes(r *mux.Router, userWrapper *wrapper.UserWrapper) {
	r.HandleFunc("/api/login", LogInHandler(userWrapper)).Methods("GET", "POST", "OPTIONS")
	r.HandleFunc("/api/signin", SignInHandler(userWrapper)).Methods("GET", "POST", "OPTIONS")
	r.HandleFunc("/api/Users", UserHandler(userWrapper)).Methods("GET", "PUT", "OPTIONS")
}

func RegisterLotRoutes(r *mux.Router, lotWrapper *wrapper.LotWrapper) {
	r.HandleFunc("/api/searchSlots", UserLot(lotWrapper)).Methods("GET", "OPTIONS")
	r.HandleFunc("/api/UserHistory/{username}", UserBookings(lotWrapper)).Methods("GET", "OPTIONS")
	r.HandleFunc("/api/bookSlot", UserBook(lotWrapper)).Methods("GET", "POST", "OPTIONS")
	r.HandleFunc("/api/releaseSlot", ReleaseSlotHandler(lotWrapper)).Methods("PUT", "OPTIONS")

	r.HandleFunc("/api/summary", GetSummaryHandler(lotWrapper)).Methods("GET", "OPTIONS")
	r.HandleFunc("/api/used-summary", UsedSummaryHandler(lotWrapper)).Methods("GET", "OPTIONS")

	r.HandleFunc("/api/login/addlot", AddLotHandler(lotWrapper)).Methods("GET", "POST", "OPTIONS")
	r.HandleFunc("/api/lots", LotHandler(lotWrapper)).Methods("GET", "POST", "OPTIONS")
	r.HandleFunc("/api/lots/{lot_id}", EditLotHandler(lotWrapper)).Methods("GET", "PUT", "OPTIONS")
	r.HandleFunc("/api/lots/spots/{lot_id}", SpotHandler(lotWrapper)).Methods("GET", "PUT", "POST", "OPTIONS")
	r.HandleFunc("/api/lot/{lot_id}", DeleteLot(lotWrapper)).Methods("DELETE", "OPTIONS")
	r.HandleFunc("/api/slots/{parking_id}", DeleteSLot(lotWrapper)).Methods("DELETE", "OPTIONS")
}
