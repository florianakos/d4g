package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
  "encoding/csv"
  "io"
  "reflect"
  "strconv"
  "strings"
  "bufio"
)

type Doctor struct {
  Id string
	Firstname string
	Lastname string
	Email string
	Gender	string
	Address string
	City string
	Phone string
	Image string
	Openings string
	Specialty string
}

func Unmarshal(reader *csv.Reader, v interface{}) error {
	record, err := reader.Read()
	if err != nil {
		return err
	}
	s := reflect.ValueOf(v).Elem()
	if s.NumField() != len(record) {
		return &FieldMismatch{s.NumField(), len(record)}
	}
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		switch f.Type().String() {
		case "string":
			f.SetString(record[i])
		case "int":
			ival, err := strconv.ParseInt(record[i], 10, 0)
			if err != nil {
				return err
			}
			f.SetInt(ival)
		default:
			return &UnsupportedType{f.Type().String()}
		}
	}
	return nil
}

type FieldMismatch struct {
	expected, found int
}

func (e *FieldMismatch) Error() string {
	return "CSV line fields mismatch. Expected " + strconv.Itoa(e.expected) + " found " + strconv.Itoa(e.found)
}

type UnsupportedType struct {
	Type string
}

func (e *UnsupportedType) Error() string {
	return "Unsupported type: " + e.Type
}


func formatOpeningHours(openingHours string) string {
  openingHours = strings.Replace(openingHours, "[{\"mon\"", "Monday", -1)
  openingHours = strings.Replace(openingHours, "\"},\"tue\"", ", Tuesday", -1)
  openingHours = strings.Replace(openingHours, "\",\"close\":\"", "-", -1)
  openingHours = strings.Replace(openingHours, "{\"open\":\"", " ", -1)
  openingHours = strings.Replace(openingHours, "\"},\"wed\"", ", Wednesday", -1)
  openingHours = strings.Replace(openingHours, "\"},\"thu\"", ", Thursday", -1)
  openingHours = strings.Replace(openingHours, "\"},\"fri\"", ", Friday", -1)
  openingHours = strings.Replace(openingHours, "\"}}]", "", -1)
  return openingHours
}


func insertDoc(doc Doctor, db *sql.DB) bool {
  temp := doc.Openings
  doc.Openings = formatOpeningHours(temp)

  //INSERT CSV DATA INTO DB
  insertEnigma := fmt.Sprintf(`INSERT INTO enigma
												          (id,first_name,last_name,email,gender,address,city,phone,image,openings,specialty)
									             VALUES
												          ("%s", "%s", "%s", "%s", "%s", "%s", "%s", "%s", "%s", "%s", "%s");`,
                                doc.Id, doc.Firstname, doc.Lastname, doc.Email, doc.Gender, doc.Address, doc.City,
                                doc.Phone, doc.Image, doc.Openings, doc.Specialty)
	_, err := db.Exec(insertEnigma)
	if err != nil {
			log.Println(err)
	}
	fmt.Println("row added")
  return true
}


func main() {
  // OPEN DATABASE
  db, err := sql.Open("sqlite3", "./enigma.db")
  if err != nil {
		log.Fatal(err)
  }
  defer db.Close()

  file, err := os.Open("./data.csv")
  if err != nil {
      log.Fatal(err)
  }
  defer file.Close()

  scanner := bufio.NewScanner(file)
  for scanner.Scan() {
      var reader = csv.NewReader(strings.NewReader(scanner.Text()))
      reader.Comma = ','
      var doc Doctor
      for {
    		err := Unmarshal(reader, &doc)
    		if err == io.EOF {
    			break
    		}
    		if err != nil {
    			panic(err)
    		}
        insertDoc(doc, db)
    	}


  }
  if err := scanner.Err(); err != nil {
      log.Fatal(err)
  }
}
