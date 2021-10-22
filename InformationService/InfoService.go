package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"os"
)

type Car struct {
	ID    string    `json:"id"`
	Brand string `json:"brand"`
	Model string `json:"model"`
	HP    string `json:"horse_power"`
}
type Configuration struct{
	UserDB string `json:"UserDB"`
	PasswordDB string `json:"PasswordDB"`
	Server string `json:"Server"`
	Port string `json:"Port"`
	Database string `json:"Database"`
}
func read_config() Configuration {
	file,_ := os.Open("../config.json")
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println("error:", err)
		}
	}(file)
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err := decoder.Decode(&configuration)
	if err!= nil{
		fmt.Println("error:", err)
	}
	return configuration
}
func get_record(db_driver *sql.DB, text_id string) Car {
	/* Prepare variable Car to be filled */
	var car Car
	/* Get query and results */
	query := fmt.Sprintf("SELECT * FROM stock WHERE id='%s'", text_id)
	results, err := db_driver.Query(query)
	if err != nil{
		panic(err.Error())
	}
	/* Assign the entry found to the car variable */
	for results.Next(){
		//There should be only one row. Scan the results and bind them to Car
		err = results.Scan(&car.ID, &car.Brand, &car.Model, &car.HP)
		if err != nil{
			panic(err.Error())
		}

	}
	/* If car brand is empty, it did not retrieve correctly. Throw an exception to parent function */
	if car.Brand == ""{
		panic("The Id was not in the DB. Perhaps the type comparison fails?")
	}
	return car
}
func connect_to_database() *sql.DB{
	config := read_config()
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", config.UserDB , config.PasswordDB,
		config.Server, config.Port, config.Database)
	print(dataSourceName)
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil{
		panic(err.Error())
	}
	return db
}



func processRequest(w http.ResponseWriter, r *http.Request){
	// Create SQL driver and defer its closure until the end of main
	dbDriver := connect_to_database()
	r.Close = true
	defer func(db_driver *sql.DB) {
		err := db_driver.Close()
		if err != nil {
			log.Fatal("DB Cannot be closed")
		}
	}(dbDriver)

	/* In case there are no cars, an exception is risen and processed. Return HTTP Response accordingly.
	I know it could be done otherwise but I wanted to try the recover() func*/
	defer func() {no_cars := recover()
		if no_cars != nil{
			w.WriteHeader(http.StatusBadRequest)
			log.Print("[REQUESTFAILED] Bad Request received.")
			return

		}}()

	/* Begin processing the request. Get the id value */
	idRequested := r.FormValue("id")
	/*Retrieve the record with that ID. Notice that this will throw an exception when no entries are found.
	This exception is treated by the recover() */
	carRetrieved := get_record(dbDriver, idRequested)
	/*If not interrupted by the panic and recover process, there is an entry. Prepare an ANSWER
	(200 Ok is set by default) with the data. JSON is fine. Structured content (as Car) is 'jsonable' */
	w.Header().Set("Content-Type", "application/json")
	carJson,_ := json.Marshal(carRetrieved)
	/* Send and get possible errors */
	_, err := w.Write(carJson)

	if err != nil {
		log.Fatal("Failed Reply")

	}

	/* Logging purposes */
	log.Print(fmt.Sprintf("[REQUESTED] Car ID: %s Car Brand: %s Car Model: %s, Car HP: %s",
		carRetrieved.ID, carRetrieved.Brand, carRetrieved.Model, carRetrieved.HP))

}

func handleRequests(){
	http.HandleFunc("/", processRequest)
	log.Fatal(http.ListenAndServe(":8081", nil))
}
func main() {handleRequests()}