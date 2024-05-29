package routes

import (
	"database/sql"

	"github.com/danilosano/web-golang-api/cmd/handler"
	"github.com/danilosano/web-golang-api/internal/customer"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Router interface {
	MapRoutes()
}

type router struct {
	eng *gin.Engine
	rg  *gin.RouterGroup
	db  *sql.DB
}

func NewRouter(eng *gin.Engine, db *sql.DB) Router {
	return &router{eng: eng, db: db}
}

func (r *router) MapRoutes() {
	r.setGroup()

	r.buildSwaggerRoutes()
	r.buildCustomerRoutes()
}

func (r *router) setGroup() {
	r.rg = r.eng.Group("/api/v1")
}

func (r *router) buildSwaggerRoutes() {
	r.rg.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

func (r *router) buildCustomerRoutes() {
	repo := customer.NewRepository(r.db)
	service := customer.NewService(repo)
	handler := handler.NewCustomerHandler(service)
	sections := r.rg.Group("/customers")
	{
		sections.POST("/", handler.Store)
		sections.GET("/", handler.GetAll)
		sections.GET("/:id", handler.Get)
		sections.PUT("/:id", handler.Update)
		sections.DELETE("/:id", handler.Delete)
	}
}
