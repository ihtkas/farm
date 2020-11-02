package store

// TODO: Evaluate otehr ways to make this externalise. Consider remodelling in updates to app in production, migration etc.,

const (
	productTable     = "product"
	userProductTable = "user_product"

	productIdColumn         = "product_id"
	userIdColumn            = "user_id"
	nameColumn              = "name"
	expiryColumn            = "expiry"
	quantityColumn          = "quantity"
	minQuantityColumn       = "min_quantity"
	pricePerQuantityColumn  = "price_per_quantity"
	pickupLocLatColumn      = "pickup_loc_lat"
	pickupLocLonColumn      = "pickup_loc_lon"
	descriptionColumn       = "description"
	tagsColumn              = "tags"
	pickUpLocColumn         = "pickup_loc"
	productProtoColumn      = "product_proto"
	insertionTimeUUIDColumn = "insertion_time"

	createProductCasandraQuery = "CREATE TABLE IF NOT EXISTS " + productTable + " (" +
		productIdColumn + " UUID PRIMARY KEY," +
		productProtoColumn + " blob)"

	createUserProductCasandraQuery = "CREATE TABLE IF NOT EXISTS " + userProductTable + " (" +
		userIdColumn + " UUID," +
		insertionTimeUUIDColumn + " UUID," +
		productProtoColumn + " blob," +
		"PRIMARY KEY ((" + userIdColumn + "), " + insertionTimeUUIDColumn + "))" +
		"WITH CLUSTERING ORDER BY (" + insertionTimeUUIDColumn + " DESC)"

	createProductPGQuery = "CREATE TABLE IF NOT EXISTS " + productTable + " (" +
		productIdColumn + " uuid DEFAULT uuid_generate_v4() PRIMARY KEY," +
		nameColumn + " varchar," +
		quantityColumn + " int," +
		tagsColumn + " text[]," +
		pickUpLocColumn + " geography)"

	createPickUpLocIndexPGQuery = "CREATE INDEX IF NOT EXISTS product_gindx ON " + productTable + " USING GIST (" + pickUpLocColumn + ")"

	insertProductPGQuery = "INSERT INTO product (" + nameColumn + ", " + quantityColumn + ", " + tagsColumn + ", " + pickUpLocColumn + ") values ($1, $2, $3, $4) returning (" + productIdColumn + ")"

	insertProductCassandraQuery = "INSERT INTO " + productTable + "(" +
		productIdColumn + "," +
		productProtoColumn + ") values (?, ?)"

	insertUserProductCassandraQuery = "INSERT INTO " + userProductTable + "(" +
		userIdColumn + "," +
		insertionTimeUUIDColumn + "," +
		productProtoColumn + ") values (?, ?, ?)"

	nearByProductsPGQuery = "SELECT " + productIdColumn + ", TRUNC(ST_Distance(" + pickUpLocColumn + ", ref_geoloc)) AS distance" +
		"	FROM " + productTable +
		" CROSS JOIN (" +
		"SELECT ST_MakePoint($1, $2)::geography AS ref_geoloc) AS r " +
		"WHERE ST_DWithin(" + pickUpLocColumn + ", ref_geoloc, $3)" +
		"ORDER BY ST_Distance(" + pickUpLocColumn + ", ref_geoloc) LIMIT $4 OFFSET $5"

	selectProductCassandraQuery = "SELECT " + productProtoColumn +

		" FROM " + productTable + " WHERE " + productIdColumn + "=?"
	selectProductByUserAfterTimeCassandraQuery = "SELECT " + productProtoColumn + ", " + insertionTimeUUIDColumn +
		" FROM " + userProductTable + " WHERE " + userIdColumn + "=? and " + insertionTimeUUIDColumn + ">? limit ?"
	selectProductByUserCassandraQuery = "SELECT " + productProtoColumn + ", " + insertionTimeUUIDColumn +
		" FROM " + userProductTable + " WHERE " + userIdColumn + "=? limit ?"
)
