import * as Knex from "knex"

export async function up(knex: Knex): Promise<void> {
  return knex.schema.alterTable("nonce_locks", function (table) {
    table.string("vc_id", 255).defaultTo("").notNullable()
  })
}

export async function down(knex: Knex): Promise<void> {
  return knex.schema.alterTable("nonce_locks", function (table) {
    table.dropColumn("vc_id")
  })
}
