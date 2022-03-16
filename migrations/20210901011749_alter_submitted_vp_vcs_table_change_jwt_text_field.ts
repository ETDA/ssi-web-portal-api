import * as Knex from "knex"

export async function up(knex: Knex): Promise<void> {
  return knex.schema.alterTable("submitted_vp_vcs", (table) => {
    table.text("jwt").notNullable().alter()
  })
}

export async function down(knex: Knex): Promise<void> {
  return knex.schema.alterTable("submitted_vp_vcs", (table) => {
    table.string("jwt", 255).notNullable().alter()
  })
}
