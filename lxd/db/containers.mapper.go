package db

// The code below was generated by lxd-generate - DO NOT EDIT!

import (
	"database/sql"
	"fmt"

	"github.com/lxc/lxd/lxd/db/cluster"
	"github.com/lxc/lxd/lxd/db/query"
	"github.com/lxc/lxd/shared/api"
	"github.com/pkg/errors"
)

var _ = api.ServerEnvironment{}

var containerObjects = cluster.RegisterStmt(`
SELECT containers.id, projects.name AS project, containers.name, nodes.name AS node, containers.type, containers.architecture, containers.ephemeral, containers.creation_date, containers.stateful, containers.last_use_date, coalesce(containers.description, '')
  FROM containers JOIN projects ON project_id = projects.id JOIN nodes ON node_id = nodes.id
  ORDER BY projects.id, containers.name
`)

var containerObjectsByType = cluster.RegisterStmt(`
SELECT containers.id, projects.name AS project, containers.name, nodes.name AS node, containers.type, containers.architecture, containers.ephemeral, containers.creation_date, containers.stateful, containers.last_use_date, coalesce(containers.description, '')
  FROM containers JOIN projects ON project_id = projects.id JOIN nodes ON node_id = nodes.id
  WHERE containers.type = ? ORDER BY projects.id, containers.name
`)

var containerObjectsByProjectAndType = cluster.RegisterStmt(`
SELECT containers.id, projects.name AS project, containers.name, nodes.name AS node, containers.type, containers.architecture, containers.ephemeral, containers.creation_date, containers.stateful, containers.last_use_date, coalesce(containers.description, '')
  FROM containers JOIN projects ON project_id = projects.id JOIN nodes ON node_id = nodes.id
  WHERE project = ? AND containers.type = ? ORDER BY projects.id, containers.name
`)

var containerObjectsByNodeAndType = cluster.RegisterStmt(`
SELECT containers.id, projects.name AS project, containers.name, nodes.name AS node, containers.type, containers.architecture, containers.ephemeral, containers.creation_date, containers.stateful, containers.last_use_date, coalesce(containers.description, '')
  FROM containers JOIN projects ON project_id = projects.id JOIN nodes ON node_id = nodes.id
  WHERE node = ? AND containers.type = ? ORDER BY projects.id, containers.name
`)

var containerObjectsByProjectAndNodeAndType = cluster.RegisterStmt(`
SELECT containers.id, projects.name AS project, containers.name, nodes.name AS node, containers.type, containers.architecture, containers.ephemeral, containers.creation_date, containers.stateful, containers.last_use_date, coalesce(containers.description, '')
  FROM containers JOIN projects ON project_id = projects.id JOIN nodes ON node_id = nodes.id
  WHERE project = ? AND node = ? AND containers.type = ? ORDER BY projects.id, containers.name
`)

var containerObjectsByProjectAndName = cluster.RegisterStmt(`
SELECT containers.id, projects.name AS project, containers.name, nodes.name AS node, containers.type, containers.architecture, containers.ephemeral, containers.creation_date, containers.stateful, containers.last_use_date, coalesce(containers.description, '')
  FROM containers JOIN projects ON project_id = projects.id JOIN nodes ON node_id = nodes.id
  WHERE project = ? AND containers.name = ? ORDER BY projects.id, containers.name
`)

var containerObjectsByProjectAndNameAndType = cluster.RegisterStmt(`
SELECT containers.id, projects.name AS project, containers.name, nodes.name AS node, containers.type, containers.architecture, containers.ephemeral, containers.creation_date, containers.stateful, containers.last_use_date, coalesce(containers.description, '')
  FROM containers JOIN projects ON project_id = projects.id JOIN nodes ON node_id = nodes.id
  WHERE project = ? AND containers.name = ? AND containers.type = ? ORDER BY projects.id, containers.name
`)

var containerProfilesRef = cluster.RegisterStmt(`
SELECT project, name, value FROM containers_profiles_ref ORDER BY project, name
`)

