import * as Knex from 'knex'

export async function up (knex: Knex): Promise<void> {
  return knex.schema.alterTable('users', function (table) {
    table.string('verify_token', 255)
  })
}

export async function down (knex: Knex): Promise<void> {
  return knex.schema.alterTable('users', function (table) {
    table.dropColumn('verify_token')
  })
}

