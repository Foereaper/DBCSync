package main

import (
	"fmt"
)

func syncItemTemplate(conns *DBConnections, cfg SyncJobConfig) error {
	minEntry := cfg.MinEntry

	query := `
        SELECT 
            entry AS itemID,
            class AS ItemClass,
            subclass AS ItemSubClass,
            SoundOverrideSubclass AS sound_override_subclassid,
            Material AS MaterialID,
            displayid AS ItemDisplayInfo,
            InventoryType AS InventorySlotID,
            sheath AS SheathID
        FROM item_template
        WHERE entry >= ?
    `

	rows, err := conns.World.Query(query, minEntry)
	if err != nil {
		return fmt.Errorf("query world.item_template: %w", err)
	}
	defer rows.Close()

	tx, err := conns.DBC.Begin()
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	stmt, err := tx.Prepare(`
        REPLACE INTO item 
            (itemID, ItemClass, ItemSubClass, sound_override_subclassid, MaterialID, ItemDisplayInfo, InventorySlotID, SheathID)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?)
    `)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("prepare REPLACE: %w", err)
	}
	defer stmt.Close()
    
    synced := 0
	for rows.Next() {
		var itemID, itemClass, itemSubClass, soundOverrideSubClassID, materialID, itemDisplayInfo, inventorySlotID, sheathID int

		if err := rows.Scan(
			&itemID,
			&itemClass,
			&itemSubClass,
			&soundOverrideSubClassID,
			&materialID,
			&itemDisplayInfo,
			&inventorySlotID,
			&sheathID,
		); err != nil {
			tx.Rollback()
			return fmt.Errorf("scan row: %w", err)
		}

		if _, err := stmt.Exec(
			itemID,
			itemClass,
			itemSubClass,
			soundOverrideSubClassID,
			materialID,
			itemDisplayInfo,
			inventorySlotID,
			sheathID,
		); err != nil {
			tx.Rollback()
			return fmt.Errorf("replace row: %w", err)
		}
        synced++
	}

	if err := rows.Err(); err != nil {
		tx.Rollback()
		return fmt.Errorf("iterate rows: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit: %w", err)
	}

	fmt.Printf("Synced item_template → item, %d items\n", synced)
	return nil
}
