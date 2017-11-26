package main

import (
	"fmt"
	"net/http"
	"net/url"
	"log"
	"io/ioutil"
	"database/sql"
  _ "github.com/mattn/go-sqlite3"
	"html/template"
	"strings"
)


// struct to hold data retrieved from SQL queries...
type Doctor struct {
	Id template.HTML
	Firstname string
	Lastname string
	Email template.HTML
	Gender	string
	Address string
	City template.HTML
	Phone string
	Image template.HTML
	Openings template.HTML
	Specialty string
}

// global variables that need to be accessible in all functions...
var db *sql.DB
var err error
var searchTemplate *template.Template
var dentistTemplate *template.Template

// check errors returned at various points in the code... panic if something happens...
func check(e error) {
    if e != nil {
        panic(e)
    }
}


// create the SQL query based on HTML form parameters...
func buildQuery(values url.Values) string {
		queryString := " WHERE "
		for k := range values {
			if values[k][0] != "" {
					if (k == "gender" && values[k][0] == "both") || (k == "specialty" && values[k][0] == "all") {
						continue
					}
					queryString = queryString + k + "='" + values[k][0] + "' AND "
			}

		}
		if values["gender"][0] == "both" && values["first_name"][0] == ""  && values["last_name"][0] ==  "" && values["email"][0] == ""  && values["address"][0] == ""  && values["city"][0] == "" && values["phone"][0] == ""  && values["specialty"][0] == "" {
				return ""
		}
		if values["gender"][0] == "" && values["first_name"][0] == ""  && values["last_name"][0] ==  "" && values["email"][0] == ""  && values["address"][0] == ""  && values["city"][0] == "" && values["phone"][0] == ""  && values["specialty"][0] == "all" {
				return ""
		}
		return queryString[:len(queryString)-5]
}

// Insert data into the "/search" result page dynamically
func insertImageAndCity(docs []Doctor) []Doctor {
	for i, element := range docs {
		if docs[i].Id == "N/A" {
			docs[i].Image = "<img alt=\"Missing\" src=\"" + element.Image + "\" />"
		} else {
			docs[i].Image = "<img alt=\"Missing\" src=\"" + element.Image + "\" />"
			docs[i].City = "<a href=\"https://www.google.com/maps?q=" + element.City + "\" target=\"_blank\">" + element.City + "</a>"
			docs[i].Email = "<a href=\"mailto:"  + element.Email + "?Subject=Dentist%20appointment\" target=\"_top\">" + element.Email + "</a>"
			docs[i].Openings = "<a href=\"/dentist?id=" + element.Id + "\">Click to view</a>"
		}
  }
	return docs
}

// Insert data into the "/dentist?id=X" page dynamically
func insertDentistImageAndCity(w http.ResponseWriter, doc Doctor) {
	temp1 := doc.Image
	doc.Image = "<img alt=\"Missing\" src=\"" + temp1 + "\" />"
	temp2 := doc.City
	doc.City = "<a href=\"https://www.google.com/maps?q=" + temp2 + "\" target=\"_blank\">" + temp2 + "</a>"
	temp3 := doc.Email
	doc.Email = "<a href=\"mailto:"  + temp3 + "?Subject=Dentist%20appointment\" target=\"_top\">" + temp3 + "</a>"
	temp4 := doc.Openings
	doc.Openings = template.HTML(strings.Replace(string(temp4), ",", "<br>", -1))
	dentistTemplate.Execute(w, doc)
}

