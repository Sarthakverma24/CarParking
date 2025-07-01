package wrapper

import (
	"CarParking/db"
	"context"
	"database/sql"
	"fmt"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"log"
	"math"
	"time"
)

type LotWrapper struct {
	Db db.Database
}

type Lot struct {
	LotID    int    `json:"id"`
	Location string `json:"location"`
	Price    int    `json:"price"`
	Address  string `json:"address"`
	Pincode  string `json:"pincode"`
	Spots    int    `json:"spots"`
	SpotList []Spot `json:"spot_list"`
}

type Spot struct {
	ParkingID int    `json:"parking_id"`
	LotID     int    `json:"lot_id"`
	Status    bool   `json:"status"`
	Address   string `json:"address"`
}
type ReservedSpot struct {
	ID          string    `json:"id"`
	SpotID      int       `json:"spot_id"`
	UserID      string    `json:"user_id"`
	VehicleNo   string    `json:"vehicle_no"`
	Location    string    `json:"location"`
	Parking     time.Time `json:"parking"`
	Leaving     time.Time `json:"leaving"`
	ParkingCost int       `json:"parking_cost"`
	Status      bool      `json:"status"`
}
type LotRevenue struct {
	Location string `json:"location"`
	Amount   int    `json:"amount"`
}

type SummaryResponse struct {
	Revenue   []LotRevenue `json:"revenue"`
	Occupied  int          `json:"occupied"`
	Available int          `json:"available"`
}

type SummaryUsage struct {
	Location string `json:"location"`
	Count    int    `json:"count"`
}

func (d *LotWrapper) SetLot(ctx context.Context, lot Lot) (string, error) {
	// Insert into lot table
	query := `INSERT INTO lot (prime_location_name, price, address, pincode, spots) VALUES ($1, $2, $3, $4, $5)`
	err := d.Db.SetData(ctx, query, lot.Location, lot.Price, lot.Address, lot.Pincode, lot.Spots)
	if err != nil {
		return "", err
	}

	// Fetch the inserted lot ID
	lotId, err := d.GetLotId(ctx, lot)
	if err != nil {
		return "", err
	}
	log.Println("lotId", lotId)

	// Insert 'lot.Spots' number of spots
	insertQuery := `INSERT INTO spot (lot_id, status,address) VALUES ($1, $2,$3)`
	for i := 0; i < lot.Spots; i++ {
		err = d.Db.SetData(ctx, insertQuery, lotId, true, lot.Address)
		if err != nil {
			return "", err
		}
	}

	return lotId, nil
}

func (d *LotWrapper) GetLotId(ctx context.Context, lot Lot) (string, error) {
	if d.Db == nil {
		db.Logger.Error("Database is nil")
		return "", errors.New("internal server error")
	}

	query := `SELECT lot_id FROM lot WHERE prime_location_name = $1 LIMIT 1`
	result, err := d.Db.GetData(ctx, query, lot.Location)
	if err != nil {
		db.Logger.Error("Failed to execute query", zap.Error(err))
		return "", err
	}

	rows, ok := result.(*sql.Rows)
	if !ok {
		db.Logger.Error("Expected *sql.Rows from GetData")
		return "", errors.New("invalid query result")
	}
	defer rows.Close()

	var lotID string
	if rows.Next() {
		err := rows.Scan(&lotID)
		if err != nil {
			db.Logger.Error("Failed to scan lot_id", zap.Error(err))
			return "", err
		}
		return lotID, nil
	}

	// No rows found
	return "", errors.New("lot not found")
}

func (d *LotWrapper) EditLot(ctx context.Context, lot Lot) (Lot, error) {
	query := `UPDATE "lot"
				SET 
				  prime_location_name = $1,
				  price = $2,
				  address = $3,
				  pincode = $4,
				  spots = $5
				WHERE lot_id = $6;`

	err := d.Db.UpdateData(ctx, query,
		lot.Location,
		lot.Price,
		lot.Address,
		lot.Pincode,
		lot.Spots,
		lot.LotID,
	)

	if err != nil {
		return Lot{}, err
	}
	return lot, nil
}

