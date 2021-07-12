package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/lib/pq"
	"github.com/subosito/gotenv"
)

var db *sql.DB
var rLT, rLD, pD float64 //rider lat, long and pickup Dist
var pT, trT time.Time    // pickup, travel time
var der int              // dr_id
var rdrid int = 1111     // to maintain ride_id
var uL string            // user location name

type Ride struct {
	R_id           int     `json:"r_id"`
	Dr_id          int     `json:"dr_id"`
	Starting_point string  `json:"starting_point"`
	Ending_point   string  `json:"ending_point"`
	Distance       float64 `json:"distance"`
	Travel_time    string  `json:"travel_time"`
	Pickup_time    string  `json:"pickup_time"`
	CreatedAt      string  `json:"createdAt"`
	UpdatedAt      string  `json:"updatedAt"`
}

var rides []Ride

type Cab struct {
	Cab_id      int     `json:"cab_id"`
	Driver_name string  `json:"driver_name"`
	Model       string  `json:"model"`
	Reg_no      string  `json:"reg_no"`
	Status      string  `json:"status"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	CreatedAt   string  `json:"createdAt"`
	UpdatedAt   string  `json:"updatedAt"`
}

var cabs []Cab

type Payment struct {
	Ride_id   int     `json:"ride_id"`
	P_type    string  `json:"p_type"`
	Time      string  `json:"time"`
	Cost      float64 `json:"cost"`
	Status    string  `json:"status"`
	CreatedAt string  `json:"createdAt"`
	UpdatedAt string  `json:"updatedAt"`
}

var payments Payment

type Dr struct {
	Id_dr    int `json:"id_dr"`
	Cab_id   int `json:"cab_id"`
	Rider_id int `json:"rider_id"`
}

var drs Dr

type Rider struct {
	Cust_id   int     `json:"cust_id"`
	Password  string  `json:"password"`
	Name      string  `json:"name"`
	Phone     int     `json:"phone"`
	CreatedAt string  `json:"createdAt"`
	UpdatedAt string  `json:"updatedAt"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

var rdr Rider

type Location struct {
	Loc_id    int     `json:"loc_id"`
	Name      string  `json:"name"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	CreatedAt string  `json:"createdAt"`
	UpdatedAt string  `json:"updatedAt"`
}

var loci Location

func init() { // initializing values in .env file
	gotenv.Load()
}

func logFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func distn(lat1 float64, lng1 float64, lat2 float64, lng2 float64, unit ...string) float64 {
	const PI float64 = 3.141592653589793

	radlat1 := float64(PI * lat1 / 180)
	radlat2 := float64(PI * lat2 / 180)

	theta := float64(lng1 - lng2)
	radtheta := float64(PI * theta / 180)

	dist := math.Sin(radlat1)*math.Sin(radlat2) + math.Cos(radlat1)*math.Cos(radlat2)*math.Cos(radtheta)

	if dist > 1 {
		dist = 1
	}

	dist = math.Acos(dist)
	dist = dist * 180 / PI
	dist = dist * 60 * 1.1515

	if len(unit) > 0 {
		if unit[0] == "K" {
			dist = dist * 1.609344
		} else if unit[0] == "N" {
			dist = dist * 0.8684
		}
	}

	return dist
}

func checkCount(rows *sql.Rows) (count int) {
	for rows.Next() {
		err := rows.Scan(&count)
		logFatal(err)
	}
	return count
}

func getRides(w http.ResponseWriter, r *http.Request) {
	var ride Ride    //variable type is struct Ride
	rides = []Ride{} //assigning empty values mapped to attributes of struct Ride

	rows, err := db.Query("select * from ride") // passing query
	logFatal(err)                               //checking errors

	defer rows.Close()

	for rows.Next() { // assigning values to slice rides and checking errors
		err := rows.Scan(&ride.R_id, &ride.Dr_id, &ride.Starting_point, &ride.Ending_point, &ride.Distance, &ride.Travel_time, &ride.Pickup_time, &ride.CreatedAt, &ride.UpdatedAt)
		logFatal(err)

		rides = append(rides, ride)
	}

	json.NewEncoder(w).Encode(&rides) //encoding values of slice rides to json

}

func getCab(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)

	var riders Rider
	rdr := []Rider{}

	var Cabs Cab   // new Cabs variable of type Cab
	cabs = []Cab{} // already declared cabs variable initialized to empty and mapped to struct Cab

	rid, _ := strconv.Atoi(params["cust_id"])

	rows, err := db.Query("select * from Rider where cust_id=$1", rid)
	logFatal(err)

	defer rows.Close()

	for rows.Next() {

		err := rows.Scan(&riders.Cust_id, &riders.Name, &riders.Password, &riders.Phone, &riders.CreatedAt, &riders.UpdatedAt, &riders.Latitude, &riders.Longitude)
		logFatal(err)

		rdr = append(rdr, riders)
	}

	json.NewEncoder(w).Encode(&riders)

	var lt, ln float64 = 0, 0

	for _, rd := range rdr {
		lt = rd.Latitude
		ln = rd.Longitude
	}

	var min float64 = 0
	var cbid int = 0

	rLT = lt
	rLD = ln

	rows, err = db.Query("select * from Cab")
	logFatal(err)

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&Cabs.Cab_id, &Cabs.Driver_name, &Cabs.Model, &Cabs.Reg_no, &Cabs.Status, &Cabs.Latitude, &Cabs.Longitude, &Cabs.CreatedAt, &Cabs.UpdatedAt)
		logFatal(err)

		cabs = append(cabs, Cabs)

	}

	for _, cab := range cabs {

		if min < distn(cab.Latitude, cab.Longitude, lt, ln) {
			min = distn(cab.Latitude, cab.Longitude, lt, ln)
			cbid = cab.Cab_id
		}
	}

	pD = min

	for _, cab := range cabs {
		if cab.Cab_id == cbid {
			json.NewEncoder(w).Encode(&cab)
		}
	}

	//updating DR db with cab_id and rider_id

	// var ddr Dr

	nr, _ := db.Query("select count(*) as count from Dr")
	nn := (checkCount(nr))

	der = nn + 1

	rows, err = db.Query("insert into Dr(id_dr,cab_id,rider_id) values($1,$2,$3)", der, cbid, rid)
	logFatal(err)

}

func bookCaB(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)

	var loc Location
	locis := []Location{}

	rows, err := db.Query("select * from Location where name=$1", params["name"])
	logFatal(err)

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&loc.Loc_id, &loc.Name, &loc.Latitude, &loc.Longitude, &loc.CreatedAt, &loc.UpdatedAt)
		logFatal(err)

		locis = append(locis, loc)
	}

	var lt, ln float64 = 0, 0
	var nm string

	for _, rd := range locis {
		lt = rd.Latitude
		ln = rd.Longitude
		nm = rd.Name
	}

	dt := distn(rLT, rLD, lt, ln)
	t := time.Now()
	var as time.Duration = (time.Duration)(math.Round(pD / 35 * 60))
	// as := 15
	var pt time.Time = t.Add(time.Minute * as)

	pT = pt

	var am time.Duration = (time.Duration)(math.Round(dt / 32 * 60))

	var tt time.Time = (pt.Add(time.Minute * am))
	trT = tt

	cost := dt * 12

	rows, err = db.Query("select * from Location")
	logFatal(err)

	var umin float64 = 0

	for rows.Next() {
		err = rows.Scan(&loc.Loc_id, &loc.Name, &loc.Latitude, &loc.Longitude, &loc.CreatedAt, &loc.UpdatedAt)
		logFatal(err)

		if umin < distn(rLT, rLD, loc.Latitude, loc.Latitude) {
			umin = distn(rLT, rLD, loc.Latitude, loc.Latitude)
			uL = loc.Name
		}

	}

	kk := rdrid + 1

	var rd Ride
	bRide := []Ride{}

	rows, err = db.Query("insert into Ride values($1,$2,$3,$4,$5,$6,$7,&8,&9)", kk, der, uL, nm, dt, trT, pT, time.Now(), "0000-00-00 00:00:00")
	logFatal(err)

	for rows.Next() {
		err = rows.Scan(&rd.R_id, &rd.Dr_id, &rd.Starting_point, &rd.Ending_point, &rd.Distance, &rd.Travel_time, &rd.Pickup_time, &rd.CreatedAt, &rd.UpdatedAt)
		logFatal(err)

		bRide = append(bRide, rd)
	}

	json.NewDecoder(r.Body).Decode(&bRide)

	rows, err = db.Query("insert into Payment values($1,$2,$3,$4,$5,$6,$7)", kk, "UPI", trT, cost, "paid", time.Now(), "0000-00-00 00:00:00")
	logFatal(err)

	var py Payment
	pay := []Payment{}

	for rows.Next() {
		err = rows.Scan(&py.Ride_id, &py.P_type, &py.Time, &py.Cost, &py.Status, &py.CreatedAt, &py.UpdatedAt)
		logFatal(err)

		pay = append(pay, py)
	}

	json.NewEncoder(w).Encode(pay)

}

func getHistory(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	var ddr Dr
	drr := []Dr{}
	var Costs Payment
	payments := []Payment{}

	rows, err := db.Query("select * from Dr where rider_id=$1", params["rider_id"])
	logFatal(err)

	defer rows.Close()

	for rows.Next() { // getting dr_id from DR by giving rider_id
		err := rows.Scan(&ddr.Id_dr, &ddr.Cab_id, &ddr.Rider_id)
		logFatal(err)

		drr = append(drr, ddr)
	}

	k := 0

	for _, o := range drr {
		k = o.Id_dr
	}

	var ride Ride    //variable type is struct Ride
	rides = []Ride{} //assigning empty values mapped to attributes of struct Ride

	rows, err = db.Query("select * from Ride where dr_id=$1", k)
	logFatal(err)

	defer rows.Close()

	for rows.Next() { // assigning values to slice rides and checking errors
		err := rows.Scan(&ride.R_id, &ride.Dr_id, &ride.Starting_point, &ride.Ending_point, &ride.Distance, &ride.Travel_time, &ride.Pickup_time, &ride.CreatedAt, &ride.UpdatedAt)
		logFatal(err)

		rides = append(rides, ride)
	}

	kr := 0

	for _, o := range rides {
		kr = o.R_id
	}

	rows, err = db.Query("select * from Payment where ride_id=$1", kr)
	logFatal(err)

	for rows.Next() {
		err := rows.Scan(&Costs.Ride_id, &Costs.P_type, &Costs.Time, &Costs.Cost, &Costs.Status, &Costs.CreatedAt, &Costs.UpdatedAt)
		logFatal(err)

		payments = append(payments, Costs)
	}

	defer rows.Close()

	json.NewEncoder(w).Encode(&rides)
	json.NewEncoder(w).Encode(&payments)

}

func main() {

	dbUrl, err := pq.ParseURL(("postgres://bayoeycm:6aoB3F9L80EYnxikSjhs7OnjzgzNSbNm@batyr.db.elephantsql.com/bayoeycm"))

	//dbUrl, err := pq.ParseURL(os.Getenv("DB_URL")) // connecting to DB
	logFatal(err) //checking for error

	//log.Println(dbUrl)

	db, err = sql.Open("postgres", dbUrl)

	err = db.Ping()
	logFatal(err)

	router := mux.NewRouter() //for routing

	router.HandleFunc("/rides", getRides).Methods("GET") // (endpoint, handlerFunction)
	router.HandleFunc("/getCab/{rider_id}", getCab).Methods("GET")
	router.HandleFunc("/ridesHistory/{rider_id}", getHistory).Methods("GET")
	router.HandleFunc("/bookRide/{ending_point}", bookCaB).Methods("POST")
	// router.HandleFunc("/rides", addCab).Methods("POST")
	// router.HandleFunc("/rides", addLocation).Methods("POST")
	// router.HandleFunc("/rides", updateUser).Methods("PUT")
	// router.HandleFunc("/rides/{cust_id}", removeUser).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8588", router)) // start server

}
