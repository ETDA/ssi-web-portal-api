import * as Knex from "knex"

export async function up(knex: Knex): Promise<void> {
  return knex.schema.alterTable("nonce_locks", (table) => {
    table.renameColumn("did_address", "id")
    table.renameColumn("is_locked", "is_done")
    table.dateTime("created_at").notNullable()
    table.dateTime("updated_at").notNullable()
  })
}

export async function down(knex: Knex): Promise<void> {
  return knex.schema.alterTable("nonce_locks", (table) => {
    table.renameColumn("id", "did_address")
    table.renameColumn("is_done", "is_locked")
    table.dropColumn("created_at")
    table.dropColumn("updated_at")
  })
}
