package main

import (
	"CarParking/db"
	"CarParking/loggers"
	"CarParking/wrapper"
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	_ "go.uber.org/zap"
	"log"
	_ "log"
	"net/http"
	"strconv"
	"time"
)

var logger = loggers.Logger

type BookingRequest struct {
	LotID     int    `json:"lot_id"`
	SpotID    int    `json:"spot_id"`
	UserID    string `json:"user_id"`
	VehicleNo string `json:"vehicle_no"`
	Location  string `json:"location"`
}

func enableCORS(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		enableCORS(&w)
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func LogInHandler(userWrapper *wrapper.UserWrapper) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		enableCORS(&w)
		w.Header().Set("Content-Type", "application/json")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		logger.Info("✅ login called")

		var user wrapper.User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			loggers.Sugar().Errorw("Invalid login request", "error", err)
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}
		loggers.Sugar().Infow("Login request", "user", user)
		username, err := userWrapper.Login(r.Context(), user)
		if err != nil {
			loggers.Sugar().Warnw("Invalid credentials", "error", err)
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
			return
		}

		loggers.Sugar().Infow("Login successful", "username", username)
		json.NewEncoder(w).Encode(map[string]string{
			"username": username,
		})
	}
}
func SignInHandler(userWrapper *wrapper.UserWrapper) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		enableCORS(&w)
		w.Header().Set("Content-Type", "application/json")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		loggers.Sugar().Info("✅ sign in called")

		var user wrapper.SigninUser
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			loggers.Sugar().Errorw("Invalid signin request", "error", err)
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		username, err := userWrapper.SignIn(r.Context(), user)
		if err != nil {
			loggers.Sugar().Errorw("Failed to sign in", "error", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		loggers.Sugar().Infow("Sign in successful", "username", username)
		json.NewEncoder(w).Encode(map[string]string{
			"username": username,
		})
	}
}

func AddLotHandler(lotWrapper *wrapper.LotWrapper) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		enableCORS(&w)
		w.Header().Set("Content-Type", "application/json")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		loggers.Sugar().Info("✅ ADD LOT called")

		var lot wrapper.Lot
		if err := json.NewDecoder(r.Body).Decode(&lot); err != nil {
			loggers.Sugar().Errorw("Invalid lot request", "error", err)
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		lotID, err := lotWrapper.SetLot(r.Context(), lot)
		if err != nil {
			loggers.Sugar().Errorw("Failed to add lot", "error", err)
			http.Error(w, "Invalid details", http.StatusUnauthorized)
			return
		}

		loggers.Sugar().Infow("Lot added", "lot_id", lotID)
		json.NewEncoder(w).Encode(map[string]string{
			"lot_id": lotID,
		})
	}
}

func LotHandler(lotWrapper *wrapper.LotWrapper) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		enableCORS(&w)
		w.Header().Set("Content-Type", "application/json")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		loggers.Sugar().Info("✅ lot handler called")

		lots, err := lotWrapper.GetLot(r.Context())
		if err != nil {
			loggers.Sugar().Errorw("Failed to fetch lots", "error", err)
			http.Error(w, "Failed to fetch lots", http.StatusInternalServerError)
			return
		}
		log.Println(lots)

		loggers.Sugar().Infow("Lots fetched", "count", len(lots))
		if err := json.NewEncoder(w).Encode(lots); err != nil {
			loggers.Sugar().Errorw("Failed to encode lots", "error", err)
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
	}
}

func UserHandler(userWrapper *wrapper.UserWrapper) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		enableCORS(&w)
		w.Header().Set("Content-Type", "application/json")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		loggers.Sugar().Info("✅ Users handler called")

		users, err := userWrapper.GetUsers(r.Context())
		if err != nil {
			loggers.Sugar().Errorw("Failed to fetch users", "error", err)
			http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
			return
		}

		loggers.Sugar().Infow("Users fetched", "count", len(users))
		if err := json.NewEncoder(w).Encode(users); err != nil {
			loggers.Sugar().Errorw("Failed to encode users", "error", err)
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
	}
}

