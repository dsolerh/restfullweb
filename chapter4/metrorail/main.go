package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/dsolerh/restfullweb/chapter4/dbutils"
	"github.com/emicklei/go-restful"
	_ "github.com/go-sql-driver/mysql"
)

// DB Driver visible to whole program
var DB *sql.DB

// StationResource holds information about locations
type StationResource struct {
	ID          int
	Name        string
	OpeningTime time.Time
	ClosingTime time.Time
}

// ScheduleResource links both trains and stations
type ScheduleResource struct {
	ID          int
	TrainID     int
	StationID   int
	ArrivalTime time.Time
}

func main() {
	DB, err := sql.Open("mysql", "test:pass@tcp(localhost:3306)/test?parseTime=true")
	if err != nil {
		log.Fatal("Driver creation failed!")
	}
	if err = DB.Ping(); err != nil {
		log.Fatal(err.Error())
	}
	if err = dbutils.Initialize(DB); err != nil {
		log.Fatal("Fail to initialize Tables")
	}
	defer DB.Close()

	wsContainer := restful.NewContainer()
	wsContainer.Router(restful.CurlyRouter{})

	t := TrainController{
		DB: DB,
	}
	t.Register(wsContainer)

	log.Printf("start listening on localhost:8000")
	server := &http.Server{Addr: ":8000", Handler: wsContainer}
	log.Fatal(server.ListenAndServe())
}

type TrainController struct {
	DB *sql.DB
}

// TrainResource is the model for holding rail information
type TrainResource struct {
	ID              int
	DriverName      string
	OperatingStatus bool
}

// Registeradds paths and routes to a new service instance
func (t *TrainController) Register(container *restful.Container) {
	ws := new(restful.WebService)
	ws.Path("/v1/trains").Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON)
	ws.Route(ws.GET("/{train-id}").To(t.getTrain))
	ws.Route(ws.POST("").To(t.createTrain))
	ws.Route(ws.DELETE("/{train-id}").To(t.removeTrain))
	container.Add(ws)
}

// GET http://localhost:8000/v1/trains/1
func (t TrainController) getTrain(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("train-id")
	train := TrainResource{}
	err := t.DB.QueryRow("SELECT ID, DRIVER_NAME, OPERATING_STATUS FROM train WHERE id=?", id).Scan(&train.ID, &train.DriverName, &train.OperatingStatus)

	if err != nil {
		log.Println(err)
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusNotFound, "Train could	not be found.")
	} else {
		response.WriteEntity(train)
	}
}

// POST http://localhost:8000/v1/trains
func (t TrainController) createTrain(request *restful.Request, response *restful.Response) {
	log.Println("Body: ", request.Request.Body)
	decoder := json.NewDecoder(request.Request.Body)
	var b TrainResource
	err := decoder.Decode(&b)
	if err != nil {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusInternalServerError, err.Error())
	}
	log.Println(b.DriverName, b.OperatingStatus)
	// Error handling is obvious here. So omitting...
	statement, _ := t.DB.Prepare("insert into train (DRIVER_NAME,	OPERATING_STATUS) values (?, ?)")
	result, err := statement.Exec(b.DriverName, b.OperatingStatus)
	if err == nil {
		newID, _ := result.LastInsertId()
		b.ID = int(newID)
		response.WriteHeaderAndEntity(http.StatusCreated, b)
	} else {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusInternalServerError, err.Error())
	}
}

// DELETE http://localhost:8000/v1/trains/1
func (t TrainController) removeTrain(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("train-id")
	statement, _ := t.DB.Prepare("delete from train where id=?")
	_, err := statement.Exec(id)
	if err == nil {
		response.WriteHeader(http.StatusOK)
	} else {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusInternalServerError, err.Error())
	}
}
