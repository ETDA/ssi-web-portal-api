import * as Knex from 'knex'

export async function up (knex: Knex): Promise<void> {
  return knex.schema.createTable('users', function (table) {
    table.uuid('id').primary()
    table.string('organization_id', 255).notNullable().references('id').inTable('organizations').onDelete('CASCADE')
    table.string('email', 255).notNullable().unique()
    table.string('password', 255)
    table.string('first_name', 255).notNullable()
    table.string('last_name', 255).notNullable()
    table.string('status', 255).notNullable()
    table.string('role', 255).notNullable()
    table.string('date_of_birth', 255).notNullable()
    table.dateTime('created_at').notNullable()
    table.dateTime('updated_at').notNullable()
    table.dateTime('deleted_at')
  })
}

export async function down (knex: Knex): Promise<void> {
  return knex.schema.dropTableIfExists('users')
}