func (d *LotWrapper) GetSpotsByLotID(ctx context.Context, lotID int) ([]Spot, error) {
	if d.Db == nil {
		db.Logger.Error("Database instance is nil")
		return nil, errors.New("internal server error")
	}

	db.Logger.Info("Fetching spots for lot ID", zap.Int("lot_id", lotID))

	query := `SELECT parking_id, lot_id, status, address FROM spot WHERE lot_id = $1`
	result, err := d.Db.GetData(ctx, query, lotID)
	if err != nil {
		db.Logger.Error("Failed to execute query to fetch spot data", zap.Error(err), zap.Int("lot_id", lotID))
		return nil, err
	}

	rows, ok := result.(*sql.Rows)
	if !ok {
		db.Logger.Error("Invalid result format: not *sql.Rows", zap.Any("result", result))
		return nil, errors.New("internal error: invalid query result")
	}
	defer rows.Close()

	var spots []Spot
	count := 0
	for rows.Next() {
		var spot Spot
		err := rows.Scan(
			&spot.ParkingID,
			&spot.LotID,
			&spot.Status,
			&spot.Address,
		)
		if err != nil {
			db.Logger.Error("Failed to scan spot row", zap.Error(err))
			return nil, errors.New("internal error: scan failed")
		}
		spots = append(spots, spot)
		count++
	}

	if err := rows.Err(); err != nil {
		db.Logger.Error("Row iteration error while fetching spots", zap.Error(err))
		return nil, errors.New("internal error: row iteration")
	}

	db.Logger.Info("Successfully fetched spots", zap.Int("lot_id", lotID), zap.Int("spot_count", count))

	return spots, nil
}

func (d *LotWrapper) GetLot(ctx context.Context) ([]Lot, error) {
	if d.Db == nil {
		db.Logger.Error("Database instance is nil")
		return nil, errors.New("internal server error")
	}

	query := `SELECT * FROM "lot"`
	result, err := d.Db.GetData(ctx, query)
	if err != nil {
		db.Logger.Error("Failed to fetch lot data", zap.Error(err))
		return nil, err
	}

	rows, ok := result.(*sql.Rows)
	if !ok {
		db.Logger.Error("Invalid result format: not *sql.Rows")
		return nil, errors.New("internal error: invalid query result")
	}
	defer rows.Close()

	var lots []Lot
	for rows.Next() {
		var lot Lot
		err := rows.Scan(
			&lot.LotID,
			&lot.Location,
			&lot.Price,
			&lot.Address,
			&lot.Pincode,
			&lot.Spots,
		)
		if err != nil {
			db.Logger.Error("Failed to scan lot row", zap.Error(err))
			return nil, errors.New("internal error: scan failed")
		}

		// ✅ Fetch spots for each lot
		spots, err := d.GetSpotsByLotID(ctx, lot.LotID)
		if err != nil {
			db.Logger.Error("Failed to fetch spots for lot", zap.Int("lot_id", lot.LotID), zap.Error(err))
			return nil, errors.New("internal error: could not fetch spots")
		}
		lot.SpotList = spots // ✅ set the spots

		lots = append(lots, lot)
	}

	if err := rows.Err(); err != nil {
		db.Logger.Error("Row iteration error", zap.Error(err))
		return nil, errors.New("internal error: row iteration")
	}

	return lots, nil
}

func (d *LotWrapper) DeleteLotByID(ctx context.Context, lotID int) (string, error) {
	slotQuery := `DELETE FROM spot WHERE lot_id = $1`
	err := d.Db.DeleteData(ctx, slotQuery, lotID)
	if err != nil {
		return "", err
	}
	log.Println("Deleted lot slots")

	query := `DELETE FROM lot WHERE lot_id = $1`
	err = d.Db.DeleteData(ctx, query, lotID)
	if err != nil {
		return "", err
	}
	log.Println("Deleted lot entry")

	return "Deleted", nil
}

func (d *LotWrapper) DeleteSLots(ctx context.Context, lotID int) (string, error) {
	slotQuery := `DELETE FROM spot WHERE parking_id = $1`
	err := d.Db.DeleteData(ctx, slotQuery, lotID)
	if err != nil {
		return "", err
	}
	log.Println("Deleted lot slots")
	return "Deleted", nil
}

