import * as Knex from "knex";


export async function up(knex: Knex): Promise<void> {
    return knex.schema.createTable('requested_vps', function (table) {
        table.uuid('id').primary()
        table.string('creator_id', 255).notNullable().references('id').inTable('users').onDelete('CASCADE')
        table.string('name', 255).notNullable().unique()
        table.string('status', 255).notNullable()
        table.integer('schema_count').notNullable()
        table.dateTime('created_at').notNullable()
        table.dateTime('updated_at').notNullable()
    })
}


export async function down(knex: Knex): Promise<void> {
    return knex.schema.dropTableIfExists('requested_vps')
}

