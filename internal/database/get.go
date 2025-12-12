package database

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
)

const (
	vmColumns = `id, name, namespace, cpu, ram, mac_address, network_id, image, disk_size, disk_path, created_at, labels`
	// networkColumns must match the actual table schema order
	networkColumns = `id, name, namespace, labels, mac_address, bridge, mode, net_address, netmask, dhcp, autostart, created_at`
)

// GetRecords retrieves all documents of type T from the specified collection
// that match the given namespace.
func GetVMRecords(
	ctx context.Context,
	db *sql.DB,
	namespace string,
) ([]VirtualMachine, error) {
	query := fmt.Sprintf("SELECT %s FROM %s", vmColumns, vmsTable)
	args := []any{}
	if namespace != "" {
		query += " WHERE namespace = ?"
		args = append(args, namespace)
	}
	var (
		vms       []VirtualMachine
		rawLabels string
	)
	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var vm VirtualMachine
		if err := rows.Scan(
			&vm.ID,
			&vm.Name,
			&vm.Namespace,
			&vm.CPU,
			&vm.RAM,
			&vm.MacAddress,
			&vm.NetworkID,
			&vm.Image,
			&vm.DiskSize,
			&vm.DiskPath,
			&vm.CreatedAt,
			&rawLabels,
		); err != nil {
			return nil, fmt.Errorf("scan failed: %w", err)
		}
		// parse JSON labels
		if err := json.Unmarshal([]byte(rawLabels), &vm.Labels); err != nil {
			return nil, fmt.Errorf("invalid labels JSON: %w", err)
		}
		vms = append(vms, vm)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}
	return vms, nil
}

// GetRecords retrieves all documents of type T from the specified collection
// that match the given namespace.
// GetVMByName fetches a single VirtualMachine by name.
// If namespace is non-empty, it will be included in the WHERE clause.
func GetVMByName(
	ctx context.Context,
	db *sql.DB,
	name, namespace string,
) (VirtualMachine, error) {
	// Build the base query
	query := fmt.Sprintf("SELECT %s FROM %s WHERE name = ?", vmColumns, vmsTable)
	args := []any{name}

	// Optionally filter by namespace
	if namespace != "" {
		query += " AND namespace = ?"
		args = append(args, namespace)
	}

	// Prepare a holder for the single record
	var (
		vm        VirtualMachine
		rawLabels string
	)

	// Execute the query
	row := db.QueryRowContext(ctx, query, args...)
	if err := row.Scan(
		&vm.ID,
		&vm.Name,
		&vm.Namespace,
		&vm.CPU,
		&vm.RAM,
		&vm.MacAddress,
		&vm.NetworkID,
		&vm.Image,
		&vm.DiskSize,
		&vm.DiskPath,
		&vm.CreatedAt,
		&rawLabels,
	); err != nil {
		if err == sql.ErrNoRows {
			return vm, fmt.Errorf("no VM found with name %q", name)
		}
		return vm, fmt.Errorf("failed to scan VM row: %w", err)
	}

	// Parse JSON labels
	if err := json.Unmarshal([]byte(rawLabels), &vm.Labels); err != nil {
		return vm, fmt.Errorf("invalid labels JSON for VM %q: %w", name, err)
	}

	return vm, nil
}

func GetNetworks(
	ctx context.Context,
	db *sql.DB,
	namespace string,
) ([]VirtualNetwork, error) {
	query := fmt.Sprintf("SELECT %s FROM %s", networkColumns, networksTable)
	args := []any{}
	if namespace != "" {
		query += " WHERE namespace = ?"
		args = append(args, namespace)
	}
	var (
		networks []VirtualNetwork
	)

	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var network VirtualNetwork
		var rawLabels, rawDHCP string
		if err := rows.Scan(
			&network.ID,
			&network.Name,
			&network.Namespace,
			&rawLabels,
			&network.MacAddress,
			&network.Bridge,
			&network.Mode,
			&network.NetAddress,
			&network.Netmask,
			&rawDHCP,
			&network.Autostart,
			&network.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan failed: %w", err)
		}
		// Parse JSON labels if present
		if rawLabels != "" {
			if err := json.Unmarshal([]byte(rawLabels), &network.Labels); err != nil {
				return nil, fmt.Errorf("invalid labels JSON: %w", err)
			}
		}
		// Parse JSON DHCP if present
		if rawDHCP != "" {
			if err := json.Unmarshal([]byte(rawDHCP), &network.DHCP); err != nil {
				return nil, fmt.Errorf("invalid DHCP JSON: %w", err)
			}
		}
		networks = append(networks, network)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}
	return networks, nil
}

const storeColumns = `id, name, namespace, labels, backend, artifacts_path, images_path, created_at`

// GetStores retrieves all store records from the database,
// optionally filtered by namespace.
func GetStores(
	ctx context.Context,
	db *sql.DB,
	namespace string,
) ([]Store, error) {
	query := fmt.Sprintf("SELECT %s FROM %s", storeColumns, storesTable)
	args := []any{}
	if namespace != "" {
		query += " WHERE namespace = ?"
		args = append(args, namespace)
	}

	var stores []Store
	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var (
			store     Store
			rawLabels string
		)
		if err := rows.Scan(
			&store.ID,
			&store.Name,
			&store.Namespace,
			&rawLabels,
			&store.Backend,
			&store.ArtifactsPath,
			&store.ImagesPath,
			&store.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan failed: %w", err)
		}
		// Parse JSON labels
		if rawLabels != "" {
			if err := json.Unmarshal([]byte(rawLabels), &store.Labels); err != nil {
				return nil, fmt.Errorf("invalid labels JSON: %w", err)
			}
		}

		// Fetch images for this store
		// TODO: optimize this to avoid N+1 query
		imgQuery := `SELECT id, name, version, os_profile, file, checksum, size, created_at FROM ` + imagesTable + ` WHERE store_id = ?`
		imgRows, err := db.QueryContext(ctx, imgQuery, store.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch images for store %s: %w", store.Name, err)
		}

		var images []Image
		for imgRows.Next() {
			var img Image
			img.StoreID = int64(store.ID)
			if err := imgRows.Scan(
				&img.ID, &img.Name, &img.Version, &img.OsProfile,
				&img.File, &img.Checksum, &img.Size, &img.CreatedAt,
			); err != nil {
				imgRows.Close()
				return nil, fmt.Errorf("scan image failed: %w", err)
			}
			images = append(images, img)
		}
		imgRows.Close()
		store.Images = images

		stores = append(stores, store)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}
	return stores, nil
}