func (d *LotWrapper) SearchLots(ctx context.Context, searchType, value string) ([]Lot, error) {
	if d.Db == nil {
		db.Logger.Error("Database is nil")
		return nil, errors.New("internal server error")
	}

	var query string
	var args []interface{}

	switch searchType {
	case "location":
		query = `SELECT lot_id, prime_location_name, price, address, pincode, spots FROM lot WHERE prime_location_name ILIKE '%' || $1 || '%'`
		args = append(args, value)
	case "pincode":
		query = `SELECT lot_id, prime_location_name, price, address, pincode, spots FROM lot WHERE pincode = $1`
		args = append(args, value)
	default:
		return nil, errors.New("invalid search type")
	}

	result, err := d.Db.GetData(ctx, query, args...)
	if err != nil {
		db.Logger.Error("Failed to execute query", zap.Error(err))
		return nil, err
	}

	rows, ok := result.(*sql.Rows)
	if !ok {
		db.Logger.Error("Invalid query result", zap.Any("result", result))
		return nil, errors.New("invalid query result")
	}
	defer rows.Close()

	var lots []Lot
	for rows.Next() {
		var lot Lot
		err := rows.Scan(&lot.LotID, &lot.Location, &lot.Price, &lot.Address, &lot.Pincode, &lot.Spots)
		if err != nil {
			db.Logger.Error("Failed to scan lot row", zap.Error(err))
			return nil, errors.New("scan error")
		}

		// ✅ Fetch SpotList for the lot
		spots, err := d.GetSpotsByLotID(ctx, lot.LotID)
		if err != nil {
			db.Logger.Warn("Failed to fetch spots for lot", zap.Int("lot_id", lot.LotID), zap.Error(err))
			// Optional: continue even if spots fail
			lot.SpotList = []Spot{}
		} else {
			lot.SpotList = spots
		}

		lots = append(lots, lot)
	}

	if err := rows.Err(); err != nil {
		db.Logger.Error("Row iteration error", zap.Error(err))
		return nil, err
	}

	db.Logger.Info("Search complete", zap.Int("result_count", len(lots)))
	return lots, nil
}

func (d *LotWrapper) GetBookings(ctx context.Context, username string) ([]ReservedSpot, error) {
	if d.Db == nil {
		db.Logger.Error("Database instance is nil")
		return nil, errors.New("internal server error")
	}

	query := `SELECT id,location , vehicle_no, spot_id, user_id, parking, leaving, parking_cost,status FROM reserved_spot WHERE user_id = $1`
	result, err := d.Db.GetData(ctx, query, username)
	if err != nil {
		db.Logger.Error("Failed to fetch booking data", zap.Error(err))
		return nil, err
	}

	rows, ok := result.(*sql.Rows)
	if !ok {
		db.Logger.Error("Invalid result format: not *sql.Rows")
		return nil, errors.New("internal error: invalid query result")
	}
	defer rows.Close()

	var bookings []ReservedSpot
	for rows.Next() {
		var booking ReservedSpot
		err := rows.Scan(
			&booking.ID,
			&booking.Location,
			&booking.VehicleNo,
			&booking.SpotID,
			&booking.UserID,
			&booking.Parking,
			&booking.Leaving,
			&booking.ParkingCost,
			&booking.Status,
		)
		if err != nil {
			db.Logger.Error("Failed to scan booking row", zap.Error(err))
			return nil, errors.New("internal error: scan failed")
		}
		bookings = append(bookings, booking)
	}

	if err := rows.Err(); err != nil {
		db.Logger.Error("Row iteration error", zap.Error(err))
		return nil, err
	}

	return bookings, nil
}

