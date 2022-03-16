import * as Knex from "knex";

export async function up(knex: Knex): Promise<void> {
    return knex.schema.alterTable("requested_vps", (table) => {
        table.string("qr_code_id", 255).notNullable()
    })
}


export async function down(knex: Knex): Promise<void> {
    return knex.schema.alterTable("requested_vps", (table) => {
        table.dropColumn("qr_code_id")
    })
}
