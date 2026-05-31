package main

//Esse comando também é um executável, por isso usa package main.
//go run ./cmd/migrate.

import (
	"context"      //Usado para criar timeout da migration.
	"database/sql" //Várias funções recebem -> db *sql.DB.
	"log"          //Usado para imprimir mensagens no terminal.
	"os"
	"path/filepath" //Usado para buscar arquivos.
	"sort"          //Usado para ordenar migrations.
	"strings"
	"time" //Usado para definir timeout.

	"github.com/jomatheusdev/go-task-api/internal/database"
	//Usado para abrir conexão com o banco.

	_ "github.com/go-sql-driver/mysql"
	//Necessário para o sql.Open("mysql", dsn) funcionar dentro do database.Open().
)

func main() {
	db, err := database.Open()
	if err != nil {
		log.Fatal("failed to connect database:", err)
	}
	defer db.Close()
	//Mesma lógica da API. Se não conseguir conectar, para o comando.

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	//context.Background() -> Cria um contexto base.
	//Como isso não é request HTTP, usamos Background.
	//context.withTimeout -> Cria um contexto que expira depois de 30 segundos.
	//Se alguma migration travar, ela não fica infinita.

	defer cancel()
	//Libera recursos do contexto ao final.

	if err := ensureSchemaMigrationsTable(ctx, db); err != nil {
		//Antes de saber quais migrations já rodaram, precisamos ter a tabela que guarda esse histórico.
		//Essa função vai criar a tabela se ela não existir.

		log.Fatal("failed to ensure schema_migrations table:", err)
	}

	appliedMigrations, err := loadAppliedMigrations(ctx, db)
	//Essa função busca no banco tudo que já foi aplicado.

	if err != nil {
		log.Fatal("failed to load applied migrations:", err)
	}

	files, err := filepath.Glob("migrations/*.sql")
	//Isso busca todos os arquivos que batem com o padrão: "migrations/*.sql".
	//Ex -> migrations/001_create_tasks_table.sql

	if err != nil {
		log.Fatal("failed to list migration files:", err)
	}

	sort.Strings(files)
	//Ordena os arquivos por nome. -> 001 vem antes de 002.

	if len(files) == 0 {
		log.Println("no migration files found")
		return
	}
	//Se não houver arquivos, ele avisa e sai.

	pendingCount := 0
	//Essa variável conta quantas migrations serão realmente aplicadas.
	//Se todas já foram aplicadas, ela continua 0.

	for _, file := range files {
		version := filepath.Base(file)
		//Se file for -> migrations/001_create_tasks_table.sql.
		//Retorna -> 001_create_tasks_table.sql. (nome que será salvo no banco).

		if appliedMigrations[version] {
			//Se o map contém essa version, pulamos.

			log.Println("skipping already applied migration:", version)
			continue
			//significa -> não execute o resto deste loop, vá para o próximo arquivo.

		}

		pendingCount++
		//Aumenta o contador de migrations pendentes.

		log.Println("running migration:", version)

		if err := runMigrationFile(ctx, db, file, version); err != nil {
			//Chama a função que lê o arquivo, executa os comandos e registra a migration.

			log.Println("failed to run migration:", err)
		}

		log.Println("applied migration:", version)
	}

	if pendingCount == 0 {
		log.Println("no pending migrations")
		return
	}

	log.Println("migration completed successfully")
}

func ensureSchemaMigrationsTable(ctx context.Context, db *sql.DB) error {
	//Executa essa SQL.

	_, err := db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version VARCHAR(255) PRIMARY KEY,
			applied_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		);
	`)

	return err
}

func loadAppliedMigrations(ctx context.Context, db *sql.DB) (map[string]bool, error) {
	rows, err := db.QueryContext(ctx, `
		SELECT version
		FROM schema_migrations
		ORDER BY version
	`)
	//QueryContext -> Podemos receber várias linhas.

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	appliedMigrations := map[string]bool{}
	//Cria um map vazio.
	//Esse map será usado para consulta rápida.
	//Ex -> if appliedMigrations["001_create_tasks_table.sql"] {}

	for rows.Next() {
		//Para cada linha do banco -> pega version e marca appliedMigrations[version] = true.

		var version string

		if err := rows.Scan(&version); err != nil {
			return nil, err
		}

		appliedMigrations[version] = true
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return appliedMigrations, nil
}

func runMigrationFile(ctx context.Context, db *sql.DB, file string, version string) error {
	//Essa função aplica uma migration.
	//1. lê arquivo
	//2. separa comandos SQL
	//3. abre transação
	//4. executa comandos
	//5. registra no schema_migrations
	//6. confirma transação

	content, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	statements := splitSQLStatements(string(content))

	tx, err := db.BeginTx(ctx, nil)
	//BeginTx -> começa uma transação.
	//tx -> representa a transação aberta.

	if err != nil {
		return err
	}

	defer tx.Rollback()
	//Garante que, se a função sair antes do Commit, a transação será desfeita.

	for _, statement := range statements {
		//Executar statements.

		if strings.TrimSpace(statement) == "" {
			continue
		}

		if _, err := tx.ExecContext(ctx, statement); err != nil {
			return err
		}
	}

	if _, err := tx.ExecContext(ctx, `
		INSERT INTO schema_migrations (version)
		VALUES (?)
	`, version); err != nil {
		return err
		//Registrar migration.

	}

	return tx.Commit()
	//Commit da transação.
	//Confirma a transação.
	//Se der certo: SQL aplicado, schema_migrations atualizado.
	//Se der erro antes do commit: rollback desfaz

}

func splitSQLStatements(content string) []string {
	return strings.Split(content, ";")
}
