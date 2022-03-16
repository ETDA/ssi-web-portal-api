import * as Knex from 'knex'

export async function up (knex: Knex): Promise<void> {
  return knex.schema.alterTable('users', function (table) {
    table.dropForeign(['organization_id'])
    table.dropColumn('organization_id')
  })

}

export async function down (knex: Knex): Promise<void> {
  return knex.schema.alterTable('users', function (table) {
    table.string('organization_id', 255).notNullable().references('id').inTable('organizations').onDelete('CASCADE')
  })
}

