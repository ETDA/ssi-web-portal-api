import * as Knex from 'knex'

export async function up (knex: Knex): Promise<void> {
  return knex.schema.alterTable('vcs', function (table) {
    table.string('schema_name', 255).defaultTo('').notNullable()
  })
}

export async function down (knex: Knex): Promise<void> {
  return knex.schema.alterTable('vcs', function (table) {
    table.dropColumn('schema_name')
  })
}
