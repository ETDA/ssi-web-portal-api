import * as Knex from "knex"

export async function up(knex: Knex): Promise<void> {
  return knex.schema.alterTable("vcs", function (table) {
    table.text("rejected_reason").defaultTo("").notNullable()
  })
}

export async function down(knex: Knex): Promise<void> {
  return knex.schema.alterTable("vcs", function (table) {
    table.dropColumn("rejected_reason")
  })
}
