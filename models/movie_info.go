// movie_info
package models

import (
	"regexp" //导入正则表达式所需要的包
	"strings"

	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

var (
	db orm.Ormer
)

type MovieInfo struct {
	Id                   int64
	Movie_id             int64
	Movie_name           string
	Movie_pic            string
	Movie_director       string
	Movie_writer         string
	Movie_country        string
	Movie_language       string
	Movie_main_character string
	Movie_type           string
	Movie_on_time        string
	Movie_span           string
	Movie_grade          string
	_Create_time         string
	//与数据库的字段相对应(第一个字母大写表示可供外部使用，在进入数据库时所有大写会变为小写，)
}

func init() {
	orm.Debug = true // 是否开启调试模式 调试模式下会打印出sql语句
	//注册mysql
	orm.RegisterDataBase("default", "mysql", "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8", 30)
	orm.RegisterModel(new(MovieInfo))
	db = orm.NewOrm()
}
func AddMovie(movie_info *MovieInfo) (int64, error) {
	movie_info.Id = 0
	id, err := db.Insert(movie_info)
	return id, err

}

//利用正则匹配将导演匹配出来
func GetMovieDirector(movieHtml string) string {
	if movieHtml == "" {
		return ""
	}

	//注意时1前面的符号
	reg := regexp.MustCompile(`<a.*?rel="v:directedBy">(.*?)</a>`)
	// 负数表示不限次数
	result := reg.FindAllStringSubmatch(movieHtml, -1)
	if len(result) == 0 {
		return ""
	}
	return string(result[0][1])
}

//利用正则匹配将电影名匹配出来
func GetMovieName(movieHtml string) string {
	if movieHtml == "" {
		return ""
	}
	//注意时1前面的符号
	reg := regexp.MustCompile(`<span\s*property="v:itemreviewed">(.*?)</span>`)
	// 负数表示不限次数。
	result := reg.FindAllStringSubmatch(movieHtml, -1)
	if len(result) == 0 {
		return ""
	}

	return string(result[0][1])
}

//利用正则将演员匹配出来
func GetMovieMainCharacters(movieHtml string) string {

	reg := regexp.MustCompile(`<a.*?rel="v:starring">(.*?)</a>`)
	//加一个问号代表贪婪匹配，尽量少的匹配
	// 负数表示不限次数。
	result := reg.FindAllStringSubmatch(movieHtml, -1)
	if len(result) == 0 {
		return ""
	}
	mainCharacters := ""
	for _, v := range result {
		mainCharacters += v[1] + "/"

		//每一次返回是一个下表为1时是结果
	}
	//去掉末尾的斜杠
	return strings.Trim(mainCharacters, "/")

}

//利用正则匹配将电影评分匹配出来
func GetMovieGrade(movieHtml string) string {
	if movieHtml == "" {
		return ""
	}
	//注意时1前面的符号
	reg := regexp.MustCompile(`<strong.*?property="v:average">(.*?)</strong>`)

	result := reg.FindAllStringSubmatch(movieHtml, -1)
	if len(result) == 0 {
		return ""
	}
	return string(result[0][1])
}

//利用正则将电影类型匹配出来
func GetMovieGenre(movieHtml string) string {
	if movieHtml == "" {
		return ""
	}
	//注意时1前面的符号
	reg := regexp.MustCompile(`<span\s+property="v:genre">(.*?)</span>`)

	result := reg.FindAllStringSubmatch(movieHtml, -1)
	if len(result) == 0 {
		return ""
	}
	movieTypes := ""
	for _, v := range result {
		movieTypes += v[1] + "/"

		//每一次返回是一个下表为1时是结果
	}
	//去掉末尾的斜杠
	return strings.Trim(movieTypes, "/")

}

func GetMovieOnTime(movieHtml string) string {
	if movieHtml == "" {
		return ""
	}
	//注意时1前面的符号  下面有一个转移符号\(  为了匹配（
	reg := regexp.MustCompile(`<span.*?property="v:initialReleaseDate".*?>(.*?)\(.*</span>`)

	result := reg.FindAllStringSubmatch(movieHtml, -1)
	if len(result) == 0 {
		return ""
	}
	return string(result[0][1])
}

//获取电影时长
func GetMovieRunningTime(movieHtml string) string {
	if movieHtml == "" {
		return ""
	}
	//注意时1前面的符号
	reg := regexp.MustCompile(`<span.*?property="v:runtime".*?>(.*?)</span>`)
	//加一个问号代表贪婪匹配，尽量少的匹配
	// 负数表示不限次数。
	result := reg.FindAllStringSubmatch(movieHtml, -1)
	if len(result) == 0 {
		return ""
	}
	return string(result[0][1])
}

//获取编剧信息
func GetMovieWriter(movieHtml string) string {
	if movieHtml == "" {
		return ""
	}

	reg := regexp.MustCompile(`<a.*?href="/celebrity/.*?/">(.*?)</a>`)

	result := reg.FindAllStringSubmatch(movieHtml, -1)
	if len(result) == 0 {
		return ""
	}
	writers := ""
	for _, v := range result {
		writers += v[1] + "/"

		//每一次返回是一个下表为1时是结果
	}
	//去掉末尾的斜杠
	return strings.Trim(writers, "/")
}

//中文匹配[\u4e00-\u9fa5]  go语言中 [\p{han}]

//获取电影地区
func GetMovieCountry(movieHtml string) string {
	if movieHtml == "" {
		return ""
	}
	//注意时1前面的符号
	reg := regexp.MustCompile(`[\p{Han}]{4}/.*:</span>\s?(.*?)<`)
	//加一个问号代表贪婪匹配，尽量少的匹配
	// 负数表示不限次数。
	result := reg.FindAllStringSubmatch(movieHtml, -1)
	if len(result) == 0 {
		return ""
	}
	return string(result[0][1])
}

/*
//获取语言
func GetMovieLanguage(movieHtml string) string {
	if movieHtml == "" {
		return ""
	}
	//注意时1前面的符号
	reg := regexp.MustCompile(`[\p{Han}]{2}:</span>\s?([[\p{Han}]{5}]|[\p{Han}]{2}])<`)
	//加一个问号代表贪婪匹配，尽量少的匹配
	// 负数表示不限次数。
	result := reg.FindAllStringSubmatch(movieHtml, -1)

	return string(result[0][1])
}
*/
//爬取不停获取电影url
func GetMovieUrls(movieHtml string) []string {
	reg := regexp.MustCompile(`<a.*?href="(https://movie.douban.com/.*?)"`)
	result := reg.FindAllStringSubmatch(movieHtml, -1)

	var movieSets []string
	for _, v := range result {

		movieSets = append(movieSets, v[1])
	}

	return movieSets
}
