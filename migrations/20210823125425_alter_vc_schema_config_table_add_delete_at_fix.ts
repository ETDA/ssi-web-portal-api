import * as Knex from "knex";


export async function up(knex: Knex): Promise<void> {
  return knex.schema.alterTable("vc_schema_configs", (table) => {
    table.dateTime("deleted_at").nullable().defaultTo(null)
  })
}


export async function down(knex: Knex): Promise<void> {
  return knex.schema.alterTable("vc_schema_configs", (table) => {
    table.dropColumn("deleted_at")
  })
}
