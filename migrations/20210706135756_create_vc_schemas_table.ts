import * as Knex from "knex"

export async function up(knex: Knex): Promise<void> {
  return knex.schema.createTable("vc_schemas", function (table) {
    table.uuid("id").primary()
    table.uuid("organization_id").references("id").inTable("organizations")
    table.string("schema_name", 255).notNullable()
    table.string("schema_type", 255).notNullable().unique()
    table.jsonb("schema_body").notNullable()
    table.boolean("public").notNullable()
    table.string("version", 255).notNullable()
    table.string("created_by", 255)
    table.dateTime("created_at").notNullable()
    table.dateTime("updated_at").notNullable()
    table.dateTime("deleted_at")
  })
}

export async function down(knex: Knex): Promise<void> {
  return knex.schema.dropTableIfExists("vc_schemas")
}
