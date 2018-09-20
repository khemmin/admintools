package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mssql"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

var db *gorm.DB
var err error

type User struct {
	gorm.Model
	Name  string `json:"name"`
	Email string
}

type Config struct {
	gorm.Model
	ConfName      string `json:"name"`
	ConfHost      string `json:"host"`
	ConfPort      string `json:"port"`
	ConfDatabase  string `json:"database"`
	ConfUsername  string `json:"username"`
	ConfPassword  string `json:"password"`
	ConfNameStore string `json:"namestore"`
}

type Cat struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type Dog struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type Hamster struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

func InitialMigration() {
	db, err := gorm.Open("mssql", "sqlserver://sa:Admin1234@192.9.58.78:1433?database=admin_tools")

	if err != nil {
		fmt.Println(err.Error())
		panic("Failed to connect to database")
	}
	defer db.Close()

	db.AutoMigrate(&Config{})
}

func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello from the website!")
}

func getCats(c echo.Context) error {
	catName := c.QueryParam("name")
	catType := c.QueryParam("type")

	dataType := c.Param("data")

	if dataType == "string" {
		return c.String(http.StatusOK, fmt.Sprintf("your cat name is : %s\nand his type is : %s\n", catName, catType))
	}

	if dataType == "json" {
		return c.JSON(http.StatusOK, map[string]string{
			"name": catName,
			"type": catType,
		})
	}

	return c.JSON(http.StatusBadRequest, map[string]string{
		"error": "you need to lets us know if you want json or string data",
	})

}

func addCats(c echo.Context) error {
	cat := Cat{}

	defer c.Request().Body.Close()

	b, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		log.Printf("Failed reading the request body for addCats: %s", err)
		return c.String(http.StatusInternalServerError, "")
	}

	err = json.Unmarshal(b, &cat)
	if err != nil {
		log.Printf("Failed unmarshaling in addCats: %s", err)
		return c.String(http.StatusInternalServerError, "")
	}

	log.Printf("this is your cat: %#v", cat)
	return c.String(http.StatusOK, "We got your cat!")
}

func addDogs(c echo.Context) error {
	dog := Dog{}

	defer c.Request().Body.Close()

	err := json.NewDecoder(c.Request().Body).Decode(&dog)
	if err != nil {
		log.Printf("Failed processing addDog request: %s", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	log.Printf("this is your dog: %#v", dog)
	return c.String(http.StatusOK, "We got your dog!")
}

func addHamster(c echo.Context) error {
	hamster := Hamster{}

	err := c.Bind(&hamster)

	if err != nil {
		log.Printf("Failed processing addHamster request: %s", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	log.Printf("this is your hamster: %#v", hamster)
	return c.String(http.StatusOK, "We got your hamster!")
}

func AllUsers(c echo.Context) error {
	users := User{}

	err := c.Bind(&users)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid connection string")
	}
	db, err = gorm.Open("mssql", "sqlserver://sa:Admin1234@192.9.58.78:1433?database=admin_tools")
	defer db.Close()

	if err != nil {
		return c.String(http.StatusBadRequest, "Cannot connect to database")
	}
	return c.String(http.StatusOK, "Success")

}

func addUsers(c echo.Context) error {
	users := User{}
	err := c.Bind(&users)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid connection string")
	}
	db, err = gorm.Open("mssql", "sqlserver://sa:Admin1234@192.9.58.78:1433?database=admin_tools")
	defer db.Close()
	log.Println(users.Name)
	log.Println(users.Email)
	db.Create(&users)
	db.Save(&users)

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(db)
		return c.JSON(http.StatusCreated, &users)
	}
	return c.String(http.StatusOK, "add Success")
}

func AllConfig(c echo.Context) error {
	config := Config{}

	err := c.Bind(&config)

	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid connection string")
	}
	db, err = gorm.Open("mssql", "sqlserver://sa:Admin1234@192.9.58.78:1433?database=admin_tools")
	defer db.Close()

	if err != nil {
		return c.String(http.StatusBadRequest, "Cannot connect to database")
	}
	return c.String(http.StatusOK, "Success")
}

