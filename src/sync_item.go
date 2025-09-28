// Copyright (c) 2025 DBCsync
//
// DBCsync is licensed under the MIT License.
// See the LICENSE file for details.

package main

import (
	"fmt"
)

func syncItemTemplate(conns *DBConnections, cfg SyncJobConfig) error {
	minEntry := cfg.MinEntry

	query := `
        SELECT 
            entry AS id,
            class AS class,
            subclass AS subclass,
            SoundOverrideSubclass AS sound_override_subclass,
            Material AS material,
            displayid AS display_id,
            InventoryType AS inventory_type,
            sheath AS sheath
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
        REPLACE INTO Item 
            (id, class, subclass, sound_override_subclass, material, display_id, inventory_type, sheath)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?)
    `)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("prepare REPLACE: %w", err)
	}
	defer stmt.Close()
    
    synced := 0
	for rows.Next() {
		var id, class, subclass, sound_override_subclass, material, display_id, inventory_type, sheath int

		if err := rows.Scan(
			&id,
			&class,
			&subclass,
			&sound_override_subclass,
			&material,
			&display_id,
			&inventory_type,
			&sheath,
		); err != nil {
			tx.Rollback()
			return fmt.Errorf("scan row: %w", err)
		}

		if _, err := stmt.Exec(
			id,
			class,
			subclass,
			sound_override_subclass,
			material,
			display_id,
			inventory_type,
			sheath,
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