func (d *LotWrapper) InsertBooking(ctx context.Context, booking ReservedSpot) error {
	insertQuery := `
		INSERT INTO reserved_spot 
		(id, spot_id, user_id, vehicle_no, location, parking, leaving, parking_cost, status) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	err := d.Db.SetData(ctx, insertQuery,
		booking.ID,
		booking.SpotID,
		booking.UserID,
		booking.VehicleNo,
		booking.Location,
		booking.Parking,
		booking.Leaving,
		booking.ParkingCost,
		booking.Status,
	)
	if err != nil {
		return err
	}

	updateQuery := `
		UPDATE spot
		SET status = false
		WHERE parking_id = $1
	`

	err = d.Db.UpdateData(ctx, updateQuery, booking.SpotID)
	if err != nil {
		return err
	}

	return nil
}

func (d *LotWrapper) ReleaseSpot(ctx context.Context, booking ReservedSpot) (*ReservedSpot, error) {
	if d.Db == nil {
		db.Logger.Error("Database instance is nil")
		return nil, errors.New("internal server error")
	}

	// Step 1: Get price, parking time, and vehicle_no from DB
	query := `
		SELECT l.price, rs.parking, rs.vehicle_no
		FROM reserved_spot rs
		JOIN spot s ON rs.spot_id = s.parking_id
		JOIN lot l ON s.lot_id = l.lot_id
		WHERE rs.id = $1
	`
	result, err := d.Db.GetData(ctx, query, booking.ID)
	if err != nil {
		db.Logger.Error("❌ Failed to fetch price and parking time", zap.Error(err))
		return nil, err
	}
	rows, ok := result.(*sql.Rows)
	if !ok {
		return nil, errors.New("invalid query result format")
	}
	defer rows.Close()

	var pricePerHour int
	var parkingTime time.Time
	var vehicleNo string

	if rows.Next() {
		if err := rows.Scan(&pricePerHour, &parkingTime, &vehicleNo); err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("booking not found")
	}

	// Step 2: Calculate parking duration and cost
	currentTime := time.Now().UTC()
	duration := currentTime.Sub(parkingTime)
	if duration < 0 {
		duration = 0
	}
	hours := int(math.Ceil(duration.Hours()))
	if hours < 1 {
		hours = 1
	}
	parkingCost := hours * pricePerHour

	// Step 3: Update reserved_spot
	updateBookingQuery := `
		UPDATE reserved_spot
		SET leaving = $1,
			parking_cost = $2,
			status = false
		WHERE id = $3
	`
	err = d.Db.UpdateData(ctx, updateBookingQuery, currentTime, parkingCost, booking.ID)
	if err != nil {
		db.Logger.Error("❌ Failed to update reserved_spot", zap.Error(err))
		return nil, err
	}

	// Step 4: Update spot status
	updateSpotQuery := `
		UPDATE spot
		SET status = true
		WHERE parking_id = $1
	`
	err = d.Db.UpdateData(ctx, updateSpotQuery, booking.SpotID)
	if err != nil {
		db.Logger.Error("❌ Failed to update spot status", zap.Error(err))
		return nil, err
	}

	// Step 5: Return updated booking info
	booking.Parking = parkingTime
	booking.Leaving = currentTime
	booking.ParkingCost = parkingCost
	booking.Status = false
	booking.VehicleNo = vehicleNo

	db.Logger.Info("✅ Spot released and booking updated",
		zap.Int("hours", hours),
		zap.Int("cost", parkingCost),
		zap.String("vehicle_no", vehicleNo),
		zap.Time("parking_time", parkingTime),
		zap.Time("leaving_time", currentTime),
	)

	return &booking, nil
}

func (d *LotWrapper) GetSummaryData(ctx context.Context) (*SummaryResponse, error) {
	summary := &SummaryResponse{}

	// 1. Revenue per parking lot
	revenueQuery := `
		SELECT l.prime_location_name, COALESCE(SUM(rs.parking_cost), 0)
		FROM lot l
		LEFT JOIN spot s ON l.lot_id = s.lot_id
		LEFT JOIN reserved_spot rs ON rs.spot_id = s.parking_id
		GROUP BY l.prime_location_name
	`
	res1, err := d.Db.GetData(ctx, revenueQuery)
	if err != nil {
		return nil, err
	}
	defer res1.(*sql.Rows).Close()

	for res1.(*sql.Rows).Next() {
		var rev LotRevenue
		if err := res1.(*sql.Rows).Scan(&rev.Location, &rev.Amount); err != nil {
			return nil, err
		}
		summary.Revenue = append(summary.Revenue, rev)
	}

	// 2. Count occupied and available using one query
	statusCountQuery := `
		SELECT status, COUNT(*) 
		FROM spot 
		GROUP BY status
	`
	res2, err := d.Db.GetData(ctx, statusCountQuery)
	if err != nil {
		return nil, err
	}
	defer res2.(*sql.Rows).Close()

	for res2.(*sql.Rows).Next() {
		var status bool
		var count int
		if err := res2.(*sql.Rows).Scan(&status, &count); err != nil {
			return nil, err
		}
		if status {
			summary.Available = count
		} else {
			summary.Occupied = count
		}
	}

	return summary, nil
}

func (d *LotWrapper) GetUsedParkingSummary(ctx context.Context) ([]SummaryUsage, error) {
	query := `
		SELECT l.prime_location_name, COUNT(*)
		FROM reserved_spot rs
		JOIN spot s ON rs.spot_id = s.parking_id
		JOIN lot l ON s.lot_id = l.lot_id
		WHERE rs.status = false -- already used
		GROUP BY l.prime_location_name
	`

	result, err := d.Db.GetData(ctx, query)
	if err != nil {
		return nil, err
	}

	rows, ok := result.(*sql.Rows)
	if !ok {
		return nil, fmt.Errorf("unexpected result format")
	}
	defer rows.Close()

	var summary []SummaryUsage
	for rows.Next() {
		var su SummaryUsage
		if err := rows.Scan(&su.Location, &su.Count); err != nil {
			return nil, err
		}
		summary = append(summary, su)
	}

	return summary, nil
}
