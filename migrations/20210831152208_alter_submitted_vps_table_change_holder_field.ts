import * as Knex from "knex"

export async function up(knex: Knex): Promise<void> {
  return knex.schema.alterTable("submitted_vps", function (table) {
    table.renameColumn("holder_did", "holder")
  })
}

export async function down(knex: Knex): Promise<void> {
  return knex.schema.alterTable("submitted_vps", function (table) {
    table.renameColumn("holder", "holder_did")
  })
}
