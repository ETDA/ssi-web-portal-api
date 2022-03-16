import * as Knex from 'knex'

export async function up (knex: Knex): Promise<void> {
  return knex.schema.createTable('organization_api_keys', function (table) {
    table.uuid('id').primary()
    table.uuid('organization_id').notNullable().references('id').inTable('organizations').onDelete('CASCADE')
    table.string('name', 255)
    table.string('key', 255)
    table.boolean('read')
    table.boolean('write')
    table.dateTime('created_at').notNullable()
    table.dateTime('updated_at').notNullable()
    table.dateTime('deleted_at')
  })
}

export async function down (knex: Knex): Promise<void> {
  return knex.schema.dropTableIfExists('api_keys')
}
