import * as Knex from "knex";


export async function up(knex: Knex): Promise<void> {
  return knex.schema.alterTable('organizations', function (table) {
    table.string('name',255).notNullable()
  })

}


export async function down(knex: Knex): Promise<void> {
  return knex.schema.alterTable('organizations', function (table) {
    table.dropColumn('name')
  })
}

