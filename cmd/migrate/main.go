package main

//Esse comando também é um executável, por isso usa package main.
//go run ./cmd/migrate.

import (
	"context"       //Usado para criar timeout da migration.
	"log"           //Usado para imprimir mensagens no terminal.
	"os"            //Usado para ler arquivo.
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

	for _, file := range files {
		log.Println("running migrating", file)
		//Percorre cada migration encontrada.

		content, err := os.ReadFile(file)
		if err != nil {
			log.Fatal("failed to read migration file:", err)
		}
		//Lê o conteúdo do arquivo .sql.

		statements := splitSQLStatements(string(content))
		//Transforma o conteúdo do arquivo em uma lista de comandos.
		//A função faz -> strings.Split(content, ";").
		//strings.Split(content, ";") -> separa por ;.
		//Ex -> CREATE TABLE tasks (...); = statement 1: CREATE TABLE tasks (...).

		for _, statement := range statements {
			if strings.TrimSpace(statement) == "" {
				continue
			}
			//Ignora vazio.

			if _, err := db.ExecContext(ctx, statement); err != nil {
				//Executa o SQL no banco.

				log.Fatal("failed to execute migration:", err)
			}
		}
	}

	log.Println("migration completed successfully")
}

func splitSQLStatements(content string) []string {
	return strings.Split(content, ";")
}
