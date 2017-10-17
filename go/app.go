package main

import (
	"strconv"
	"net/http"
	"github.com/labstack/echo"
	"fmt"
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
	orderId   int
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


func initDB() {
	orderId = 0

	artist = []Artist{

		Artist{artistName: "NHN48",
			ticketNames: []string{"西武ドームライブ", "東京ドームライブ"},
			ticketIds:   []int{1, 2},
		},

		Artist{artistName: "はだいろクローバーZ",
			ticketNames: []string{"さいたまスーパーアリーナライブ", "横浜アリーナライブ", "西武ドームライブ"},
			ticketIds:   []int{3, 4, 5},
		},
	}
	ticket = []Ticket{

		Ticket{artistName: "NHN48",
			ticketName:     "西武ドームライブ",
			variationNames: []string{"アリーナ席", "スタンド席"},
			variationIds:   []int{1, 2},
		},

		Ticket{artistName: "NHN48",
			ticketName:     "東京ドームライブ",
			variationNames: []string{"アリーナ席", "スタンド席"},
			variationIds:   []int{3, 4},
		},

		Ticket{artistName: "はだいろクローバーZ",
			ticketName:     "さいたまスーパーアリーナライブ",
			variationNames: []string{"アリーナ席", "スタンド席"},
			variationIds:   []int{5, 6},
		},

		Ticket{artistName: "はだいろクローバーZ",
			ticketName:     "横浜アリーナライブ",
			variationNames: []string{"アリーナ席", "スタンド席"},
			variationIds:   []int{7, 8},
		},

		Ticket{artistName: "はだいろクローバーZ",
			ticketName:     "西武ドームライブ",
			variationNames: []string{"アリーナ席", "スタンド席"},
			variationIds:   []int{9, 10},
		},
	}

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
	for i := 0; i < len(artist); i++ {
		ret += fmt.Sprintf(`<li><span class="artist_name">%s</span></li>`,artist[i].artistName)
	}
	return ret
}

func GenIndexHTML(r *Render) string {
	artlist := getArtistList()
	return fmt.Sprintf(`<h1>TOP</h1><ul>%s</ul>`, artlist)
}

func GenSoldOutHTML(r *Render) string {
	return `<span class="result" data-result="failure">売り切れました。</span>`
}

func GenTicketHTML(r *Render) string {
	return ""
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
	initDB()
	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		r := &Render{}
		return c.HTML(http.StatusOK, GenHTML("index", r));
	})

	e.GET("/artist/:artist_id", func(c echo.Context) error {
		r := &Render{}
		r.artistId = atoi(c.Param("artist_id"))
		return c.HTML(http.StatusOK, GenHTML("artist", r));
	})

	e.GET("/ticket/:ticket_id", func(c echo.Context) error {
		r := &Render{}
		r.ticketId = atoi(c.Param("ticket_id"))
		return c.HTML(http.StatusOK, GenHTML("ticket", r));
	})

	e.POST("/buy", func(c echo.Context) error {
		r := Render{}
		r.variationId = atoi(c.Param("artist_id"))
		r.memberId = atoi(c.Param("member_id"))

		// 更新処理

		return nil
	})

	e.GET("admin", func(c echo.Context) error {
		r := &Render{}
		return c.HTML(http.StatusOK, GenHTML("admin", r));
	})

	e.GET("/admin/order.csv", func(c echo.Context) error {
		// content-type: text/csv?
		return c.String(http.StatusOK, csv);
	})

	e.Logger.Fatal(e.Start(":5000"))
}
