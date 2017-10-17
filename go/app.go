package main

import (
	"strconv"

	"github.com/labstack/echo"
)

type Artist struct {
	artistName  string
	ticketNames []string
	ticketIds   []int
}

type Ticket struct {
	artistName     string
	ticketName     string
	variationNames []string
	variationIds   []int
}

var (
	counter   []int
	soldList  []string
	recentId  int
	orderdId  int
	csv       string
	emptySold bool

	inArtistId int
	inTicketId int

	artist []Artist // artist[artist_id] = Artist
	ticket []Ticket // ticket[ticket_id] = Ticket
)

func itoa(a int) string {
	return strconv.Itoa(a)
}
func atoi(a string) int {
	b, _ := strconv.Atoi(a)
	return b
}

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
	ret := `
	<ul>
	<li>
	<a href="/admin/order.csv">注文CSV</a>
	</li>
	<li>
	<form method="POST">
	<input type="submit" value="データ初期化" />
	</form>
	</li>
	</ul>
	`
	return ret
}

func GenArtistHTML() string {
	ret := `<h2>`
	ret += artist[inArtistId].artistName
	ret += `</h2>`
	ret += `<ul>`
	for i := 0; i < len(ticket); i++ {
		ret += `<li class="ticket">`
		ret += `<a href="/ticket/`
		ret += itoa(artist[inArtistId].ticketIds[i])
		ret += `">`
		ret += artist[inArtistId].ticketNames[i]
		ret += `</a>残り<span class="count">`
		ret += itoa(counter[artist[inArtistId].ticketIds[i]])
		ret += `</span>枚`
	}
	ret += `</li>`
	return ret
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
		variation_id := atoi(c.Param("artist_id"))
		member_id := atoi(c.Param("member_id"))

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
