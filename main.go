package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/google/uuid"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Define a model struct
type User struct {
	ID        string    `json:"id"`
	Name string `json:"name" binding:"required,min=5,max=10"`
	Age  int    `json:"age" binding:"required,min=0,max=25"`
	CreatedAt time.Time `json:"createdAt"`

}



	
	var db *gorm.DB

func main() {
	// Connect to MySQL database
	dsn := "root:Sainath@2000@tcp(127.0.0.1:3306)/testdata"

	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	
	// AutoMigrate will create the table. You can also manually create the table
	db.AutoMigrate(&User{})

	
	// Initialize Gin
	router := gin.Default()

	// Define routes

	router.POST("/user", createUser)

	// Run the server
	router.Run(":9087")

}

// Handler functions
func createUser(c *gin.Context) {
	var user []User
    err:=	c.ShouldBindJSON(&user)
    if err!=nil{
		c.JSON(http.StatusBadRequest,gin.H{"error":err.Error})
	}


	WriteData(user)
    users:= ReadData()

	input:=make(chan User)
	go InsertData2(input)
	go InsertData1(input)
	for _,val:=range users {
	     input<-val
	}

	c.JSON(http.StatusCreated, user)
}


func InsertData1(user chan User){


	for {
		data:=<-user
	
		// Generate UUID for the user ID
        data.ID = uuid.New().String()

		// Set the timestamp
	    data.CreatedAt = time.Now()
	    fmt.Println("data passed from first router",data)
	    if err := db.Table("users").Create(&data).Error; err != nil {
			fmt.Println(err.Error())
	    }
    }

}

func InsertData2(user chan User){
	
	
	for {
		data:=<-user
		// Generate UUID for the user ID
        data.ID = uuid.New().String()

     // Set the timestamp
       data.CreatedAt = time.Now()
		fmt.Println("data passed from second router",data)
	    if err := db.Table("users").Create(&data).Error; err != nil {
			fmt.Println(err.Error())
	    }
    }

}


func WriteData(user []User) {
	// Open the file for writing
	file, err := os.Create("people.json")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	// Create a JSON encoder
	encoder := json.NewEncoder(file)

	// Encode and write the JSON data
	err = encoder.Encode(user)
	if err!=nil{
		fmt.Println("Error encoding JSON:", err)
		return
	}

	fmt.Println("JSON data written to file successfully")
}




func ReadData()[]User {
	// Open the JSON file for reading
	file, err := os.Open("people.json")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil
	}
	defer file.Close()

	// Create a slice to hold decoded data
	var people []User

	// Create a JSON decoder
	decoder := json.NewDecoder(file)

	// Decode JSON data from the file
	if err := decoder.Decode(&people); err != nil {
		fmt.Println("Error decoding JSON:", err)
		return nil
	}

	return people
}