var containerProfilesRefByProject = cluster.RegisterStmt(`
SELECT project, name, value FROM containers_profiles_ref WHERE project = ? ORDER BY project, name
`)

var containerProfilesRefByNode = cluster.RegisterStmt(`
SELECT project, name, value FROM containers_profiles_ref WHERE node = ? ORDER BY project, name
`)

var containerProfilesRefByProjectAndNode = cluster.RegisterStmt(`
SELECT project, name, value FROM containers_profiles_ref WHERE project = ? AND node = ? ORDER BY project, name
`)

var containerProfilesRefByProjectAndName = cluster.RegisterStmt(`
SELECT project, name, value FROM containers_profiles_ref WHERE project = ? AND name = ? ORDER BY project, name
`)

var containerConfigRef = cluster.RegisterStmt(`
SELECT project, name, key, value FROM containers_config_ref ORDER BY project, name
`)

var containerConfigRefByProject = cluster.RegisterStmt(`
SELECT project, name, key, value FROM containers_config_ref WHERE project = ? ORDER BY project, name
`)

var containerConfigRefByNode = cluster.RegisterStmt(`
SELECT project, name, key, value FROM containers_config_ref WHERE node = ? ORDER BY project, name
`)

var containerConfigRefByProjectAndNode = cluster.RegisterStmt(`
SELECT project, name, key, value FROM containers_config_ref WHERE project = ? AND node = ? ORDER BY project, name
`)

var containerConfigRefByProjectAndName = cluster.RegisterStmt(`
SELECT project, name, key, value FROM containers_config_ref WHERE project = ? AND name = ? ORDER BY project, name
`)

var containerDevicesRef = cluster.RegisterStmt(`
SELECT project, name, device, type, key, value FROM containers_devices_ref ORDER BY project, name
`)

var containerDevicesRefByProject = cluster.RegisterStmt(`
SELECT project, name, device, type, key, value FROM containers_devices_ref WHERE project = ? ORDER BY project, name
`)

var containerDevicesRefByNode = cluster.RegisterStmt(`
SELECT project, name, device, type, key, value FROM containers_devices_ref WHERE node = ? ORDER BY project, name
`)

var containerDevicesRefByProjectAndNode = cluster.RegisterStmt(`
SELECT project, name, device, type, key, value FROM containers_devices_ref WHERE project = ? AND node = ? ORDER BY project, name
`)

var containerDevicesRefByProjectAndName = cluster.RegisterStmt(`
SELECT project, name, device, type, key, value FROM containers_devices_ref WHERE project = ? AND name = ? ORDER BY project, name
`)

var containerID = cluster.RegisterStmt(`
SELECT containers.id FROM containers JOIN projects ON project_id = projects.id JOIN nodes ON node_id = nodes.id
  WHERE projects.name = ? AND containers.name = ?
`)

var containerCreate = cluster.RegisterStmt(`
INSERT INTO containers (project_id, name, node_id, type, architecture, ephemeral, creation_date, stateful, last_use_date, description)
  VALUES ((SELECT id FROM projects WHERE name = ?), ?, (SELECT id FROM nodes WHERE name = ?), ?, ?, ?, ?, ?, ?, ?)
`)

var containerCreateConfigRef = cluster.RegisterStmt(`
INSERT INTO containers_config (container_id, key, value)
  VALUES (?, ?, ?)
`)

var containerCreateDevicesRef = cluster.RegisterStmt(`
INSERT INTO containers_devices (container_id, name, type)
  VALUES (?, ?, ?)
`)
var containerCreateDevicesConfigRef = cluster.RegisterStmt(`
INSERT INTO containers_devices_config (container_device_id, key, value)
  VALUES (?, ?, ?)
`)

var containerRename = cluster.RegisterStmt(`
UPDATE containers SET name = ? WHERE project_id = (SELECT id FROM projects WHERE name = ?) AND name = ?
`)

