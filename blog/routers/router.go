package routers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"


	_ "github.com/baxiang/soldiers-sortie/blog/docs"
	"github.com/baxiang/soldiers-sortie/blog/middleware"
	"github.com/baxiang/soldiers-sortie/blog/pkg/export"
	"github.com/baxiang/soldiers-sortie/blog/pkg/upload"
	"github.com/baxiang/soldiers-sortie/blog/routers/api"
	"github.com/baxiang/soldiers-sortie/blog/routers/api/v1"

)

// InitRouter initialize routing information
func InitRouter() *gin.Engine {
	r := gin.Default()

	r.StaticFS("/export", http.Dir(export.GetExcelFullPath()))
	r.StaticFS("/upload/images", http.Dir(upload.GetImageFullPath()))

	r.POST("/auth", api.GetAuth)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.POST("/upload", api.UploadImage)

	apiV1 := r.Group("/api/v1")
	apiV1.Use(middleware.JWT())
	{
		//获取标签列表
		apiV1.GET("/tags", v1.GetTags)
		//新建标签
		apiV1.POST("/tags", v1.AddTag)
		//更新指定标签
		apiV1.PUT("/tags/:id", v1.EditTag)
		//删除指定标签
		apiV1.DELETE("/tags/:id", v1.DeleteTag)
		//导出标签
		r.POST("/tags/export", v1.ExportTag)
		//导入标签
		r.POST("/tags/import", v1.ImportTag)

		//获取文章列表
		apiV1.GET("/articles", v1.GetArticles)
		//获取指定文章
		apiV1.GET("/articles/:id", v1.GetArticle)
		//新建文章
		apiV1.POST("/articles", v1.AddArticle)
		//更新指定文章
		apiV1.PUT("/articles/:id", v1.EditArticle)
		//删除指定文章
		apiV1.DELETE("/articles/:id", v1.DeleteArticle)
	}

	return r
}