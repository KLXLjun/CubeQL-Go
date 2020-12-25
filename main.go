package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	bloomfilter "cubeqL-go/bloomfilter"
)

var filter = bloomfilter.NewBloomFilter()
var cylinder = make([]*urlArrayType, 0)

func main() {
	e := echo.New()
	//Gzip 压缩
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
	}))

	//日志输出
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${remote_ip} [${time_rfc3339}] ${method} ${status} ${latency_human} ${uri} ${user_agent}\n",
	}))

	//主页
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "CubeQL-Go is Run!")
	})

	e.GET("/get", getValue)
	e.GET("/set", setValue)
	e.GET("/filter_set", filterValueSet)
	e.GET("/filter_contain", filterValueContain)

	//fmt.Println(filter.Funcs[1].Seed)
	//fmt.Println(filter.Set.Count())

	//启动
	e.Logger.Fatal(e.Start(":1278"))
}

type urlArrayType struct {
	Typ     string
	Content string
}

func getValue(c echo.Context) error {
	return c.JSON(http.StatusOK, &cylinder)
}

func setValue(c echo.Context) error {
	typ := c.QueryParam("typ")
	content := c.QueryParam("url")
	t := new(urlArrayType)
	t.Typ = typ
	t.Content = content
	cylinder = append(cylinder, t)

	return c.String(http.StatusOK, "ok")
}

func filterValueSet(c echo.Context) error {
	url := c.QueryParam("url")
	filter.Add(url)
	return c.String(http.StatusOK, "ok")
}

func filterValueContain(c echo.Context) error {
	url := c.QueryParam("url")
	if filter.Contains(url) {
		return c.JSON(http.StatusOK, 1)
	}
	return c.JSON(http.StatusOK, 0)
}

func delValue(c echo.Context) error {
	cylinder = make([]*urlArrayType, 0)
	return c.String(http.StatusOK, "ok")
}
