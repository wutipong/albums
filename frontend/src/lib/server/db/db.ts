import type { DB } from 'kysely-codegen' // this is the Database interface we defined earlier
import { Pool } from 'pg'
import { Kysely, PostgresDialect } from 'kysely'
import { DATABASE_URL } from '$env/static/private'

const dialect = new PostgresDialect({
  pool: new Pool({
    connectionString: DATABASE_URL,
    max: 10,
  })
})

// Database interface is passed to Kysely's constructor, and from now on, Kysely 
// knows your database structure.
// Dialect is passed to Kysely's constructor, and from now on, Kysely knows how 
// to communicate with your database.
export const db = new Kysely<DB>({
  dialect,
})