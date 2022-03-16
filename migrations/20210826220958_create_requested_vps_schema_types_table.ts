import * as Knex from "knex"

export async function up(knex: Knex): Promise<void> {
  return knex.schema.createTable(
    "requested_vps_schema_types",
    function (table) {
      table.uuid("id").primary()
      table
        .string("requested_vp_id", 255)
        .notNullable()
        .references("id")
        .inTable("requested_vps")
        .onDelete("CASCADE")
      table.string("schema_type", 255).notNullable()
      table.boolean("is_required").notNullable()
      table.string("Noted", 255)
      table.dateTime("created_at").notNullable()
      table.dateTime("updated_at").notNullable()
    }
  )
}

export async function down(knex: Knex): Promise<void> {
  return knex.schema.dropTableIfExists("users")
}
