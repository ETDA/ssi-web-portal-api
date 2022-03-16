import * as Knex from 'knex'

export async function up(knex: Knex): Promise<void> {
    return knex.schema.createTable('organizations', function (table) {
        table.uuid('id').primary()
        table.string('juristic_id', 255).notNullable().unique()
        table.string('did_address', 255).notNullable().unique()
        table.string('encrypted_id', 255).notNullable().unique()
        table.dateTime('created_at').notNullable()
        table.dateTime('updated_at').notNullable()
        table.dateTime('deleted_at')
    })
}

export async function down(knex: Knex): Promise<void> {
    return knex.schema.dropTableIfExists('organizations')
}
