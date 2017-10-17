package main

import (
	"strconv"

	"github.com/labstack/echo"
)

var (
	counter   []uint32
	soldList  []string
	recentId  uint32
	orderdId  uint32
	csv       string
	emptySold bool
)

func get_recent_sold() string {
	ret := ""
	n := len(soldList)
	if n-10 < 0 {
		n = 0
	} else {
		n = n - 10
	}
	recent_sold := soldList[n:]
	for _, s := range recent_sold {
		ret += "<tr><td class=\"recent_variation\">"
		ret += s
		ret += "</td>\n</tr>"
	}
	return ret
}

func GenAdminHTML() string {
	return ""
}

func GenArtistHTML() string {
	return ""
}

func GenCompleteHTML() string {
	return ""
}

func GenIndexHTML() string {
	return ""
}

func GenSoldOutHTML() string {
	return ""
}

func GenTicketHTML() string {
	return ""
}

func GenHTML(content_name string) string {
	res := ` <!DOCTYPE html> <html> <head>	<title>isucon 2</title>	<meta charset="utf-8">	<link type="text/css" rel="stylesheet" href="/css/ui-lightness/jquery-ui-1.8.24.custom.css">	<link type="text/css" rel="stylesheet" href="/css/isucon2.css">	<script type="text/javascript" src="/js/jquery-1.8.2.min.js"></script>	<type="text/javascript" src="/js/jquery-ui-1.8.24.custom.min.js"></script>	<script type="text/javascript" src="/js/isucon2.js"></script>	</head>	<body>	<header>	<a href="/">	<img src="/images/isucon_title.jpg">	</a>	</header>	<div id="sidebar">`
	if !emptySold {
		res += `<table>
	<tr><th colspan="2">最近購入されたチケット</th></tr>`
		res += get_recent_sold()
		res += `</table>`
	}
	res += `
	</div>
	<div id="content">
	`
	switch content_name {
	case "admin":
	case "artist":
	case "complete":
	case "index":
	case "soldout":
	case "ticket":
	}
	res += `
	</div>
	</body>
	</html>
	`
	return res
}

func main() {
	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return nil
	})

	e.GET("/artist/<int:artist_id>", func(c echo.Context) error {

		return nil
	})

	e.GET("/ticket/<int:ticket_id>", func(c echo.Context) error {
		return nil
	})

	e.POST("/buy", func(c echo.Context) error {
		variation_id, _ := strconv.Atoi(c.Param("artist_id"))
		member_id, _ := strconv.Atoi(c.Param("member_id"))

		// 更新処理

		return nil
	})

	e.GET("admin", func(c echo.Context) error {
		return nil
	})

	e.GET("/admin/order.csv", func(c echo.Context) error {
		return nil
	})

	e.Logger.Fatal(e.Start(":5000"))
}
