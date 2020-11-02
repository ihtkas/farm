package store

// TODO: Evaluate otehr ways to make this externalise. Consider remodelling in updates to app in production, migration etc.,

const (
	idColumn               = "id"
	nameColumn             = "name"
	expiryColumn           = "expiry"
	quantityColumn         = "quantity"
	minQuantityColumn      = "min_quantity"
	pricePerQuantityColumn = "price_per_quantity"
	pickupLocLatColumn     = "pickup_loc_lat"
	pickupLocLonColumn     = "pickup_loc_lon"
	descriptionColumn      = "description"
	tagsColumn             = "tags"
	pickUpLocColumn        = "pickup_loc"
	productTable           = "product"
	productProtoColumn     = "product_proto"

	createProductCasandraQuery = "CREATE TABLE IF NOT EXISTS " + productTable + " (" +
		idColumn + " UUID PRIMARY KEY," +
		productProtoColumn + " blob)"

	createProductPGQuery = "CREATE TABLE IF NOT EXISTS " + productTable + " (" +
		idColumn + " uuid DEFAULT uuid_generate_v4() PRIMARY KEY," +
		nameColumn + " varchar," +
		quantityColumn + " int," +
		tagsColumn + " text[]," +
		pickUpLocColumn + " geography)"

	createPickUpLocLocIndexPGQuery = "CREATE INDEX IF NOT EXISTS product_gindx ON " + productTable + " USING GIST (" + pickUpLocColumn + ")"

	insertProductPGQuery = "INSERT INTO product (" + nameColumn + "+, " + quantityColumn + "+, " + tagsColumn + "+, " + pickUpLocColumn + ") values ($1, $2, $3, $4) returning (id)"

	insertProductCassandraQuery = "INSERT INTO product (" +
		idColumn + "," +
		productProtoColumn + ") values (?, ?)"

	nearByProductsPGQuery = "SELECT " + idColumn + ", TRUNC(ST_Distance(" + pickUpLocColumn + ", ref_geoloc)) AS distance" +
		"	FROM " + productTable +
		" CROSS JOIN (" +
		"SELECT ST_MakePoint($1, $2)::geography AS ref_geoloc) AS r " +
		"WHERE ST_DWithin(" + pickUpLocColumn + ", ref_geoloc, $3)" +
		"ORDER BY ST_Distance(" + pickUpLocColumn + ", ref_geoloc) LIMIT $4 OFFSET $5"

	selectProductCassandraQuery = "SELECT " + productProtoColumn +
		" FROM " + productTable + " WHERE " + idColumn + "=?"
)
