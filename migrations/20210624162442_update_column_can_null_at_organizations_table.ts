import * as Knex from 'knex'

export async function up (knex: Knex): Promise<void> {
  return knex.schema.alterTable('organizations', function (table) {
    table.string('did_address', 255).nullable().alter()
    table.string('encrypted_id', 255).nullable().alter()
  })
}

export async function down (knex: Knex): Promise<void> {
  return knex.schema.alterTable('organizations', function (table) {
    table.string('did_address', 255).notNullable().unique().alter()
    table.string('encrypted_id', 255).notNullable().unique().alter()
  })
}

