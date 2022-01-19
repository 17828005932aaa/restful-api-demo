package impl

const (
	insertResourceSQL = `INSERT INTO resource (
		id,vendor,region,zone,create_at,expire_at,category,type,instance_id,
		name,description,status,update_at,sync_at,sync_accout,public_ip,
		private_ip,pay_type,resource_hash,describe_hash
	) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?);`

	insertDescribeSQL = `INSERT INTO host (
		resource_id,cpu,memory,gpu_amount,gpu_spec,os_type,os_name,
		serial_number,image_id,internet_max_bandwidth_out,
		internet_max_bandwidth_in,key_pair_name,security_groups
	) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?);`

	queryHostSQL = `SELECT * FROM resource as r LEFT JOIN host h ON r.id=h.resource_id`

	updateResourceSQL = `UPDATE resource SET vendor=?,region=?,zone=?,expire_at=?,name=?,description=?,update_at=? WHERE id = ?`

	updateHostSQL = `UPDATE host SET cpu=?,memory=? WHERE resource_id = ?`

	deleteResourceSQL = `DELETE FROM resource WHERE id = ?`

	deleleHostSQL = `DELETE FROM host WHERE resource_id = ?`
)