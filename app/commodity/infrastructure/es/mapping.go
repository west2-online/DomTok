package es

const mapping = `{
	"mappings": {
		"properties": {
			"id": {"type": "keyword"},
			"name": {"type": "keyword"},
			"creator_id": {"type": "keyword"},
			"description": {"type": "text"},
			"category_id": {"type": "keyword"},
			"price": {"type": "float"},
			"for_sale": {"type": "keyword"},
			"shipping": {"type": "float"},
			"created_at": {"type": "date"},
			"updated_at": {"type": "date"},
			"deleted_at": {"type": "date"}}
}}`
