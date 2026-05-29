package main

//todo arquivo Go começa com um package.
//Um package é um agrupamento de código.
//package main = este projeto gera um executável.
//Se fosse uma biblioteca, poderia ser package task ou package users.

import (
	//import serve para trazer código pronto de outros pacotes.

	"database/sql"
	//Pacote padrão do Go para trabalhar com bancos SQL.
	//Não é específico do MySQL.
	//Serve -> abrir conexão, configurar pool, executar queries, fazer ping, fechar conexão.

	"encoding/json" //Serve para ler e escrever JSON.
	"log"           //Serve para imprimir logs no terminal.
	"net/http"      //Esse é o pacote HTTP nativo do Go. Permite criar servidor web sem framework externo.
	"os"            //Serve para acessar coisas do sistema operacional.
	"time"          //Serve para trabalhar com data e hora.

	"github.com/jomatheusdev/go-task-api/internal/task"
	//Acessa o pacote interno e usa task.NewStore().

	_ "github.com/go-sql-driver/mysql"
	//"_" -> importe esse pacote apenas pelos efeitos colaterais dele.
)

func main() {
	//func main() = Função inicial do programa.

	db, err := openDatabase()
	//Chama nossa função para abrir conexão com o banco.
	//Retorna -> db = conexão/pool do banco, err = erro, se falhar.

	if err != nil {
		log.Fatal("failed to connect database", err)
	}
	//Se não conseguir conectar no banco, a API nem inicia.

	defer db.Close()
	//Fecha o banco quando o programa encerrar.
	//defer = execute isso no final da função.

	store := task.NewStore()
	//Cria nossa Store em memória (Lugar para guardar tarefas).

	taskHandler := task.NewHandler(store)
	//Cria o handler de tasks e entrega a Store para ele.
	//main cria a Store -> main entrega a Store para o Handler -> Handler usa a Store.

	mux := http.NewServeMux()
	//Cria um roteador HTTP.
	//mux = multiplexer -> roteador de rotas HTTP.

	mux.HandleFunc("GET /health", healthHandler(db))
	//mux.HandleFunc = Registra uma rota.
	//Formato -> mux.HandleFunc("MÉTODO /caminho", função).
	//"GET /healt" -> Quando alguém fizer GET /health, execute essa função.
	//GET = Método HTTP para buscar dados.
	// "/health" = Rota de saúde da aplicação.

	taskHandler.RegisterRoutes(mux)
	//Registra rotas de task.
	//GET /tasks, POST /tasks, GET /tasks/{id}, PATCH /tasks/{id}, DELETE /tasks/{id}.

	log.Println("API running on http://localhost:8080")
	//Mostra no terminal -> saber que o servidor iniciou.

	log.Fatal(http.ListenAndServe(":8080", mux))
	//Inicia o servidor HTTP.
	//":8080" -> escute na porta 8080.
	//mux -> roteador com as rotas: GET /health, GET /tasks, POST /tasks.
	//ListenAndServe = fica rodando enquanto o servidor está ativo.
	//Se der erro -> log.Fatal -> imprime o erro e encerra o programa.

}

func openDatabase() (*sql.DB, error) {
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

func healthHandler(db *sql.DB) http.HandlerFunc {
	//Recebe -> db *sql.DB.
	//Retorna -> http.HandlerFunc.
	//Cria um handler HTTP usando o banco.

	return func(w http.ResponseWriter, r *http.Request) {
		//porque retornar a função? -> porque mux.HandleFunc espera = func(w http.ResponseWriter, r *http.Request).
		//func(w http.ResponseWriter, r *http.Request) {} -> função anônima.
		//w http.ResponseWriter -> w = Usado para escrever a resposta HTTP.
		//w http.ResponseWriter -> com ele voce retorna -> status code, headers, JSON, texto, erro.
		//r *http.Request -> r = representa a requisição recebida.
		//r *http.Request -> Nele você encontra -> método HTTP, rota, body, headers, query params, contexto.
		//status int -> status code.
		//data any -> qualquer dado para virar JSON.

		status := http.StatusOK
		//Status inicial.

		response := map[string]any{
			//Objeto de resposta.

			"status":   "ok",
			"time":     time.Now(),
			"database": "ok",
		}

		if err := db.Ping(); err != nil {
			status = http.StatusServiceUnavailable
			//Mudamos o status HTTP -> 503 (serviço indisponível).

			response["status"] = "degraded"
			response["database"] = "unavailable"
			//Mudamos a resposta.

		}

		w.Header().Set("Content-Type", "application/json")
		//Define o header -> Content-Type:application/json.
		//Content-Type:application/json -> avisa que a resposta está em JSON.

		w.WriteHeader(status)
		//Define o status HTTP -> 200(OK), 201(Created), 400(Bad Request).

		if err := json.NewEncoder(w).Encode(response); err != nil {
			//json.NewEncoder(w).Encode(data) -> Transforma o dado Go em JSON e escreve na resposta.

			log.Println("failed to write health response:", err)
			//Se der erro ao gerar JSON, ele mostra no terminal.

		}
	}
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
