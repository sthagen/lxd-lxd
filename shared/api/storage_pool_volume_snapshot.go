package api

import (
	"time"
)

// StorageVolumeSnapshotsPost represents the fields available for a new LXD storage volume snapshot
//
// swagger:model
//
// API extension: storage_api_volume_snapshots.
type StorageVolumeSnapshotsPost struct {
	// Snapshot name
	// Example: snap0
	Name string `json:"name" yaml:"name"`

	// Description of the storage volume snapshot
	// Example: My custom snapshot
	Description string `json:"description" yaml:"description"`

	// When the snapshot expires (gets auto-deleted)
	// Example: 2021-03-23T17:38:37.753398689-04:00
	//
	// API extension: custom_volume_snapshot_expiry
	ExpiresAt *time.Time `json:"expires_at" yaml:"expires_at"`
}

// StorageVolumeSnapshotPost represents the fields required to rename/move a LXD storage volume snapshot
//
// swagger:model
//
// API extension: storage_api_volume_snapshots.
type StorageVolumeSnapshotPost struct {
	// New snapshot name
	// Example: snap1
	Name string `json:"name" yaml:"name"`

	// Initiate volume snapshot migration
	// Example: false
	//
	// API extension: storage_api_remote_volume_snapshot_copy
	Migration bool `json:"migration" yaml:"migration"`

	// Migration target (for push mode)
	//
	// API extension: storage_api_remote_volume_snapshot_copy
	Target *StorageVolumePostTarget `json:"target" yaml:"target"`
}

// StorageVolumeSnapshot represents a LXD storage volume snapshot
//
// swagger:model
//
// API extension: storage_api_volume_snapshots.
type StorageVolumeSnapshot struct {
	WithEntitlements `yaml:",inline"`

	// Snapshot name
	// Example: snap0
	Name string `json:"name" yaml:"name"`

	// Description of the storage volume
	// Example: My custom volume
	Description string `json:"description" yaml:"description"`

	// The content type (filesystem or block)
	// Example: filesystem
	//
	// API extension: custom_block_volumes
	ContentType string `json:"content_type" yaml:"content_type"`

	// Volume snapshot creation timestamp
	// Example: 2021-03-23T20:00:00-04:00
	// API extension: storage_volumes_created_at
	CreatedAt time.Time `json:"created_at" yaml:"created_at"`

	// When the snapshot expires (gets auto-deleted)
	// Example: 2021-03-23T17:38:37.753398689-04:00
	//
	// API extension: custom_volume_snapshot_expiry
	ExpiresAt *time.Time `json:"expires_at" yaml:"expires_at"`

	// Storage volume configuration map (refer to doc/storage.md)
	// Example: {"zfs.remove_snapshots": "true", "size": "50GiB"}
	Config map[string]string `json:"config" yaml:"config"`
}

// StorageVolumeSnapshotPut represents the modifiable fields of a LXD storage volume
//
// swagger:model
//
// API extension: storage_api_volume_snapshots.
type StorageVolumeSnapshotPut struct {
	// Description of the storage volume
	// Example: My custom volume
	Description string `json:"description" yaml:"description"`

	// When the snapshot expires (gets auto-deleted)
	// Example: 2021-03-23T17:38:37.753398689-04:00
	//
	// API extension: custom_volume_snapshot_expiry
	ExpiresAt *time.Time `json:"expires_at" yaml:"expires_at"`
}

// Writable converts a full StorageVolumeSnapshot struct into a StorageVolumeSnapshotPut struct (filters read-only fields).
func (storageVolumeSnapshot *StorageVolumeSnapshot) Writable() StorageVolumeSnapshotPut {
	return StorageVolumeSnapshotPut{
		Description: storageVolumeSnapshot.Description,
		ExpiresAt:   storageVolumeSnapshot.ExpiresAt,
	}
}

// SetWritable sets applicable values from StorageVolumeSnapshotPut struct to StorageVolumeSnapshot struct.
func (storageVolumeSnapshot *StorageVolumeSnapshot) SetWritable(put StorageVolumeSnapshotPut) {
	storageVolumeSnapshot.Description = put.Description
	storageVolumeSnapshot.ExpiresAt = put.ExpiresAt
}