func addConfig(c echo.Context) error {
	config := Config{}
	err := c.Bind(&config)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid connection string")
	}
	db, err = gorm.Open("mssql", "sqlserver://sa:Admin1234@192.9.58.78:1433?database=admin_tools")
	defer db.Close()
	db.Create(&config)
	db.Save(&config)

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(db)
		return c.JSON(http.StatusCreated, &config)
	}
	return c.String(http.StatusOK, "add Success")
}

func selectUserOnProc(c echo.Context) error {
	db, err := gorm.Open("mssql", "sqlserver://sa:Admin1234@192.9.58.78:1433?database=admin_tools")
	defer db.Close()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Cannot connect to database")
	}
	type sp struct {
		Name  string
		Email string
	}
	var sps []sp
	db.Raw("EXEC SelectAllUsers").Scan(&sps)
	return c.JSON(http.StatusOK, sps)
}

func selectViewUser(c echo.Context) error {
	db, err := gorm.Open("mssql", "sqlserver://sa:Admin1234@192.9.58.78:1433?database=admin_tools")
	defer db.Close()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Cannot connect to database")
	}
	type sp struct {
		Id    string
		Name  string
		Email string
	}
	var sps []sp
	db.Raw("SELECT * FROM ViewAllUser").Scan(&sps)
	return c.JSON(http.StatusOK, sps)
}

func getDataTypeUser(c echo.Context) error {
	db, err := gorm.Open("mssql", "sqlserver://sa:Admin1234@192.9.58.78:1433?database=admin_tools")
	defer db.Close()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Cannot connect to database")
	}
	type sp struct {
		Field_Name  string
		Data_Type   string
		Length_Size string
	}
	var sps []sp
	db.Raw("SELECT OBJECT_SCHEMA_NAME (c.object_id) SchemaName, o.Name AS Table_Name, c.Name AS Field_Name, t.Name AS Data_Type, t.max_length AS Length_Size, t.precision AS Precision FROM sys.columns c INNER JOIN sys.objects o ON o.object_id = c.object_id LEFT JOIN  sys.types t on t.user_type_id  = c.user_type_id WHERE --o.type = 'U' o.Name = 'users' ORDER BY o.Name, c.Name").Scan(&sps)
	return c.JSON(http.StatusOK, sps)
}

func getTypeStoreProc(c echo.Context) error {
	db, err := gorm.Open("mssql", "sqlserver://sa:Admin1234@192.9.58.78:1433?database=admin_tools")
	defer db.Close()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Cannot connect to database")
	}
	type sp struct {
		Parameter_name string
		Type           string
		Length         string
		Prec           string
		Scale          string
		Param_order    string
		Collation      string
	}
	var sps []sp
	db.Raw("EXEC AllUsers").Scan(&sps)

	return c.JSON(http.StatusOK, sps)
}

func main() {
	fmt.Println("Welcome to the server")
	InitialMigration()
	e := echo.New()
	e.Use(middleware.CORS())
	e.Use(middleware.Logger())
	e.GET("/", hello)
	e.GET("/cats/:data", getCats)
	e.GET("/users", AllUsers)
	e.GET("/configs", AllConfig)
	e.GET("/userselect", selectUserOnProc)
	e.GET("/viewuser", selectViewUser)
	e.GET("/getdatatype", getDataTypeUser)
	e.GET("/gettypeuser", getTypeStoreProc)

	e.POST("/configs", addConfig)
	e.POST("/users", addUsers)
	e.POST("/cats", addCats)
	e.POST("/dogs", addDogs)
	e.POST("/hamsters", addHamster)

	e.Logger.Fatal(e.Start(":1323"))
}
