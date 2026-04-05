import postgres from 'postgres';
import { DATABASE_URL } from '$env/static/private';

const connectionString = DATABASE_URL || 'postgresql://user:password@localhost:5432/db';

let sql: postgres.Sql;

export function initDb() {
  sql = postgres(connectionString);
}

export function cleanupDb() {
  if (sql) {
    sql.end();
  }
}

export function getDb() {
  if (!sql) {
    throw new Error('Database not initialized. Call initDb() first.');
  }
  return sql;
}   

// Initialize on module load
initDb();

// Cleanup on exit
process.on('exit', cleanupDb);
process.on('SIGINT', () => {
  cleanupDb();
  process.exit();
});
process.on('SIGTERM', () => {
  cleanupDb();
  process.exit();
});