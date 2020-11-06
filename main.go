package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"
)

type accountBook struct {
	db *sql.DB
}

func main() {
	//DB接続
	db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/accountBook?parseTime=true&loc=Asia%2FTokyo")
	if err != nil {
		panic(err)
	}

	db.SetConnMaxLifetime(time.Minute * 3) //ドライバーによる接続を閉じる時間の設定
	db.SetMaxOpenConns(5)                  //同時接続の上限
	db.SetMaxIdleConns(5)                  //コネクションプールのアイドル数

	defer db.Close()

	ab := new(accountBook)
	ab.db = db

	//ハンドラの登録
	http.Handle("/resources/", http.StripPrefix("/resources/", http.FileServer(http.Dir("resources/"))))
	http.HandleFunc("/", ab.getListHandler)
	http.HandleFunc("/addPurchasedItem", ab.addPurchasedItemHandler)
	http.HandleFunc("/updatePurchasedItem", ab.updatePurchasedItemHandler)
	http.HandleFunc("/addUsedItem", ab.addUsedItemHandler)

	//サーバの起動
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}
