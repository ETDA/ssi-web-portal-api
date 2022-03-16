import * as Knex from "knex";


export async function up(knex: Knex): Promise<void> {
    return knex.schema.alterTable("vc_qr_tokens", (table) => {
        table.string("did_address", 255).notNullable()
    })
}


export async function down(knex: Knex): Promise<void> {
    return knex.schema.alterTable("vc_qr_tokens", (table) => {
        table.dropColumn("did_address")
    })
}

