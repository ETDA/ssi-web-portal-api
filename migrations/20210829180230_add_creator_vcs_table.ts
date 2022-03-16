import * as Knex from "knex"

export async function up(knex: Knex): Promise<void> {
  return knex.schema.alterTable("vcs", (table) => {
    table
      .string("creator_id", 255)
      .notNullable()
      .references("id")
      .inTable("users")
      .onDelete("CASCADE")
  })
}

export async function down(knex: Knex): Promise<void> {
  return knex.schema.alterTable("vcs", (table) => {
    table.dropColumn("creator_id")
  })
}