// Handles the search results page if there were no results found according to
// the request specified in the query > EMPTY page <
func parseEmptySearch(w http.ResponseWriter) {
	searchTemplate.Execute(w, insertImageAndCity([]Doctor{
						Doctor{Id: "Not found",
						Firstname: "Not found",
						Lastname: "Not found",
						Email: "Not found",
						Gender: "Not found",
						Address: "Not found",
						City: "Not found",
						Phone: "Not found",
						Image: `data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAABAAAAAQCAYAA
						AAf8/9hAAAABGdBTUEAAK/INwWK6QAAABl0RVh0U29mdHdhcmUAQWRvYmUgSW1hZ2VS
						ZWFkeXHJZTwAAAHdSURBVDjLpZNraxpBFIb3a0ggISmmNISWXmOboKihxpgUNGWNSpv
						aS6RpKL3Ry//Mh1wgf6PElaCyzq67O09nVjdVlJbSDy8Lw77PmfecMwZg/I/GDw3DCo
						8HCkZl/RlgGA0e3Yfv7+DbAfLrW+SXOvLTG+SHV/gPbuMZRnsyIDL/OASziMxkkKkUQ
						TJJsLaGn8/iHz6nd+8mQv87Ahg2H9Th/BxZqxEkEgSrq/iVCvLsDK9awtvfxb2zjD2A
						RID+lVVlbabTgWYTv1rFL5fBUtHbbeTJCb3EQ3ovCnRC6xAgzJtOE+ztheYIEkqbFaS
						3vY2zuIj77AmtYYDusPy8/zuvunJkDKXM7tYWTiyGWFjAqeQnAD6+7ueNx/FLpRGAru
						7mcoj5ebqzszil7DggeF/DX1nBN82rzPqrzbRayIsLhJqMPT2N83Sdy2GApwFqRN7jF
						PL0tF+10cDd3MTZ2AjNUkGCoyO6y9cRxfQowFUbpufr1ct4ZoHg+Dg067zduTmEbq4y
						i/UkYidDe+kaTcP4ObJIajksPd/eyx3c+N2rvPbMDPbUFPZSLKzcGjKPrbJaDsu+dQO
						3msfZzeGY2TCvKGYQhdSYeeJjUt21dIcjXQ7U7Kv599f4j/oF55W4g/2e3b8AAAAASU
						VORK5CYII=`,
						Openings: "Not found",
						Specialty: "Not found"}}))
}

// handles requests to URL "/" (main page)
func mainPageHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Printf("Request arrived to URL: %s\n\n", r.URL)
	dat, err := ioutil.ReadFile("./html/index.html")
	check(err)
	fmt.Fprintf(w,string(dat))
}

// Handles request made to "/search?param1=X&param2=Y&..."
func searchHandler(w http.ResponseWriter, r *http.Request) {
  r.ParseForm()
	fmt.Printf("Request arrived to URL: %s\n\n", r.URL)
	var DoctorsFound []Doctor
	var newDoc Doctor
	baseQuery := `SELECT * FROM enigma`
	rows, err := db.Query(baseQuery + buildQuery(r.Form))
	check(err)
 	for rows.Next() {
        err = rows.Scan(&newDoc.Id, &newDoc.Firstname, &newDoc.Lastname, &newDoc.Email,
							&newDoc.Gender, &newDoc.Address, &newDoc.City, &newDoc.Phone,
							&newDoc.Image, &newDoc.Openings, &newDoc.Specialty)
				check(err)
				DoctorsFound = append(DoctorsFound, newDoc)
  }
	if len(DoctorsFound) == 0 {
		parseEmptySearch(w)
		return
	}
	searchTemplate.Execute(w, insertImageAndCity(DoctorsFound))
}


// Handles requests made to URL: "/dentist?id=X"
func dentistHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Printf("Request arrived to URL: %s\n\n", r.URL)
	var DoctorFound Doctor
	rows, err := db.Query("SELECT * FROM enigma WHERE id ='" + r.Form["id"][0] + "'")
	check(err)
	rows.Next()
  err = rows.Scan(&DoctorFound.Id, &DoctorFound.Firstname, &DoctorFound.Lastname,
					&DoctorFound.Email,	&DoctorFound.Gender, &DoctorFound.Address,
					&DoctorFound.City, &DoctorFound.Phone, &DoctorFound.Image,
					&DoctorFound.Openings, &DoctorFound.Specialty)
	check(err)
	insertDentistImageAndCity(w, DoctorFound)
}


func main() {
		// OPEN Database
		db, err = sql.Open("sqlite3", "./database/enigma.db")
		check(err)
		// Parse HTML templates that will be populated dynamically
		searchTemplate = template.Must(template.ParseFiles("./html/search.html"))
		dentistTemplate = template.Must(template.ParseFiles("./html/dentist.html"))

		// Define URL handler functions
		http.HandleFunc("/", mainPageHandler)
    http.HandleFunc("/search", searchHandler)
		http.HandleFunc("/dentist", dentistHandler)
		fs := http.FileServer(http.Dir("./html/static"))
		http.Handle("/static/", http.StripPrefix("/static", fs))

		// Start  server
		fmt.Println("Server is running...")
		err = http.ListenAndServe(":8080", nil)
		if err != nil {
			log.Fatal("ListenAndServe: ", err)
		}
		db.Close()
}
