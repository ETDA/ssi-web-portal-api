import * as Knex from "knex"

export async function up(knex: Knex): Promise<void> {
  return knex.schema.createTable("submitted_vps", (table) => {
    table.uuid("id").primary()
    table
      .string("requested_vp_id", 255)
      .notNullable()
      .references("id")
      .inTable("requested_vps")
      .onDelete("CASCADE")
    table.string("holder_did", 255).notNullable()
    table.text("jwt").notNullable()
    table.string("tags", 255).notNullable()
    table.integer("document_count").notNullable()
    table.dateTime("created_at").notNullable()
    table.dateTime("updated_at").notNullable()
  })
}

export async function down(knex: Knex): Promise<void> {
  return knex.schema.dropTableIfExists("submitted_vps")
}
