package bucket

const (
	createQuery = `
		CREATE (b:Bucket {
			id: $id,
            version: 0,
			name: $name,
			description: $description,
			key: $key,
			ownerID: $ownerID,
			status: $status,
			type: $type,
			metadata: $metadata,
			createdAt: $createdAt,
			updatedAt: $updatedAt
		})
		RETURN b
	`

	getByIDQuery = `
		MATCH (b:Bucket {id: $id})
		OPTIONAL MATCH (b)-[:CONTAINS]->(a:Asset)
		RETURN b, collect(a) as assets
	`

	getBySlugQuery = `
		MATCH (b:Bucket {slug: $slug})
		OPTIONAL MATCH (b)-[:CONTAINS]->(a:Asset)
		RETURN b, collect(a) as assets
	`

	updateQuery = `
        MATCH (b:Bucket {id: $id})
        SET b.version = CASE WHEN $expectedVersion IS NULL THEN coalesce(b.version, 0) + 1 ELSE CASE WHEN b.version = $expectedVersion THEN b.version + 1 ELSE b.version END END,
            b.name = $name,
			b.description = $description,
			b.ownerID = $ownerID,
			b.status = $status,
			b.type = $type,
			b.metadata = $metadata,
			b.updatedAt = $updatedAt
		RETURN b
	`

	deleteQuery = `
		MATCH (b:Bucket {id: $id})
		OPTIONAL MATCH (b)-[r:CONTAINS]->(a:Asset)
		DELETE r, b
	`

	listQuery = `
		MATCH (b:Bucket)
		RETURN b
		ORDER BY b.createdAt DESC
		SKIP $offset
		LIMIT $limit
	`

	searchQuery = `
		MATCH (b:Bucket)
		WHERE b.name CONTAINS $query OR b.description CONTAINS $query
		RETURN b
		ORDER BY b.createdAt DESC
		SKIP $offset
		LIMIT $limit
	`

	getByOwnerIDQuery = `
		MATCH (b:Bucket {ownerID: $ownerID})
		RETURN b
		ORDER BY b.createdAt DESC
		LIMIT $limit
	`

	addAssetQuery = `
		MATCH (b:Bucket {id: $bucketID})
		MATCH (a:Asset {id: $assetID})
		MERGE (b)-[:CONTAINS]->(a)
		RETURN b, a
	`

	removeAssetQuery = `
		MATCH (b:Bucket {id: $bucketID})-[r:CONTAINS]->(a:Asset {id: $assetID})
		DELETE r
	`

	getAssetIDsQuery = `
		MATCH (b:Bucket {id: $bucketID})-[:CONTAINS]->(a:Asset)
		RETURN a.id
		ORDER BY a.createdAt DESC
		LIMIT $limit
	`

	getByKeyQuery = `
		MATCH (b:Bucket {key: $key})
		OPTIONAL MATCH (b)-[:CONTAINS]->(a:Asset)
		RETURN b, collect(a) as assets
	`
	hasAssetQuery = `
		MATCH (b:Bucket {id: $bucketID})-[:CONTAINS]->(a:Asset {id: $assetID})
		RETURN count(a) as count
	`

	assetCountQuery = `
		MATCH (b:Bucket {id: $bucketID})-[:CONTAINS]->(a:Asset)
		RETURN count(a) as count
	`

	findByTypeQuery = `
		MATCH (b:Bucket {type: $type})
		RETURN b
		ORDER BY b.createdAt DESC
		SKIP $offset
		LIMIT $limit
	`

	findByStatusQuery = `
		MATCH (b:Bucket {status: $status})
		RETURN b
		ORDER BY b.createdAt DESC
		SKIP $offset
		LIMIT $limit
	`

	countQuery = `
		MATCH (b:Bucket)
		RETURN count(b) as count
	`

	countByOwnerIDQuery = `
		MATCH (b:Bucket {ownerID: $ownerID})
		RETURN count(b) as count
	`

	countByTypeQuery = `
		MATCH (b:Bucket {type: $type})
		RETURN count(b) as count
	`

	existsQuery = `
		MATCH (b:Bucket {id: $id})
		RETURN count(b) as count
	`

	existsByKeyQuery = `
		MATCH (b:Bucket {key: $key})
		RETURN count(b) as count
	`
)
