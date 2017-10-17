package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"sync"

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

	mutex sync.Mutex
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

func initDB() {
	orderId = 0
	recentId = 0

	counter = make([]int, 114514)
	soldList = make([]string, 114514)

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

	variation = []Variation{
		Variation{artistName: "NHN48",
			ticketName:    "西武ドームライブ",
			variationName: "アリーナ席",
		},

		Variation{artistName: "NHN48",
			ticketName:    "西武ドームライブ",
			variationName: "スタンド席",
		},

		Variation{artistName: "NHN48",
			ticketName:    "東京ドームライブ",
			variationName: "アリーナ席",
		},

		Variation{artistName: "NHN48",
			ticketName:    "東京ドームライブ",
			variationName: "スタンド席",
		},

		Variation{artistName: "はだいろクローバーZ",
			ticketName:    "さいたまスーパーアリーナライブ",
			variationName: "アリーナ席",
		},

		Variation{artistName: "はだいろクローバーZ",
			ticketName:    "さいたまスーパーアリーナライブ",
			variationName: "スタンド席",
		},

		Variation{artistName: "はだいろクローバーZ",
			ticketName:    "横浜アリーナライブ",
			variationName: "アリーナ席",
		},

		Variation{artistName: "はだいろクローバーZ",
			ticketName:    "横浜アリーナライブ",
			variationName: "スタンド席",
		},

		Variation{artistName: "はだいろクローバーZ",
			ticketName:    "西武ドームライブ",
			variationName: "アリーナ席",
		},

		Variation{artistName: "はだいろクローバーZ",
			ticketName:    "西武ドームライブ",
			variationName: "スタンド席",
		},
	}
}

func get_recent_sold() string {
	ret := ""
	n := orderId - 10
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
	ret += artist[r.artistId-1].artistName
	ret += `</h2>`
	ret += `<ul>`
	for i := 0; i < len(artist[r.artistId-1].ticketIds); i++ {
		id := artist[r.artistId-1].ticketIds[i]
		ret += `<li class="ticket">`
		ret += `<a href="/ticket/`
		ret += itoa(id + 1)
		ret += `">`
		ret += artist[r.artistId-1].ticketNames[i]
		ret += `</a>残り<span class="count">`
		ret += itoa(4096*2 - (counter[id*2-1-1] + counter[id*2-1]))
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
	for i, art := range artist {
		ret += fmt.Sprintf(`<li><a href="/artist/%d"><span class="artist_name">%s</span></a></li>`, i+1, art.artistName)
	}
	return ret
}

func GenIndexHTML(r *Render) string {
	artlist := getArtistList()
	return fmt.Sprintf(`<h1>TOP</h1><ul>%s</ul>`, artlist)
}

func GenTicketHTML(r *Render) string {
	ret := fmt.Sprintf(`<h2> %s : %s </h2> <ul> `, ticket[r.ticketId-1].artistName, ticket[r.ticketId-1].ticketName)

	for i, v := range ticket[r.ticketId-1].variationIds {

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
	`, r.ticketId, v, ticket[r.ticketId-1].variationNames[i], v, 4096-counter[v])

	}
	ret += `</ul><h3>席状況</h3>`

	for i, v := range ticket[r.ticketId-1].variationIds {
		ret += fmt.Sprintf(` <h4>%s</h4> <table class="seats" data-variationid="%d"> `, ticket[r.ticketId-1].variationNames[i], v)

		for row := 0; row < 64; row++ {
			ret += `<tr>`
			for col := 0; col < 64; col++ {
				if row*64+col <= counter[v] {
					ret += fmt.Sprintf(`<td id="%2d-%2d" class="unavailable"></td>`, row, col)
				} else {
					ret += fmt.Sprintf(`<td id="%2d-%2d" class="available"></td>`, row, col)
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
	initDB()
	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		r := &Render{}
		return c.HTML(http.StatusOK, GenHTML(indexHTML, r))
	})

	e.GET("/artist/:artist_id", func(c echo.Context) error {
		r := &Render{
			artistId: atoi(c.Param("artist_id")),
		}
		return c.HTML(http.StatusOK, GenHTML(artistHTML, r))
	})

	e.GET("/ticket/:ticket_id", func(c echo.Context) error {
		r := &Render{
			ticketId: atoi(c.Param("ticket_id")),
		}
		return c.HTML(http.StatusOK, GenHTML(ticketHTML, r))
	})

	e.POST("/buy", func(c echo.Context) error {
		r := &Render{
			variationId: atoi(c.Param("variation_id")),
			memberId:    atoi(c.Param("member_id")),
		}

		mutex.Lock()
		orderId++
		if counter[r.variationId] == 4096 {
			return c.HTML(http.StatusOK, GenHTML(soldOutHTML, r))
		}
		counter[r.variationId]++
		ctr := counter[r.variationId]
		r.seatId = fmt.Sprintf("%02d-%02d", ctr/64, ctr%64)
		push(fmt.Sprintf("%s %s %s</td>\n<td class=\"recent_seat_id\">%s",
			variation[r.variationId].artistName, variation[r.variationId].ticketName, variation[r.variationId].variationName, r.seatId))

		csv += fmt.Sprintf("%d,%s,%s,%s\n",
			orderId, r.memberId, r.seatId, r.variationId, time.Now().Format("%Y-%m-%d %X"))

		mutex.Unlock()
		return c.HTML(http.StatusOK, GenHTML(completeHTML, r))
	})

	e.GET("/admin", func(c echo.Context) error {
		return c.HTML(http.StatusOK, GenHTML(adminHTML, &Render{}))
	})

	e.POST("/admin", func(c echo.Context) error {
		return c.Redirect(302, "/admin")
	})

	e.GET("/admin/order.csv", func(c echo.Context) error {
		return c.String(http.StatusOK, csv)
	})

	e.HTTPErrorHandler = func(err error, c echo.Context) {
		fmt.Print("404: " + c.Path())
	}

	e.Logger.Fatal(e.Start(":5000"))
}