func EditLotHandler(lotWrapper *wrapper.LotWrapper) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		enableCORS(&w)
		w.Header().Set("Content-Type", "application/json")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		loggers.Sugar().Info("✅ EditLotHandler called")

		vars := mux.Vars(r)
		id := vars["lot_id"]

		lotID, err := strconv.Atoi(id)
		if err != nil {
			loggers.Sugar().Errorw("Invalid lot_id", "lot_id", id)
			http.Error(w, "Invalid lot_id", http.StatusBadRequest)
			return
		}

		var lot wrapper.Lot
		if err := json.NewDecoder(r.Body).Decode(&lot); err != nil {
			loggers.Sugar().Errorw("Invalid request body", "error", err)
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		lot.LotID = lotID

		loggers.Sugar().Infow("Editing lot", "lot", lot)

		_, err = lotWrapper.EditLot(r.Context(), lot)
		if err != nil {
			loggers.Sugar().Errorw("Failed to update lot", "error", err)
			http.Error(w, "Failed to update lot", http.StatusInternalServerError)
			return
		}

		if err := json.NewEncoder(w).Encode(lot); err != nil {
			loggers.Sugar().Errorw("Failed to encode response", "error", err)
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
	}
}

func SpotHandler(lotWrapper *wrapper.LotWrapper) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		enableCORS(&w)
		w.Header().Set("Content-Type", "application/json")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		loggers.Sugar().Info("✅ SpotHandler called")

		vars := mux.Vars(r)
		id := vars["lot_id"]

		lotID, err := strconv.Atoi(id)
		if err != nil {
			loggers.Sugar().Errorw("Invalid lot_id", "lot_id", id)
			http.Error(w, "Invalid lot_id", http.StatusBadRequest)
			return
		}

		if r.Method == http.MethodGet {
			loggers.Sugar().Infow("Fetching spots for lot", "lot_id", lotID)
			spots, err := lotWrapper.GetSpotsByLotID(r.Context(), lotID)
			if err != nil {
				loggers.Sugar().Errorw("Failed to get spots", "error", err)
				http.Error(w, "Failed to get spots", http.StatusInternalServerError)
				return
			}
			if err := json.NewEncoder(w).Encode(spots); err != nil {
				loggers.Sugar().Errorw("Failed to encode response", "error", err)
				http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			}
			return
		}

		// For PUT or POST, expect a request body
		var lot wrapper.Lot
		if err := json.NewDecoder(r.Body).Decode(&lot); err != nil {
			loggers.Sugar().Errorw("Invalid request body", "error", err)
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		lot.LotID = lotID

		loggers.Sugar().Infow("Updating lot spots", "lot_id", lotID)

		// Here you would handle PUT/POST logic — this part depends on your app
		// For now, just respond with OK
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"status": "updated successfully",
		})
	}
}

func DeleteLot(lotWrapper *wrapper.LotWrapper) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		enableCORS(&w)
		w.Header().Set("Content-Type", "application/json")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		loggers.Sugar().Info("✅ Delete called")

		vars := mux.Vars(r)
		id := vars["lot_id"]
		ID, err := strconv.Atoi(id)
		if err != nil {
			loggers.Sugar().Errorw("Invalid lot_id", "lot_id", id)
			http.Error(w, "Invalid lot_id", http.StatusBadRequest)
			return
		}

		loggers.Sugar().Infow("Deleting lot", "lot_id", ID)

		_, err = lotWrapper.DeleteLotByID(r.Context(), ID)
		if err != nil {
			loggers.Sugar().Errorw("Failed to get spots", "error", err)
			http.Error(w, "Failed to get spots", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"status": "deleted successfully",
		})
	}
}

func DeleteSLot(lotWrapper *wrapper.LotWrapper) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		enableCORS(&w)
		w.Header().Set("Content-Type", "application/json")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		loggers.Sugar().Info("✅ Delete slot called")

		vars := mux.Vars(r)
		id := vars["parking_id"]
		ID, err := strconv.Atoi(id)
		if err != nil {
			loggers.Sugar().Errorw("Invalid lot_id", "lot_id", id)
			http.Error(w, "Invalid lot_id", http.StatusBadRequest)
			return
		}

		loggers.Sugar().Infow("Deleting lot", "lot_id", ID)

		_, err = lotWrapper.DeleteSLots(r.Context(), ID)
		if err != nil {
			loggers.Sugar().Errorw("Failed to get spots", "error", err)
			http.Error(w, "Failed to get spots", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"status": "deleted successfully",
		})
	}
}

