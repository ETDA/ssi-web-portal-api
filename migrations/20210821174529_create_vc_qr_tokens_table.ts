import * as Knex from "knex"

export async function up(knex: Knex): Promise<void> {
  return knex.schema.createTable("vc_qr_tokens", function (table) {
    table.uuid("id").primary()
    table.string("token", 255).notNullable()
    table.jsonb("cids").notNullable()
    table.dateTime("created_at").notNullable()
    table.dateTime("deleted_at")
  })
}

export async function down(knex: Knex): Promise<void> {
  return knex.schema.dropTableIfExists("vc_qr_tokens")
}
