import * as Knex from "knex"

export async function up(knex: Knex): Promise<void> {
  return knex.schema.alterTable("nonce_locks", (table) => {
    table.specificType("created_at", "DATETIME(6)").notNullable().alter()
    table.specificType("updated_at", "DATETIME(6)").notNullable().alter()
  })
}

export async function down(knex: Knex): Promise<void> {
  return knex.schema.alterTable("nonce_locks", (table) => {
    table.dateTime("created_at").notNullable()
    table.dateTime("updated_at").notNullable()
  })
}
