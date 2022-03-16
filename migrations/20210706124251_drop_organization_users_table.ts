import * as Knex from 'knex'

export async function up (knex: Knex): Promise<void> {
  return knex.schema.dropTable('organization_users')
    .then(() => knex.schema.alterTable('users', function (table) {
      table.string('organization_id', 255).notNullable().references('id').inTable('organizations').onDelete('CASCADE')
    }))
}

export async function down (knex: Knex): Promise<void> {
  return knex.schema.createTable('organization_users', function (table) {
    table.uuid('id').primary()
    table.string('user_id', 255).notNullable().references('id').inTable('users').onDelete('CASCADE')
    table.string('organization_id', 255).notNullable().references('id').inTable('organizations').onDelete('CASCADE')
    table.string('role', 255).notNullable()
    table.dateTime('created_at').notNullable()
    table.dateTime('updated_at').notNullable()
  }).then(() => knex.schema.alterTable('users', function (table) {
    table.dropForeign(['organization_id'])
    table.dropColumn('organization_id')
  }))
}

