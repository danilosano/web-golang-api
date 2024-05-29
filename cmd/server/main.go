package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"

	"github.com/danilosano/web-golang-api/cmd/routes"
	"github.com/danilosano/web-golang-api/docs"
)

// @title MELI Bootcamp Fresh Food API
// @version 1.0
// @description This API handle MELI Fresh Food.
//
// @contact.name API Fresh Food.
// @contact.url https://developers.mercadolibre.com.ar/support
//
// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("error loading .env file: %s\n", err.Error())
	}
	db, err := sql.Open("mysql", fmt.Sprintf("%s?parseTime=true", os.Getenv("MYSQL_CONNECTION_STRING")))
	if err != nil {
		log.Fatalf("error opening database connection: %s\n", err.Error())
	}
	docs.SwaggerInfo.Host = os.Getenv("HOST")

	r := gin.Default()
	router := routes.NewRouter(r, db)
	router.MapRoutes()
	r.Run()
}
