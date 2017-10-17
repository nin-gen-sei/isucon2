package main

import (
	"strconv"

	"fmt"
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

type Render struct {
	artistId    int
	ticketId    int
	variationId int
	memberId    int
	seatId      int
}

var (
	counter   []int
	soldList  []string
	recentId  int
	orderdId  int
	csv       string
	emptySold bool

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

func GenAdminHTML(r *Render) string {
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

func GenArtistHTML(r *Render) string {
	ret := `<h2>`
	ret += artist[r.artistId].artistName
	ret += `</h2>`
	ret += `<ul>`
	for i := 0; i < len(ticket); i++ {
		ret += `<li class="ticket">`
		ret += `<a href="/ticket/`
		ret += itoa(artist[r.artistId].ticketIds[i])
		ret += `">`
		ret += artist[r.artistId].ticketNames[i]
		ret += `</a>残り<span class="count">`
		ret += itoa(counter[artist[r.artistId].ticketIds[i]])
		ret += `</span>枚`
	}
	ret += `</li>`
	return ret
}

func GenCompleteHTML(r *Render) string {
	ret := `<h2>予約完了</h2>`
	ret += `会員ID:<span class="member_id">`
	ret += itoa(r.memberId)
	ret += `</span>で<span class="result" data-result="success">&quot;<span class="seat">`
	ret += itoa(r.seatId)
	ret += `</span>&quot;の席を購入しました。</span>`
	return ret
}

func getArtistList() string {
	ret := ""
	for i := 0; i <= len(artist); i++ {
		ret += fmt.Sprint(`<li><span class="artist_name">%s</span></li>`, artist[i].artistName)
	}
	return ret
}

func GenIndexHTML(r *Render) string {
	artlist := getArtistList()
	return fmt.Sprint(`<h1>TOP</h1><ul>%s</ul>`, artlist)
}

func GenTicketHTML(r *Render) string {
	ret := fmt.Sprintf(`<h2> %s : %s </h2> <ul> `, ticket[r.ticketId].artistName, ticket[r.ticketId].ticketName)

	for i, v := range ticket[r.ticketId].variationIds {

		ret += fmt.Sprintf(`
	  <li class="variation">
	  <form method="POST" action="/buy">
	  <input type="hidden" name="ticket_id" value="%d">
	  <input type="hidden" name="variation_id" value="%d">
	  <span class="variation_name">%s </span> 残り<span class="vacancy" id="vacancy_%d">%d</span>席
	  <input type="text" name="member_id" value="">
	  <input type="submit" value="購入">
	  </form>
	  </li>
	`, r.ticketId, v, ticket[r.ticketId].variationNames[i], v, 4096-counter[v])

	}
	ret += `</ul><h3>席状況</h3>`

	for i, v := range ticket[r.ticketId].variationIds {
		ret += fmt.Sprintf(` <h4>%s</h4> <table class="seats" data-variationid="%d"> `, ticket[r.ticketId].variationNames[i], v)

		for row := 0; row < 64; row++ {
			ret += `<tr>`
			for col := 0; col < 64; col++ {
				if row*64+col <= counter[v] {
					ret += fmt.Sprintf(`<td id="%2d-%2d" class="available"></td>`, row, col)
				} else {
					ret += fmt.Sprintf(`<td id="%2d-%2d" class="unavailable"></td>`, row, col)
				}
			}
			ret += "</tr>"
		}
		ret += "</table>"
	}

	return ret
}

func GenSoldOutHTML(r *Render) string {
	return `<span class="result" data-result="failure">売り切れました。</span>`

}

func GenHTML(content_name string, r *Render) string {
	res := `<!DOCTYPE html> <html> <head>	<title>isucon 2</title>	<meta charset="utf-8">	<link type="text/css" rel="stylesheet" href="/css/ui-lightness/jquery-ui-1.8.24.custom.css">	<link type="text/css" rel="stylesheet" href="/css/isucon2.css">	<script type="text/javascript" src="/js/jquery-1.8.2.min.js"></script>	<type="text/javascript" src="/js/jquery-ui-1.8.24.custom.min.js"></script>	<script type="text/javascript" src="/js/isucon2.js"></script>	</head>	<body>	<header>	<a href="/">	<img src="/images/isucon_title.jpg">	</a>	</header>	<div id="sidebar">`
	if !emptySold {
		res += `<table><tr><th colspan="2">最近購入されたチケット</th></tr>`
		res += get_recent_sold()
		res += `</table>`
	}
	res += `</div><div id="content">`
	switch content_name {
	case "admin":
		res += GenAdminHTML(r)
	case "artist":
		res += GenArtistHTML(r)
	case "complete":
		res += GenCompleteHTML(r)
	case "index":
		res += GenIndexHTML(r)
	case "soldout":
		res += GenSoldOutHTML(r)
	case "ticket":
		res += GenTicketHTML(r)
	}
	res += `</div></body></html>`
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
		r := Render{}
		r.variationId = atoi(c.Param("artist_id"))
		r.memberId = atoi(c.Param("member_id"))

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
