package asset

func buildAssetSaveQuery() string {
	return `
	MERGE (a:Asset {id: $id})
	ON CREATE SET
		a.slug = $slug,
		a.title = $title,
		a.description = $description,
		a.type = $type,
		a.genre = $genre,
		a.genres = $genres,
		a.tags = $tags,
		a.createdAt = $createdAt,
		a.updatedAt = $updatedAt,
		a.ownerId = $ownerId,
		a.parentId = $parentId,
		a.videos = $videos,
		a.images = $images,
		a.credits = $credits,
		a.publishRule = $publishRule,
		a.metadata = $metadata
	ON MATCH SET
		a.slug = $slug,
		a.title = $title,
		a.description = $description,
		a.type = $type,
		a.genre = $genre,
		a.genres = $genres,
		a.tags = $tags,
		a.updatedAt = $updatedAt,
		a.ownerId = $ownerId,
		a.parentId = $parentId,
		a.videos = $videos,
		a.images = $images,
		a.credits = $credits,
		a.publishRule = $publishRule,
		a.metadata = $metadata
	`
}

func buildAssetFindByIDQuery() string {
	return `
	MATCH (a:Asset {id: $id})
	RETURN a
	`
}

func buildAssetFindBySlugQuery() string {
	return `
	MATCH (a:Asset {slug: $slug})
	RETURN a
	`
}

func buildAssetDeleteQuery() string {
	return `
	MATCH (a:Asset {id: $id})
	DETACH DELETE a
	`
}

func buildAssetListQuery() string {
	return `
	MATCH (a:Asset)
	RETURN a
	ORDER BY a.createdAt DESC
	LIMIT $limit
	`
}

func buildAssetSearchQuery() string {
	return `
	CALL db.index.fulltext.queryNodes("assetSearch", $query) YIELD node
	WITH node AS a
	RETURN a
	ORDER BY a.createdAt DESC
	LIMIT $limit
	`
}

func buildParentRelationshipQuery() string {
	return `
	MATCH (child:Asset {id: $childID})
	MATCH (parent:Asset {id: a.parentId})
	MERGE (child)-[:IS_CHILD_OF]->(parent)
	`
}

func buildAssetFindByOwnerIDQuery() string {
	return `
	MATCH (a:Asset {ownerId: $ownerId})
	RETURN a
	ORDER BY a.createdAt DESC
	LIMIT $limit
	`
}

func buildAssetFindByParentIDQuery() string {
	return `
	MATCH (a:Asset {parentId: $parentId})
	RETURN a
	ORDER BY a.createdAt DESC
	LIMIT $limit
	`
}

func buildAssetFindByTypeQuery() string {
	return `
	MATCH (a:Asset {type: $type})
	RETURN a
	ORDER BY a.createdAt DESC
	LIMIT $limit
	`
}

func buildAssetFindByGenreQuery() string {
	return `
	MATCH (a:Asset {genre: $genre})
	RETURN a
	ORDER BY a.createdAt DESC
	LIMIT $limit
	`
}

func buildAssetFindByTagQuery() string {
	return `
	MATCH (a:Asset)
	WHERE $tag IN a.tags
	RETURN a
	ORDER BY a.createdAt DESC
	LIMIT $limit
	`
}