var containerDelete = cluster.RegisterStmt(`
DELETE FROM containers WHERE project_id = (SELECT id FROM projects WHERE name = ?) AND name = ?
`)

// ContainerList returns all available containers.
func (c *ClusterTx) ContainerList(filter ContainerFilter) ([]Container, error) {
	// Result slice.
	objects := make([]Container, 0)

	// Check which filter criteria are active.
	criteria := map[string]interface{}{}
	if filter.Project != "" {
		criteria["Project"] = filter.Project
	}
	if filter.Name != "" {
		criteria["Name"] = filter.Name
	}
	if filter.Node != "" {
		criteria["Node"] = filter.Node
	}
	if filter.Type != -1 {
		criteria["Type"] = filter.Type
	}

	// Pick the prepared statement and arguments to use based on active criteria.
	var stmt *sql.Stmt
	var args []interface{}

	if criteria["Project"] != nil && criteria["Name"] != nil && criteria["Type"] != nil {
		stmt = c.stmt(containerObjectsByProjectAndNameAndType)
		args = []interface{}{
			filter.Project,
			filter.Name,
			filter.Type,
		}
	} else if criteria["Project"] != nil && criteria["Node"] != nil && criteria["Type"] != nil {
		stmt = c.stmt(containerObjectsByProjectAndNodeAndType)
		args = []interface{}{
			filter.Project,
			filter.Node,
			filter.Type,
		}
	} else if criteria["Project"] != nil && criteria["Name"] != nil {
		stmt = c.stmt(containerObjectsByProjectAndName)
		args = []interface{}{
			filter.Project,
			filter.Name,
		}
	} else if criteria["Node"] != nil && criteria["Type"] != nil {
		stmt = c.stmt(containerObjectsByNodeAndType)
		args = []interface{}{
			filter.Node,
			filter.Type,
		}
	} else if criteria["Project"] != nil && criteria["Type"] != nil {
		stmt = c.stmt(containerObjectsByProjectAndType)
		args = []interface{}{
			filter.Project,
			filter.Type,
		}
	} else if criteria["Type"] != nil {
		stmt = c.stmt(containerObjectsByType)
		args = []interface{}{
			filter.Type,
		}
	} else {
		stmt = c.stmt(containerObjects)
		args = []interface{}{}
	}

	// Dest function for scanning a row.
	dest := func(i int) []interface{} {
		objects = append(objects, Container{})
		return []interface{}{
			&objects[i].ID,
			&objects[i].Project,
			&objects[i].Name,
			&objects[i].Node,
			&objects[i].Type,
			&objects[i].Architecture,
			&objects[i].Ephemeral,
			&objects[i].CreationDate,
			&objects[i].Stateful,
			&objects[i].LastUseDate,
			&objects[i].Description,
		}
	}

	// Select.
	err := query.SelectObjects(stmt, dest, args...)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to fetch containers")
	}

	// Fill field Config.
	configObjects, err := c.ContainerConfigRef(filter)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to fetch field Config")
	}

	for i := range objects {
		value := configObjects[objects[i].Name]
		if value == nil {
			value = map[string]string{}
		}
		objects[i].Config = value
	}

	// Fill field Devices.
	devicesObjects, err := c.ContainerDevicesRef(filter)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to fetch field Devices")
	}

	for i := range objects {
		value := devicesObjects[objects[i].Name]
		if value == nil {
			value = map[string]map[string]string{}
		}
		objects[i].Devices = value
	}

	// Fill field Profiles.
	profilesObjects, err := c.ContainerProfilesRef(filter)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to fetch field Profiles")
	}

	for i := range objects {
		value := profilesObjects[objects[i].Name]
		if value == nil {
			value = []string{}
		}
		objects[i].Profiles = value
	}

	return objects, nil
}

// ContainerGet returns the container with the given key.
func (c *ClusterTx) ContainerGet(project string, name string) (*Container, error) {
	filter := ContainerFilter{}
	filter.Project = project
	filter.Name = name
	filter.Type = -1

	objects, err := c.ContainerList(filter)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to fetch Container")
	}

	switch len(objects) {
	case 0:
		return nil, ErrNoSuchObject
	case 1:
		return &objects[0], nil
	default:
		return nil, fmt.Errorf("More than one container matches")
	}
}

