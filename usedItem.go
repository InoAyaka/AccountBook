package main

import (
	"database/sql"
	"time"
)

type usedItem struct {
	ID             uint32
	UsedOn         *time.Time
	Timing         string
	Morning        string
	Lunch          string
	Dinner         string
	Others         string
	PurchaseItemID uint32
	Quantity       uint16
	Menu           string
}

func (ab *accountBook) getUsedItems(purchaseItemID uint32) ([]usedItem, uint16, error) {
	//使用品を取得
	usedItems, err := ab.db.Query(getUsedItemsSQL, purchaseItemID)
	if err != nil {
		return nil, 0, err
	}

	var uis []usedItem
	var usedQuantity uint16

	for usedItems.Next() {
		var ui usedItem
		if err := usedItems.Scan(
			&ui.ID, &ui.UsedOn,
			&ui.Morning, &ui.Lunch, &ui.Dinner, &ui.Others,
			&ui.Quantity, &ui.Menu,
		); err != nil {
			return nil, 0, err
		}
		usedQuantity += ui.Quantity
		uis = append(uis, ui)
	}

	return uis, usedQuantity, nil
}

func (ab *accountBook) addUsedItem(ui *usedItem, tx *sql.Tx) error {
	_, err := tx.Exec(addUsedItemSQL,
		ui.UsedOn,
		ui.Timing,
		ui.PurchaseItemID,
		ui.Quantity,
		ui.Menu,
		0,
		00002,
		00002,
	)
	if err != nil {
		return err
	}

	return nil
}

const getUsedItemsSQL = `
	SELECT
		tu.id,
		tu.used_on,
		CASE WHEN tu.used_timing_id = 001 THEN '◯' ELSE '' END AS morning,
		CASE WHEN tu.used_timing_id = 002 THEN '◯' ELSE '' END AS lunch,
		CASE WHEN tu.used_timing_id = 003 THEN '◯' ELSE '' END AS dinner,
		CASE WHEN tu.used_timing_id = 004 THEN '◯' ELSE '' END AS others,
		tu.quantity,
		tu.menu
	FROM t_used tu
	WHERE tu.is_deleted <> 1
		AND tu.purchase_item_id = ?
	ORDER BY tu.used_on,tu.used_timing_id
`

const addUsedItemSQL = `
	INSERT INTO t_used
    	(used_on, used_timing_id, purchase_item_id, quantity, menu, is_deleted, created_at, created_user_id, updated_at, updated_user_id)
	VALUES
   		(?,?,?,?,?,?,CURRENT_TIMESTAMP,?,CURRENT_TIMESTAMP,?)
`
