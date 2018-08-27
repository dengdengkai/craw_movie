package routers

import (
	"craw_movie/controllers"

	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/craw_movie", &controllers.CrawMovieController{}, "*:CrawlMovie")
}
