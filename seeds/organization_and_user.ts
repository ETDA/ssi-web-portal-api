import * as Knex from 'knex'

export async function seed (knex: Knex): Promise<void> {
  const row = await knex.table('users').first()
  if (!row) {

    await knex('organizations').insert(
      {
        'id': '4b3c391b-9299-49a9-8c98-00c5a3d0f10e',
        'juristic_id': '8b6b07a5-f091-42a1-b6ba-5c5cf95fda87',
        'name': 'Finema',
        'did_address': null,
        'encrypted_id': null,
        'created_at': (new Date()).toISOString().slice(0, 19).replace('T', ' '),
        'updated_at': (new Date()).toISOString().slice(0, 19).replace('T', ' ')
      }
    )

    // password is 12345678
    await knex('users').insert(
      {
        'id': '95fc0734-0266-4389-9c08-d955e4b05df6',
        'organization_id': '4b3c391b-9299-49a9-8c98-00c5a3d0f10e',
        'email': 'admin@finema.co',
        'password': 'pbkdf2_sha256$180000$gFva54kUgfI=$qMAh+rjuWnNrPzW8Sq1zDKX79qnXtdzqiHVcz0+G8jk=',
        'first_name': 'ปกรณ์',
        'last_name': 'ลี้สกุล',
        'status': 'ACTIVE',
        'date_of_birth': '',
        'role': 'ADMIN',
        'created_at': (new Date()).toISOString().slice(0, 19).replace('T', ' '),
        'updated_at': (new Date()).toISOString().slice(0, 19).replace('T', ' ')
      }
    )
  }
}
