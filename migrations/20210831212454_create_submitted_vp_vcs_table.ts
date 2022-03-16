import * as Knex from "knex"

export async function up(knex: Knex): Promise<void> {
  return knex.schema.createTable("submitted_vp_vcs", function (table) {
    table.uuid("id").primary()
    table
      .string("submitted_vp_id", 255)
      .notNullable()
      .references("id")
      .inTable("submitted_vps")
      .onDelete("CASCADE")
    table.string("cid", 255).notNullable()
    table.string("schema_name", 255).notNullable()
    table.string("schema_type", 255).notNullable()
    table.dateTime("issuance_date").notNullable()
    table.string("issuer", 255).notNullable()
    table.string("holder", 255).notNullable()
    table.string("jwt", 255).notNullable()
    table.string("status", 255).notNullable()
    table.boolean("verify")
    table.dateTime("created_at").notNullable()
    table.dateTime("updated_at").notNullable()
  })
}

export async function down(knex: Knex): Promise<void> {
  return knex.schema.dropTableIfExists("submitted_vp_vcs")
}
