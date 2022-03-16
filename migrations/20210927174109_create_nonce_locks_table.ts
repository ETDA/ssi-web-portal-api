import * as Knex from "knex"

export async function up(knex: Knex): Promise<void> {
  return knex.schema.createTable("nonce_locks", (table) => {
    table.string("did_address").notNullable().primary()
    table.boolean("is_locked").notNullable()
  })
}

export async function down(knex: Knex): Promise<void> {
  return knex.schema.dropTableIfExists("nonce_locks")
}

