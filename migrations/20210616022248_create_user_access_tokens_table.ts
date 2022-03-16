import * as Knex from 'knex'

export async function up (knex: Knex): Promise<void> {
  return knex.schema.createTable('user_access_tokens', function (table) {
    table.uuid('id').primary()
    table.string('user_id', 255).notNullable().references('id').inTable('users').onDelete('CASCADE')
    table.string('token', 255)
    table.json('info')
    table.dateTime('created_at').notNullable()
  })
}

export async function down (knex: Knex): Promise<void> {
  return knex.schema.dropTableIfExists('user_access_tokens')
}

