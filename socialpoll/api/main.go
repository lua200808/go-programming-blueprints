package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/tylerb/graceful.v1"
)

func main() {
	var (
		addr  = flag.String("addr", ":8080", "エンドポイントのアドレス")
		mongo = flag.String("mongo", "localhost", "MongoDB のアドレス")
	)
	flag.Parse()
	log.Println("MongoDB に接続します", *mongo)
	db, err := mgo.Dial(*mongo)
	if err != nil {
		log.Fatalln("MongoDB への接続に失敗しました:", err)
	}
	defer db.Close()
	mux := http.NewServeMux()
	mux.HandlerFunc("/polls/", withCORS(withVars(withData(db, withAPIKey(handlePolls)))))
	log.Println("Web サーバーを開始します:", *addr)
	graceful.Run(*addr, 1*time.Second, mux)
	log.Println("停止します...")
}

func withAPIKey(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !isValidAPIKey(r.URL.Query().Get("key")) {
			respondErr(w, r, http.StatusUnauthorized, "不正な API キーです")
			return
		}
		fn(w, r)
	}
}
func isValidAPIKey(key string) bool {
	return key == "abc123"
}

func withData(b *mgo.Session, f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		thisDb := d.Copy()
		defer thisDb.Close()
		SetVar(r, "db", thisDb.DB("ballots"))
		f(w, r)
	}
}

func withVars(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		OpenVars(r)
		defer CloseVars(r)
		fn(w, r)
	}
}

func withCORS(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Expose-Headers", "Location")
		fn(w, r)
	}
}