func UserLot(lotWrapper *wrapper.LotWrapper) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		enableCORS(&w)
		w.Header().Set("Content-Type", "application/json")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		loggers.Sugar().Info("✅ UserLot called")

		// ✅ Read query parameters
		searchType := r.URL.Query().Get("type")
		value := r.URL.Query().Get("value")

		if searchType == "" || value == "" {
			loggers.Sugar().Warn("Missing query parameters")
			http.Error(w, "Missing query parameters", http.StatusBadRequest)
			return
		}

		loggers.Sugar().Infow("Fetching slots", "type", searchType, "value", value)

		// ✅ You can modify this part to search based on `type` (e.g., location, pincode)
		slots, err := lotWrapper.SearchLots(r.Context(), searchType, value)
		if err != nil {
			loggers.Sugar().Errorw("Failed to get slots", "error", err)
			http.Error(w, "Failed to get slots", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(slots)
	}
}

func UserBookings(lotWrapper *wrapper.LotWrapper) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		enableCORS(&w)
		w.Header().Set("Content-Type", "application/json")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		loggers.Sugar().Info("✅ Delete slot called")

		vars := mux.Vars(r)
		username := vars["username"]

		loggers.Sugar().Infow("Booking Spot", "username", username)

		bookings, err := lotWrapper.GetBookings(r.Context(), username)
		if err != nil {
			loggers.Sugar().Errorw("Failed to get spots", "error", err)
			http.Error(w, "Failed to get spots", http.StatusInternalServerError)
			return
		}
		log.Println(bookings)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(bookings)
	}
}

func UserBook(lotWrapper *wrapper.LotWrapper) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		enableCORS(&w)
		w.Header().Set("Content-Type", "application/json")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		var req BookingRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		log.Println(req)
		booking := wrapper.ReservedSpot{
			ID:          uuid.New().String(),
			SpotID:      req.SpotID,
			UserID:      req.UserID,
			VehicleNo:   req.VehicleNo,
			Location:    req.Location,
			Parking:     time.Now(),
			Leaving:     time.Now(), // You can change this to future time later
			ParkingCost: 0,
			Status:      true,
		}

		if err := lotWrapper.InsertBooking(r.Context(), booking); err != nil {
			http.Error(w, "Failed to book slot", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(map[string]string{"message": "Slot booked successfully"})
	}
}

func ReleaseSlotHandler(lw *wrapper.LotWrapper) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		enableCORS(&w)
		w.Header().Set("Content-Type", "application/json")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		if r.Method != http.MethodPut {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var booking wrapper.ReservedSpot
		err := json.NewDecoder(r.Body).Decode(&booking)
		if err != nil {
			db.Logger.Error("Failed to decode request body", zap.Error(err))
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		updatedBooking, err := lw.ReleaseSpot(context.Background(), booking)
		if err != nil {
			db.Logger.Error("Failed to release spot", zap.Error(err))
			http.Error(w, "Failed to release spot", http.StatusInternalServerError)
			return
		}
		loggers.Sugar().Infow("Released slot", "booking", updatedBooking)
		// ✅ Respond with the updated ReservedSpot
		json.NewEncoder(w).Encode(updatedBooking)
	}
}

func GetSummaryHandler(lw *wrapper.LotWrapper) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		enableCORS(&w)
		w.Header().Set("Content-Type", "application/json")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		summary, err := lw.GetSummaryData(context.Background())
		if err != nil {
			http.Error(w, "Failed to fetch summary", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(summary)
	}
}

func UsedSummaryHandler(lw *wrapper.LotWrapper) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		enableCORS(&w)
		w.Header().Set("Content-Type", "application/json")

		summary, err := lw.GetUsedParkingSummary(context.Background())
		if err != nil {
			http.Error(w, "Failed to get usage summary", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(summary)
	}
}
