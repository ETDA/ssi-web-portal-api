import * as Knex from "knex";


export async function up(knex: Knex): Promise<void> {
    return knex.schema.createTable('mobile_user_group_users', function (table) {
        table.uuid('id').primary()
        table.string('mobile_user_group_id', 255).notNullable().references('id').inTable('mobile_user_groups').onDelete('CASCADE')
        table.string('mobile_user_id', 255).notNullable().references('id').inTable('mobile_users').onDelete('CASCADE')
        table.dateTime('created_at').notNullable()
        table.dateTime('updated_at').notNullable()
    })
}


export async function down(knex: Knex): Promise<void> {
    return knex.schema.dropTableIfExists('mobile_user_group_users')
}