// ContainerID return the ID of the container with the given key.
func (c *ClusterTx) ContainerID(project string, name string) (int64, error) {
	stmt := c.stmt(containerID)
	rows, err := stmt.Query(project, name)
	if err != nil {
		return -1, errors.Wrap(err, "Failed to get container ID")
	}
	defer rows.Close()

	// For sanity, make sure we read one and only one row.
	if !rows.Next() {
		return -1, ErrNoSuchObject
	}
	var id int64
	err = rows.Scan(&id)
	if err != nil {
		return -1, errors.Wrap(err, "Failed to scan ID")
	}
	if rows.Next() {
		return -1, fmt.Errorf("More than one row returned")
	}
	err = rows.Err()
	if err != nil {
		return -1, errors.Wrap(err, "Result set failure")
	}

	return id, nil
}

// ContainerExists checks if a container with the given key exists.
func (c *ClusterTx) ContainerExists(project string, name string) (bool, error) {
	_, err := c.ContainerID(project, name)
	if err != nil {
		if err == ErrNoSuchObject {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

// ContainerCreate adds a new container to the database.
func (c *ClusterTx) ContainerCreate(object Container) (int64, error) {
	// Check if a container with the same key exists.
	exists, err := c.ContainerExists(object.Project, object.Name)
	if err != nil {
		return -1, errors.Wrap(err, "Failed to check for duplicates")
	}
	if exists {
		return -1, fmt.Errorf("This container already exists")
	}

	args := make([]interface{}, 10)

	// Populate the statement arguments.
	args[0] = object.Project
	args[1] = object.Name
	args[2] = object.Node
	args[3] = object.Type
	args[4] = object.Architecture
	args[5] = object.Ephemeral
	args[6] = object.CreationDate
	args[7] = object.Stateful
	args[8] = object.LastUseDate
	args[9] = object.Description

	// Prepared statement to use.
	stmt := c.stmt(containerCreate)

	// Execute the statement.
	result, err := stmt.Exec(args...)
	if err != nil {
		return -1, errors.Wrap(err, "Failed to create container")
	}

	id, err := result.LastInsertId()
	if err != nil {
		return -1, errors.Wrap(err, "Failed to fetch container ID")
	}

	// Insert config reference.
	stmt = c.stmt(containerCreateConfigRef)
	for key, value := range object.Config {
		_, err := stmt.Exec(id, key, value)
		if err != nil {
			return -1, errors.Wrap(err, "Insert config for container")
		}
	}

	// Insert devices reference.
	for name, config := range object.Devices {
		typ, ok := config["type"]
		if !ok {
			return -1, fmt.Errorf("No type for device %s", name)
		}
		typCode, err := dbDeviceTypeToInt(typ)
		if err != nil {
			return -1, errors.Wrapf(err, "Device type code for %s", typ)
		}
		stmt = c.stmt(containerCreateDevicesRef)
		result, err := stmt.Exec(id, name, typCode)
		if err != nil {
			return -1, errors.Wrapf(err, "Insert device %s", name)
		}
		deviceID, err := result.LastInsertId()
		if err != nil {
			return -1, errors.Wrap(err, "Failed to fetch device ID")
		}
		stmt = c.stmt(containerCreateDevicesConfigRef)
		for key, value := range config {
			_, err := stmt.Exec(deviceID, key, value)
			if err != nil {
				return -1, errors.Wrap(err, "Insert config for container")
			}
		}
	}

	// Insert profiles reference.
	err = ContainerProfilesInsert(c.tx, int(id), object.Project, object.Profiles)
	if err != nil {
		return -1, errors.Wrap(err, "Insert profiles for container")
	}
	return id, nil
}

// ContainerProfilesRef returns entities used by containers.
func (c *ClusterTx) ContainerProfilesRef(filter ContainerFilter) (map[string][]string, error) {
	// Result slice.
	objects := make([]struct {
		Project string
		Name    string
		Value   string
	}, 0)

	// Check which filter criteria are active.
	criteria := map[string]interface{}{}
	if filter.Project != "" {
		criteria["Project"] = filter.Project
	}
	if filter.Name != "" {
		criteria["Name"] = filter.Name
	}

	// Pick the prepared statement and arguments to use based on active criteria.
	var stmt *sql.Stmt
	var args []interface{}

	if criteria["Project"] != nil && criteria["Node"] != nil {
		stmt = c.stmt(containerProfilesRefByProjectAndNode)
		args = []interface{}{
			filter.Project,
			filter.Node,
		}
	} else if criteria["Project"] != nil && criteria["Name"] != nil {
		stmt = c.stmt(containerProfilesRefByProjectAndName)
		args = []interface{}{
			filter.Project,
			filter.Name,
		}
	} else if criteria["Project"] != nil {
		stmt = c.stmt(containerProfilesRefByProject)
		args = []interface{}{
			filter.Project,
		}
	} else if criteria["Node"] != nil {
		stmt = c.stmt(containerProfilesRefByNode)
		args = []interface{}{
			filter.Node,
		}
	} else {
		stmt = c.stmt(containerProfilesRef)
		args = []interface{}{}
	}

	// Dest function for scanning a row.
	dest := func(i int) []interface{} {
		objects = append(objects, struct {
			Project string
			Name    string
			Value   string
		}{})
		return []interface{}{
			&objects[i].Project,
			&objects[i].Name,
			&objects[i].Value,
		}
	}

	// Select.
	err := query.SelectObjects(stmt, dest, args...)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to fetch string ref for containers")
	}

	// Build index by primary name.
	index := map[string][]string{}

	for _, object := range objects {
		item, ok := index[object.Name]
		if !ok {
			item = []string{}
		}

		index[object.Name] = append(item, object.Value)
	}

	return index, nil
}

// ContainerConfigRef returns entities used by containers.
func (c *ClusterTx) ContainerConfigRef(filter ContainerFilter) (map[string]map[string]string, error) {
	// Result slice.
	objects := make([]struct {
		Project string
		Name    string
		Key     string
		Value   string
	}, 0)

	// Check which filter criteria are active.
	criteria := map[string]interface{}{}
	if filter.Project != "" {
		criteria["Project"] = filter.Project
	}
	if filter.Name != "" {
		criteria["Name"] = filter.Name
	}

	// Pick the prepared statement and arguments to use based on active criteria.
	var stmt *sql.Stmt
	var args []interface{}

	if criteria["Project"] != nil && criteria["Name"] != nil {
		stmt = c.stmt(containerConfigRefByProjectAndName)
		args = []interface{}{
			filter.Project,
			filter.Name,
		}
	} else if criteria["Project"] != nil && criteria["Node"] != nil {
		stmt = c.stmt(containerConfigRefByProjectAndNode)
		args = []interface{}{
			filter.Project,
			filter.Node,
		}
	} else if criteria["Project"] != nil {
		stmt = c.stmt(containerConfigRefByProject)
		args = []interface{}{
			filter.Project,
		}
	} else if criteria["Node"] != nil {
		stmt = c.stmt(containerConfigRefByNode)
		args = []interface{}{
			filter.Node,
		}
	} else {
		stmt = c.stmt(containerConfigRef)
		args = []interface{}{}
	}

	// Dest function for scanning a row.
	dest := func(i int) []interface{} {
		objects = append(objects, struct {
			Project string
			Name    string
			Key     string
			Value   string
		}{})
		return []interface{}{
			&objects[i].Project,
			&objects[i].Name,
			&objects[i].Key,
			&objects[i].Value,
		}
	}

	// Select.
	err := query.SelectObjects(stmt, dest, args...)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to fetch  ref for containers")
	}

	// Build index by primary name.
	index := map[string]map[string]string{}

	for _, object := range objects {
		item, ok := index[object.Name]
		if !ok {
			item = map[string]string{}
		}

		index[object.Name] = item
		item[object.Key] = object.Value
	}

	return index, nil
}

// ContainerDevicesRef returns entities used by containers.
func (c *ClusterTx) ContainerDevicesRef(filter ContainerFilter) (map[string]map[string]map[string]string, error) {
	// Result slice.
	objects := make([]struct {
		Project string
		Name    string
		Device  string
		Type    int
		Key     string
		Value   string
	}, 0)

	// Check which filter criteria are active.
	criteria := map[string]interface{}{}
	if filter.Project != "" {
		criteria["Project"] = filter.Project
	}
	if filter.Name != "" {
		criteria["Name"] = filter.Name
	}

	// Pick the prepared statement and arguments to use based on active criteria.
	var stmt *sql.Stmt
	var args []interface{}

	if criteria["Project"] != nil && criteria["Name"] != nil {
		stmt = c.stmt(containerDevicesRefByProjectAndName)
		args = []interface{}{
			filter.Project,
			filter.Name,
		}
	} else if criteria["Project"] != nil && criteria["Node"] != nil {
		stmt = c.stmt(containerDevicesRefByProjectAndNode)
		args = []interface{}{
			filter.Project,
			filter.Node,
		}
	} else if criteria["Node"] != nil {
		stmt = c.stmt(containerDevicesRefByNode)
		args = []interface{}{
			filter.Node,
		}
	} else if criteria["Project"] != nil {
		stmt = c.stmt(containerDevicesRefByProject)
		args = []interface{}{
			filter.Project,
		}
	} else {
		stmt = c.stmt(containerDevicesRef)
		args = []interface{}{}
	}

	// Dest function for scanning a row.
	dest := func(i int) []interface{} {
		objects = append(objects, struct {
			Project string
			Name    string
			Device  string
			Type    int
			Key     string
			Value   string
		}{})
		return []interface{}{
			&objects[i].Project,
			&objects[i].Name,
			&objects[i].Device,
			&objects[i].Type,
			&objects[i].Key,
			&objects[i].Value,
		}
	}

	// Select.
	err := query.SelectObjects(stmt, dest, args...)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to fetch  ref for containers")
	}

	// Build index by primary name.
	index := map[string]map[string]map[string]string{}

	for _, object := range objects {
		item, ok := index[object.Name]
		if !ok {
			item = map[string]map[string]string{}
		}

		index[object.Name] = item
		config, ok := item[object.Device]
		if !ok {
			// First time we see this device, let's int the config
			// and add the type.
			deviceType, err := dbDeviceTypeToString(object.Type)
			if err != nil {
				return nil, errors.Wrapf(
					err, "unexpected device type code '%d'", object.Type)
			}
			config = map[string]string{}
			config["type"] = deviceType
			item[object.Device] = config
		}
		if object.Key != "" {
			config[object.Key] = object.Value
		}
	}

	return index, nil
}

// ContainerRename renames the container matching the given key parameters.
func (c *ClusterTx) ContainerRename(project string, name string, to string) error {
	stmt := c.stmt(containerRename)
	result, err := stmt.Exec(to, project, name)
	if err != nil {
		return errors.Wrap(err, "Rename container")
	}

	n, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "Fetch affected rows")
	}
	if n != 1 {
		return fmt.Errorf("Query affected %d rows instead of 1", n)
	}
	return nil
}

// ContainerDelete deletes the container matching the given key parameters.
func (c *ClusterTx) ContainerDelete(project string, name string) error {
	stmt := c.stmt(containerDelete)
	result, err := stmt.Exec(project, name)
	if err != nil {
		return errors.Wrap(err, "Delete container")
	}

	n, err := result.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "Fetch affected rows")
	}
	if n != 1 {
		return fmt.Errorf("Query deleted %d rows instead of 1", n)
	}

	return nil
}
