import * as Knex from "knex"

export async function up(knex: Knex): Promise<void> {
  return knex.schema.alterTable("submitted_vps", function (table) {
    table.boolean("verify")
  })
}

export async function down(knex: Knex): Promise<void> {
  return knex.schema.alterTable("submitted_vps", function (table) {
    table.dropColumn("verify")
  })
}
