import { Pool } from 'pg'
import { Kysely, PostgresDialect } from 'kysely'
import { env } from '$env/dynamic/private'
import type { DB } from '$lib/server/db_types'

const dialect = new PostgresDialect({
  pool: new Pool({
    connectionString: env.DATABASE_URL,
    max: 10,
  })
})

export const db = new Kysely<DB>({
  dialect,
})