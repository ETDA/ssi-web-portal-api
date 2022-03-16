import * as Knex from "knex"

export async function up(knex: Knex): Promise<void> {
    return knex.schema.createTable("wallet_configs", function (table) {
        table.uuid("id").primary()
        table.string("endpoint", 255).notNullable()
        table.string("access_token", 255).notNullable()
        table.dateTime("created_at").notNullable()
        table.dateTime("updated_at").notNullable()
    })
}

export async function down(knex: Knex): Promise<void> {
    return knex.schema.dropTableIfExists("wallet_configs")
}
