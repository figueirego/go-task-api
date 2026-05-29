package main

//todo arquivo Go começa com um package.
//Um package é um agrupamento de código.
//package main = este projeto gera um executável.
//Se fosse uma biblioteca, poderia ser package task ou package users.

import (
	//import serve para trazer código pronto de outros pacotes.

	"encoding/json" //Serve para ler e escrever JSON.
	"log"           //Serve para imprimir logs no terminal.
	"net/http"      //Esse é o pacote HTTP nativo do Go. Permite criar servidor web sem framework externo.
	"time"          //Serve para trabalhar com data e hora.

	"github.com/jomatheusdev/go-task-api/internal/task"
	//Acessa o pacote interno e usa task.NewStore().
)

func main() {
	//func main() = Função inicial do programa.

	store := task.NewStore()
	//Cria nossa Store em memória (Lugar para guardar tarefas).

	taskHandler := task.NewHandler(store)
	//Cria o handler de tasks e entrega a Store para ele.
	//main cria a Store -> main entrega a Store para o Handler -> Handler usa a Store.

	mux := http.NewServeMux()
	//Cria um roteador HTTP.
	//mux = multiplexer -> roteador de rotas HTTP.

	mux.HandleFunc("GET /health", healthHandler)
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

func healthHandler(w http.ResponseWriter, r *http.Request) {
	//func(w http.ResponseWriter, r *http.Request) {} -> função anônima.
	//w http.ResponseWriter -> w = Usado para escrever a resposta HTTP.
	//w http.ResponseWriter -> com ele voce retorna -> status code, headers, JSON, texto, erro.
	//r *http.Request -> r = representa a requisição recebida.
	//r *http.Request -> Nele você encontra -> método HTTP, rota, body, headers, query params, contexto.
	//status int -> status code.
	//data any -> qualquer dado para virar JSON.

	w.Header().Set("Content-Type", "application/json")
	//Define o header -> Content-Type:application/json.
	//Content-Type:application/json -> avisa que a resposta está em JSON.

	w.WriteHeader(http.StatusOK)
	//Define o status HTTP -> 200(OK), 201(Created), 400(Bad Request).

	if err := json.NewEncoder(w).Encode(map[string]any{
		//json.NewEncoder(w).Encode(data) -> Transforma o dado Go em JSON e escreve na resposta.

		"status": "ok",
		"time":   time.Now(),
	}); err != nil {
		log.Println("failed to write response:", err)
		//Se der erro ao gerar JSON, ele mostra no terminal.
	}
}
