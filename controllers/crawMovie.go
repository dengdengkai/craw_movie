package controllers

import (
	"craw_movie/models"
	//	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/httplib"
	//"strings"
	"time"
)

type CrawMovieController struct {
	beego.Controller
}

func (c *CrawMovieController) CrawlMovie() {
	//利用结构体变量定义一个变量，然后调用方法，用内部变量存储获取的信息
	var movieInfo models.MovieInfo
	//连接redis
	models.ConnectRedis("127.0.0.1:6379")

	//请求访问页面
	//爬虫入口url
	sUrl := "https://movie.douban.com/subject/26985127/"

	//现将它加入redis的队列
	models.PutinQueue(sUrl)

	for {

		//获取队列长度，第二次执行此循环时不会为空，下方有一个获取的代码段
		length := models.GetQueueLength()
		if length == 0 {
			break //如果url队列为空，则退出当前循环
		}
		//消费者从队列获取Url
		sUrl = models.PopfromQueue()
		//判断电影的url是否被访问
		if models.IsVisit(sUrl) {
			continue
		}

		req := httplib.Get(sUrl)

		//req := httplib.Get("https://movie.douban.com/subject/26985127/")
		sMovieHtml, err := req.String() //获取页面信息(详见爬虫分析)转换为字符串
		if err != nil {
			panic(err)

		}

		//通过获取名字来判断它是不是电影
		movieInfo.Movie_name = models.GetMovieName(sMovieHtml)
		//如果它是电影，则获取其他信息，并且加入mysql数据库
		if movieInfo.Movie_name != "" {
			//获取电影信息
			movieInfo.Movie_director = models.GetMovieDirector(sMovieHtml)
			movieInfo.Movie_main_character = models.GetMovieMainCharacters(sMovieHtml)
			movieInfo.Movie_type = models.GetMovieGenre(sMovieHtml)
			movieInfo.Movie_grade = models.GetMovieGrade(sMovieHtml)
			movieInfo.Movie_on_time = models.GetMovieOnTime(sMovieHtml)
			movieInfo.Movie_span = models.GetMovieRunningTime(sMovieHtml) //时长
			movieInfo.Movie_writer = models.GetMovieWriter(sMovieHtml)
			movieInfo.Movie_country = models.GetMovieCountry(sMovieHtml)
			//加入mysql数据库
			models.AddMovie(&movieInfo)

		}
		//提取该页面的所有连接
		//urls为字符串数组类型
		urls := models.GetMovieUrls(sMovieHtml) //urls变量用于存取调用获取的电影地址
		for _, url := range urls {

			//调用获放入队列的函数，让url进入队列
			models.PutinQueue(url)
			//爬取结束打印提示
			c.Ctx.WriteString("<br>" + url + "</br>")

			//遍历字符串数组将每一个电影地址加入队列，执行程序之前 先启动redis（redis-server）
			//之后链接到redis (redis-cli) ,在浏览器输入网址，在命令行输入keys *  发现多了一个url_queue  执行lrange url_queue 0 -1

		}

		//sUrl应当记录到访问set中,表明已经访问过了
		models.AddToSet(sUrl)

		//为了防止爬取速度过快，每次爬完一部电影休息一秒
		time.Sleep(time.Second)

	}

	//movieInfo.Movie_language = models.GetMovieLanguage(sMovieHtml)
	//分析，由定义的结构体变量来接受捕获的内容

	//id, _ := models.AddMovie(&movieInfo)

	//返回的id类型是int64
	//	c.Ctx.WriteString(fmt.Sprintf("%v", id))

	/*  //用于在web界面打印出爬取的信息，便于查看结果是否为所需要的，然后便于修改代码调试
	c.Ctx.WriteString(models.GetMovieName(sMovieHtml) + "|")
	c.Ctx.WriteString(models.GetMovieDirector(sMovieHtml) + "|")
	c.Ctx.WriteString(fmt.Sprintf("%v", models.GetMovieMainCharacters(sMovieHtml)) + "|")
	c.Ctx.WriteString(models.GetMovieGrade(sMovieHtml) + "|")
	c.Ctx.WriteString(fmt.Sprintf("%v", models.GetMovieGenre(sMovieHtml)) + "|")
	c.Ctx.WriteString(models.GetMovieOnTime(sMovieHtml) + "|")
	c.Ctx.WriteString(models.GetMovieRunningTime(sMovieHtml) + "|")

	c.Ctx.WriteString(fmt.Sprintf("%v", models.GetMovieWriter(sMovieHtml)) + "|")
	c.Ctx.WriteString(models.GetMovieCountry(sMovieHtml) + "|")
	//c.Ctx.WriteString(models.GetMovieLanguage(sMovieHtml) + "|")

	*/

	//爬取结束打印提示
	c.Ctx.WriteString("end of craw_movie !")

}
