package main

import (
	"strconv"

	"fmt"
	"net/http"

	"time"

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
	seatId      string
}

type Variation struct {
	artistName    string
	ticketName    string
	variationName string
}

var (
	counter  []int
	soldList []string
	recentId int
	orderId  int
	csv      string

	artist    []Artist    // artist[artist_id] = Artist
	ticket    []Ticket    // ticket[ticket_id] = Ticket
	variation []Variation // variation[variation_id] = Variaion
)

const (
	adminHTML    = "admin"
	artistHTML   = "artist"
	completeHTML = "complete"
	indexHTML    = "index"
	soldOutHTML  = "soldout"
	ticketHTML   = "ticket"
)

func itoa(a int) string {
	return strconv.Itoa(a)
}
func atoi(a string) int {
	b, _ := strconv.Atoi(a)
	return b
}
func push(s string) {
	soldList[recentId] = s
	recentId++
}

func get_recent_sold() string {
	ret := ""
	n := recentId - 10
	if n < 0 {
		n = 0
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
		id := artist[r.artistId].ticketIds[i]
		ret += `<li class="ticket">`
		ret += `<a href="/ticket/`
		ret += itoa(id)
		ret += `">`
		ret += artist[r.artistId].ticketNames[i]
		ret += `</a>残り<span class="count">`
		ret += itoa(4096*2 - (counter[id*2-1] + counter[id*2]))
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
	ret += r.seatId
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

func GenSoldOutHTML(r *Render) string {
	return `<span class="result" data-result="failure">売り切れました。</span>`
}

func GenTicketHTML(r *Render) string {
	return ""
}

func GenHTML(content_name string, r *Render) string {
	res := `<!DOCTYPE html> <html> <head>	<title>isucon 2</title>	<meta charset="utf-8">	<link type="text/css" rel="stylesheet" href="/css/ui-lightness/jquery-ui-1.8.24.custom.css">	<link type="text/css" rel="stylesheet" href="/css/isucon2.css">	<script type="text/javascript" src="/js/jquery-1.8.2.min.js"></script>	<type="text/javascript" src="/js/jquery-ui-1.8.24.custom.min.js"></script>	<script type="text/javascript" src="/js/isucon2.js"></script>	</head>	<body>	<header>	<a href="/">	<img src="/images/isucon_title.jpg">	</a>	</header>	<div id="sidebar">`
	if orderId > 0 {
		res += `<table><tr><th colspan="2">最近購入されたチケット</th></tr>`
		res += get_recent_sold()
		res += `</table>`
	}
	res += `</div><div id="content">`
	switch content_name {
	case adminHTML:
		res += GenAdminHTML(r)
	case artistHTML:
		res += GenArtistHTML(r)
	case completeHTML:
		res += GenCompleteHTML(r)
	case indexHTML:
		res += GenIndexHTML(r)
	case soldOutHTML:
		res += GenSoldOutHTML(r)
	case ticketHTML:
		res += GenTicketHTML(r)
	}
	res += `</div></body></html>`
	return res
}

func main() {
	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		r := &Render{}
		return c.String(http.StatusOK, GenHTML(indexHTML, r))
	})

	e.GET("/artist/<:artist_id>", func(c echo.Context) error {
		r := &Render{
			artistId: atoi(c.Param("artist_id")),
		}
		return c.String(http.StatusOK, GenHTML(artistHTML, r))
	})

	e.GET("/ticket/<:ticket_id>", func(c echo.Context) error {
		r := &Render{
			ticketId: atoi(c.Param("ticket_id")),
		}
		return c.String(http.StatusOK, GenHTML(ticketHTML, r))
	})

	e.POST("/buy", func(c echo.Context) error {
		r := &Render{
			variationId: atoi(c.Param("variation_id")),
			memberId:    atoi(c.Param("member_id")),
		}

		orderId++
		if counter[r.variationId] == 4096 {
			return c.String(http.StatusOK, GenHTML(soldOutHTML, r))
		}
		counter[r.variationId]++
		ctr := counter[r.variationId]
		r.seatId = fmt.Sprint("%02d-%02d", ctr/64, ctr%64)

		push(fmt.Sprint("%s %s %s</td>\n<td class=\"recent_seat_id\">%s",
			variation[r.variationId].artistName, variation[r.variationId].ticketName, variation[r.variationId].variationName, r.seatId))

		csv += fmt.Sprint("%d,%s,%s,%s\n",
			orderId, r.memberId, r.seatId, r.variationId, time.Now().Format("%Y-%m-%d %X"))

		return c.String(http.StatusOK, GenHTML(completeHTML, r))
	})

	e.GET("admin", func(c echo.Context) error {
		return c.String(http.StatusOK, GenHTML(adminHTML, &Render{}))
	})

	e.GET("/admin/order.csv", func(c echo.Context) error {
		return c.String(http.StatusOK, csv)
	})

	e.Logger.Fatal(e.Start(":5000"))
}
