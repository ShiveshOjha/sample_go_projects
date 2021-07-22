package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

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

func addHamster(c echo.Context) error { //shortest
	hamster := Hamster{}

	defer c.Request().Body.Close()

	err := c.Bind(&hamster)
	if err != nil {
		log.Printf("Failed processing addHamster request: %s", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	log.Printf("This is your hamster: %v", hamster)
	return c.String(http.StatusOK, "we got your hamster")
}

func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello from web side!")
}

func addDog(c echo.Context) error {
	dog := Dog{}

	defer c.Request().Body.Close()

	err := json.NewDecoder(c.Request().Body).Decode(&dog) // same as unmarshalling
	if err != nil {
		log.Printf("Failed processing addDog request: %s", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	log.Printf("This is your dog: %v", dog)
	return c.String(http.StatusOK, "we got your dog")

}

func getCats(c echo.Context) error {
	catName := c.QueryParam("name")
	catType := c.QueryParam("type")

	dataType := c.Param("data")

	if dataType == "string" {
		return c.String(http.StatusOK, fmt.Sprintf("Your cat's name is: %s and his type is %s\n", catName, catType))
	}

	if dataType == "json" {
		return c.JSON(http.StatusOK, map[string]string{
			"name": catName,
			"type": catType,
		})
	}

	return c.JSON(http.StatusBadRequest, map[string]string{
		"error": "Specify data type as json or string",
	})

}

func addCat(c echo.Context) error {
	cat := Cat{}

	defer c.Request().Body.Close()

	b, err := ioutil.ReadAll(c.Request().Body) // reading request body
	if err != nil {
		log.Fatalf("Failed reading the request body: %v", err)
		return c.String(http.StatusInternalServerError, "")
	}

	err = json.Unmarshal(b, &cat) // takes json response from b and converts to string and stores in cat
	if err != nil {
		log.Printf("Failed unmarshaling in addCats: %v", err)
		return c.String(http.StatusInternalServerError, "")
	}

	log.Printf("This is your cat: %v", cat)
	return c.String(http.StatusOK, "we got your cat!")
}

func main() {
	fmt.Println("Welcome to the server")

	e := echo.New() // creates router

	e.GET("/", hello)
	e.GET("/cats/:data", getCats)

	e.POST("/cats", addCat)
	e.POST("/dogs", addDog)
	e.POST("/hamsters", addHamster)

	e.Start(":8181")
}
