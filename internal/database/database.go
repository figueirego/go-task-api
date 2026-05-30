package database

//Fica responsável por abrir a conexão com o banco.

import (
	"database/sql"
	//Pacote padrão do Go para trabalhar com bancos SQL.
	//Não é específico do MySQL.
	//Serve -> abrir conexão, configurar pool, executar queries, fazer ping, fechar conexão.

	"os"   //Serve para acessar coisas do sistema operacional.
	"time" //Serve para trabalhar com data e hora.
)

func Open() (*sql.DB, error) {
	//*sql.DB = conexão/pool do banco, error = erro, se falhar.

	dsn := getenv("DATABASE_DSN", "app:app@tcp(localhost:3307)/go_task_api?parseTime=true")
	//dsn = Data Source Name -> string de conexão com o banco.
	//app:app -> usuario:senha.
	//@tcp(localhost:3307) -> conect via TCP no localhost, porta 3307.
	//"/go_task_api" -> nome do banco de dados.
	//?parseTime=true -> Diz ao driver MySQL para converter campos de data/hora para time.Time do Go.

	db, err := sql.Open("mysql", dsn)
	//sql.Open("mysql", dsn) -> Cria uma conexão/pool para o banco.
	//"mysql" -> nome do driver. -> import = _"github.com/go-sql-driver/mysql".
	//dsn = string de conexão.

	if err != nil {
		return nil, err
	}

	//Pool de conexões -> é um conjunto de conexões reutilizáveis com o banco.
	//Sem pool, a aplicação poderia abrir uma nova conexão a cada request.

	db.SetMaxOpenConns(10)
	//Define no máximo 10 conexões abertas ao mesmo tempo.

	db.SetMaxIdleConns(5)
	//Define no máximo 5 conexões paradas, prontas para reutilizar.

	db.SetConnMaxLifetime(5 * time.Minute)
	//Define que uma conexão pode viver por até 5 minutos antes de ser reciclada.

	if err := db.Ping(); err != nil {
		//db.Ping() -> Testa se o banco responde.
		//Se o MySQL estiver desligado, senha errada ou banco inexistente, retorna erro.

		db.Close()
		//Se o ping falhar, fechamos o pool antes de retornar erro.

		return nil, err
	}

	return db, nil
	//retorna o banco, retorna nil como erro.

}

func getenv(key, fallback string) string {
	//Essa função lê uma variável de ambiente.
	//key string -> nome da variável -> exemplo = DATABASE_DSN.
	//fallback string -> valor padrão caso a variável não exista.

	value := os.Getenv(key)
	//Lê a variável de ambiente, se não existir, retorna string vazia.

	if value == "" {
		return fallback
		//Se a variável não foi definida, retorna o valor padrão.

	}

	return value
}
