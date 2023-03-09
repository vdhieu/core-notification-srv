env "local" {
    // ..
    // url of rdatabase managed in this env
    url = "postgres://postgres:postgres@localhost:5433/postgres?sslmode=disable"

    // url of the dev database for atlas testing (dev database)
    // See: https://atlasgo.io/concepts/dev-database

    dev = "postgres://postgres:postgres@localhost:5433/atlasdev?sslmode=disablet"

//    schemas = ["transactions", "transaction_legs"]

    migration {
        // URL where the migration directory resides. Only filesystem directories
        // are currently supported but more options will be added in the future.
        dir = "file://migrations"
        // Format of the migration directory: atlas | flyway | liquibase | goose | golang-migrate
        format = atlas
    }
}