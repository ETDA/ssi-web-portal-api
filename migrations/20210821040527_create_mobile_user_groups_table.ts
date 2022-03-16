import * as Knex from "knex";


export async function up(knex: Knex): Promise<void> {
    return knex.schema.createTable('mobile_user_groups', function (table) {
        table.uuid('id').primary()
        table.string('name', 255).notNullable().unique()
        table.dateTime('created_at').notNullable()
        table.dateTime('updated_at').notNullable()
    })
}


export async function down(knex: Knex): Promise<void> {

    return knex.schema.dropTableIfExists('mobile_user_groups')
}

