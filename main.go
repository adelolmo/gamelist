package main

import (
	"fmt"
	"log"
	"strings"
	"net/http"
	"database/sql"
	"encoding/json"

	_ "github.com/lib/pq"
	"github.com/gorilla/mux"
)

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("postgres", "postgres://gamelist:gamelist@localhost/gamelist?sslmode=disable")
	if err != nil {
		panic(err)
	}

	if err = db.Ping(); err != nil {
		panic(err)
	}
	fmt.Println("You connected to your database.")
}

type GameRow struct {
	Id         string
	Title      string
	Genres     string
	Company    sql.NullString
	Score      int
	ReleasedAt string
}

type Games struct {
	Games []Game `json:"games"`
}

type Game struct {
	Id         string `json:"id"`
	Title      string `json:"title"`
	Genres     []string `json:"genres"`
	Company    string `json:"company"`
	Score      int `json:"score"`
	ReleasedAt string `json:"releasedAt"`
	CoverUrl   string `json:"coverUrl"`
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/games", games).Methods("GET")
	router.HandleFunc("/games/{id}", game).Methods("GET")
	router.PathPrefix("/covers/").Handler(http.StripPrefix("/covers/",
		http.FileServer(http.Dir("./covers/"))))
	log.Fatal(http.ListenAndServe(":8080", router))
}

func games(w http.ResponseWriter, req *http.Request) {
	rows, err := db.Query("SELECT id, title, genres, company, score, released_at FROM games")
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	defer rows.Close()

	games := make([]Game, 0)

	for rows.Next() {
		gr := GameRow{}
		err := rows.Scan(&gr.Id, &gr.Title, &gr.Genres, &gr.Company, &gr.Score, &gr.ReleasedAt)
		if err != nil {
			http.Error(w, "{\"message\":\"Not found\",\"code\":404}", 404)
			return
		}
		genres := strings.Split(gr.Genres, ",")
		coverUrl := fmt.Sprintf("http://%s/covers/%s.jpeg", req.Host, gr.Id)
		games = append(games,
			Game{Id:gr.Id,
				Title:gr.Title,
				Genres:genres,
				Company:gr.Company.String,
				Score:gr.Score,
				ReleasedAt:gr.ReleasedAt,
				CoverUrl:coverUrl})

	}
	if err = rows.Err(); err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	json.NewEncoder(w).Encode(Games{Games:games})
}

func game(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	fmt.Println("GET " + params["id"])

	gr := GameRow{}
	row := db.QueryRow("SELECT id, title, genres, company, score, released_at FROM games where id=$1",
		params["id"])

	err := row.Scan(&gr.Id, &gr.Title, &gr.Genres, &gr.Company, &gr.Score, &gr.ReleasedAt) // order matters
	if err != nil {
		http.Error(w, "{\"message\":\"Not found\",\"code\":404}", 404)
		return
	}
	genres := strings.Split(gr.Genres, ",")
	coverUrl := fmt.Sprintf("http://%s/covers/%s.jpeg", req.Host, gr.Id)
	g := Game{Id:gr.Id,
		Title:gr.Title,
		Genres:genres,
		Company:gr.Company.String,
		Score:gr.Score,
		ReleasedAt:gr.ReleasedAt,
		CoverUrl:coverUrl}

	json.NewEncoder(w).Encode(g)
}