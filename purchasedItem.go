package main

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

//Item 表示する購入品情報と使用品情報
type Item struct {
	ID                uint32
	Month             string
	PurchasedOn       *time.Time
	Shop              string
	Category          string
	ItemName          string
	Quantity          uint16
	PriceIncludingTax string
	ItemStatus        string
	IsDeleted         bool
	RemainingQuantity uint16
	UsedRate          string
	UsedItems         []usedItem
}

func (ab *accountBook) getPurchasedItems() ([]*Item, error) {
	PurchasedItems, err := ab.db.Query(getPurchasedItemsSQL)
	if err != nil {
		return nil, err
	}

	var is []*Item

	for PurchasedItems.Next() {
		var i Item
		if err := PurchasedItems.Scan(
			&i.ID, &i.Month,
			&i.PurchasedOn, &i.Shop,
			&i.Category, &i.ItemName, &i.Quantity,
			&i.PriceIncludingTax, &i.ItemStatus,
		); err != nil {
			return nil, err
		}

		uis, usedQuantity, err := ab.getUsedItems(i.ID)
		if err != nil {
			return nil, err
		}

		i.RemainingQuantity = i.Quantity - usedQuantity
		i.UsedRate = fmt.Sprintf("%.0f %%", (float64(usedQuantity)/float64(i.Quantity))*100.0)
		i.UsedItems = uis

		is = append(is, &i)
	}

	return is, nil

}

func (ab *accountBook) addPurchasedItem(i *Item, tx *sql.Tx) error {

	_, err := tx.Exec(addPurchasedItemSQL,
		i.Month,
		i.PurchasedOn,
		i.Shop,
		i.Category,
		i.ItemName,
		i.Quantity,
		i.PriceIncludingTax,
		001,
		0,
		00002,
		00002,
	)
	if err != nil {
		return err
	}

	return nil
}

func (ab *accountBook) updatePurchasedItem(id int, sqlSet map[string]string, tx *sql.Tx) error {

	var sql strings.Builder

	sql.WriteString("UPDATE t_purchased_items SET ")

	for k, v := range sqlSet {
		sql.WriteString(k + " = " + v + ", ")
	}
	sql.WriteString("updated_at = CURRENT_TIMESTAMP ")

	sql.WriteString("WHERE id = " + strconv.Itoa(id))

	_, err := tx.Exec(sql.String())
	if err != nil {
		return err
	}

	//個数の変更があった場合のみ、ステータスの更新を行う
	if s, ok := sqlSet["quantity"]; ok {
		inputQuantity, err := strconv.Atoi(s)
		if err != nil {
			return err
		}

		_, usedQuantity, err := ab.getUsedItems(uint32(id))
		if err != nil {
			return err
		}

		var itemStatus int
		if 0 < usedQuantity && usedQuantity < uint16(inputQuantity) {
			itemStatus = 2
		} else if usedQuantity >= uint16(inputQuantity) {
			itemStatus = 3
		} else {
			return nil
		}

		_, err = tx.Exec(updateStatusSQL,
			itemStatus,
			id,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func (ab *accountBook) updateStatus(id uint32, inputQuantity uint16, tx *sql.Tx) error {
	var is []*Item

	PurchasedItems, err := tx.Query(getPurchasedItemSQL, id)
	if err != nil {
		return err
	}

	for PurchasedItems.Next() {
		var i Item
		if err := PurchasedItems.Scan(
			&i.ID, &i.Month,
			&i.PurchasedOn, &i.Shop,
			&i.Category, &i.ItemName, &i.Quantity,
			&i.PriceIncludingTax, &i.ItemStatus,
		); err != nil {
			return err
		}

		is = append(is, &i)
	}

	_, usedQuantity, err := ab.getUsedItems(id)
	if err != nil {
		return err
	}

	var itemStatus int
	if 0 < usedQuantity+inputQuantity && usedQuantity+inputQuantity < is[0].Quantity {
		itemStatus = 2
	} else if usedQuantity+inputQuantity >= is[0].Quantity {
		itemStatus = 3
	}

	_, err = tx.Exec(updateStatusSQL,
		itemStatus,
		id,
	)
	if err != nil {
		return err
	}

	return nil
}

const getPurchasedItemsSQL = `
	SELECT
		tpi.id,
		tpi.month,
		tpi.purchased_on,
		tpi.shop,
		mc.name,
		tpi.item_name,
		tpi.quantity,
		FORMAT(tpi.price_including_tax,0) as price_including_tax,
		mis.status
	FROM t_purchased_items tpi
	LEFT JOIN m_categorys mc ON tpi.category_id = mc.id
	INNER JOIN m_item_status mis ON tpi.item_status_id = mis.id
	WHERE tpi.is_deleted <> 1
	ORDER BY tpi.month,tpi.purchased_on,tpi.shop
`

const addPurchasedItemSQL = `
	INSERT INTO t_purchased_items
		(month, purchased_on, shop, category_id, item_name, quantity, price_including_tax, item_status_id, is_deleted, created_at, created_user_id, updated_at, updated_user_id)
	VALUES
		(?,?,?,?,?,?,?,?,?,CURRENT_TIMESTAMP,?,CURRENT_TIMESTAMP,?)
`
const getPurchasedItemSQL = `
	SELECT
		tpi.id,
		tpi.month,
		tpi.purchased_on,
		tpi.shop,
		mc.name,
		tpi.item_name,
		tpi.quantity,
		FORMAT(tpi.price_including_tax,0) as price_including_tax,
		mis.status
	FROM t_purchased_items tpi
	LEFT JOIN m_categorys mc ON tpi.category_id = mc.id
	INNER JOIN m_item_status mis ON tpi.item_status_id = mis.id
	WHERE tpi.is_deleted <> 1
		AND tpi.id = ?
	ORDER BY tpi.month,tpi.purchased_on,tpi.shop
`

const updateStatusSQL = `
	update
		t_purchased_items
	SET
		item_status_id = ?,
		updated_at = CURRENT_TIMESTAMP,
		updated_user_id = 00003
	WHERE id = ?
`
