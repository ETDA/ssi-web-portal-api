import * as Knex from "knex";


export async function up(knex: Knex): Promise<void> {
  return knex.schema.createTable("vcs", function (table) {
    table.uuid("id").primary()
    table.string("cid", 255).notNullable()
    table.string("schema_type", 255).notNullable()
    table.dateTime("issuance_date").notNullable()
    table.string("issuer", 255).notNullable()
    table.string("holder", 255).notNullable()
    table.string("jwt", 255).notNullable()
    table.string("status", 255).notNullable()
  })
}


export async function down(knex: Knex): Promise<void> {
	return knex.schema.dropTableIfExists("vcs")
}

