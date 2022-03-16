import * as Knex from 'knex'

export async function up(knex: Knex): Promise<void> {
    return knex.schema.createTable('vc_schema_histories', function (table) {
        table.uuid('schema_id').notNullable()
        table.jsonb('schema_body').notNullable()
        table.string('version', 255).notNullable()
        table.dateTime('created_at').notNullable()
        table.dateTime('updated_at').notNullable()
        table.unique(['schema_id', 'version']);
    })
}

export async function down(knex: Knex): Promise<void> {
    return knex.schema.dropTableIfExists('vc_schema_histories')
}
