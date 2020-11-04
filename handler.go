package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"strconv"
	"time"
)

func (ab *accountBook) getListHandler(w http.ResponseWriter, r *http.Request) {
	//外部テンプレートファイルの読み込み
	data, err := ioutil.ReadFile("template/index.html")
	if err != nil {
		fmt.Fprintf(os.Stderr, "read file failed : %s\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	var listTmpl = template.Must(template.New("PurchasedList").
		Funcs(template.FuncMap{
			"formatDate": formatDate,
			"formatYen":  formatYen,
		}).
		Parse(string(data)))

	//データの取得
	PurchasedItems, err := ab.getPurchasedItems()
	if err != nil {
		fmt.Fprintf(os.Stderr, "getPurchasedItems() failed : %s\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 取得したitemsをテンプレートに埋め込む
	if err := listTmpl.Execute(w, PurchasedItems); err != nil {
		fmt.Fprintf(os.Stderr, "execute tamplate failed : %s\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func formatDate(t time.Time) string {
	return t.Format("2006/01/02")
}

func formatYen(str string) string {
	return fmt.Sprintf("¥ %10s", str)
}

//追加ボタン押下時のハンドラ
func (ab *accountBook) addPurchasedItemHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		code := http.StatusMethodNotAllowed
		http.Error(w, http.StatusText(code), code)
		return
	}

	//入力フォームの情報取得
	month := r.FormValue("month")
	purchasedOn, err := time.Parse("2006-01-02", r.FormValue("purchasedOn"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	shop := r.FormValue("shop")
	category := r.FormValue("category")
	itemName := r.FormValue("itemName")
	quantity, err := strconv.Atoi(r.FormValue("quantity"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var excTax int
	if r.FormValue("excTax") != "" {
		excTax, err = strconv.Atoi(r.FormValue("excTax"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	var tax float64
	if r.FormValue("tax") != "" {
		tax, err = strconv.ParseFloat(r.FormValue("tax"), 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	var incTax int
	if r.FormValue("incTax") != "" {
		incTax, err = strconv.Atoi(r.FormValue("incTax"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	//税抜、税込どちらの金額も入力されていた場合の、矛盾チェック
	if excTax == 0 && incTax == 0 {
		http.Error(w, "金額（税抜）、金額（税込）のどちらか入力してください", http.StatusBadRequest)
		return
	}

	if excTax != 0 && tax == 0 {
		http.Error(w, "金額（税抜）を入力した場合には、税金のどちらかを選択してください", http.StatusBadRequest)
		return
	}

	if excTax != 0 && incTax != 0 {
		if math.Abs(float64(excTax)-(float64(incTax)/(1.0+tax))) > 0.0001 {
			http.Error(w, "入力された金額（税抜）と金額（税込）の計算が一致しません", http.StatusBadRequest)
			return
		}
	}

	priceIncludingTax := incTax
	if priceIncludingTax == 0 {
		priceIncludingTax = int(float64(excTax) * (1.0 + tax))
	}

	pi := &Item{
		Month:             month,
		PurchasedOn:       &purchasedOn,
		Shop:              shop,
		Category:          category,
		ItemName:          itemName,
		Quantity:          uint16(quantity),
		PriceIncludingTax: strconv.Itoa(priceIncludingTax),
	}

	//レコード追加
	tx, _ := ab.db.Begin()
	defer func() {
		if recover() != nil {
			tx.Rollback()
		}
	}()
	if err := ab.addPurchasedItem(pi, tx); err != nil {
		tx.Rollback()
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tx.Commit()

	http.Redirect(w, r, "/", http.StatusFound)

}

//使用品 追加ボタン押下時のハンドラ
func (ab *accountBook) addUsedItemHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		code := http.StatusMethodNotAllowed
		http.Error(w, http.StatusText(code), code)
		return
	}

	//入力フォームの情報取得
	purchaseItemID, err := strconv.Atoi(r.FormValue("purchaseItemID"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	usedOn, err := time.Parse("2006-01-02", r.FormValue("usedOn"))
	timing := r.FormValue("timing")

	quantity, err := strconv.Atoi(r.FormValue("quantity"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	menu := r.FormValue("menu")

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ui := &usedItem{
		UsedOn:         &usedOn,
		Timing:         timing,
		PurchaseItemID: uint32(purchaseItemID),
		Quantity:       uint16(quantity),
		Menu:           menu,
	}

	tx, _ := ab.db.Begin()
	defer func() {
		if recover() != nil {
			tx.Rollback()
		}
	}()

	//レコード追加
	if err := ab.addUsedItem(ui, tx); err != nil {
		tx.Rollback()
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//ステータスの更新
	if err := ab.updateStatus(ui.PurchaseItemID, ui.Quantity, tx); err != nil {
		tx.Rollback()
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tx.Commit()
	http.Redirect(w, r, "/", http.StatusFound)

}
