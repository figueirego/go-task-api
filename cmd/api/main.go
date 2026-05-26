package main

//todo arquivo Go começa com um package.
//Um package é um agrupamento de código.
//package main = este projeto gera um executável.
//Se fosse uma biblioteca, poderia ser package task ou package users.

import (
	//import serve para trazer código pronto de outros pacotes.

	"encoding/json" //Serve para trabalhar com JSON.
	"log"           //Serve para imprimir logs no terminal.
	"net/http"      //Esse é o pacote HTTP nativo do Go. Permite criar servidor web sem framework externo.
	"strings"       //Serve para manipular textos.
	"time"          //Serve para trabalhar com data e hora.

	"github.com/jomatheusdev/go-task-api/internal/task"
	//Acessa o pacote interno e usa task.NewStore().
)

type CreateTaskRequest struct {
	//Essa stuct representa o corpo da requisição POST /tasks.

	Title string `json:"title"`
}

func main() {
	//func main() = Função inicial do programa.

	store := task.NewStore()
	//Cria nossa Store em memória (Lugar para guardar tarefas).

	mux := http.NewServeMux()
	//Cria um roteador HTTP.
	//mux = muxtiplexer -> roteador de rotas HTTP.

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		//mux.HandleFunc = Registra uma rota.
		//Formato -> mux.HandleFunc("MÉTODO /caminho", função).
		//"GET /healt" -> Quando alguém fizer GET /health, execute essa função.
		//GET = Método HTTP para buscar dados.
		// "/health" = Rota de saúde da aplicação.
		//func(w http.ResponseWriter, r *http.Request) {} -> função anônima.
		//w http.ResponseWriter -> w = Usado para escrever a resposta.
		//w http.ResponseWriter -> com ele voce retorna -> status code, headers, JSON, texto, erro.
		//r *http.Request -> r = representa a requisição recebida.
		//r *http.Request -> Nele você encontra -> método HTTP, rota, body, headers, query params, contexto.

		writeJSON(w, http.StatusOK, map[string]any{
			//Função auxiliar -> Responde JSON.
			//http.StatusOK = status HTTP 200.
			//map[string]any = objeto chave-valor -> "status": "ok", "time":   time.Now().
			//map[string]any = mapa onde as chaves são strings, e os valores podem ser qualquer tipo.

			"status": "ok",
			"time":   time.Now(),
		})
	})

	mux.HandleFunc("GET /tasks", func(w http.ResponseWriter, r *http.Request) {
		tasks := store.ListTasks()
		//Busca as tarefas guardadas na memória.

		writeJSON(w, http.StatusOK, tasks)
		//Responde a lista em JSON.

	})

	mux.HandleFunc("POST /tasks", func(w http.ResponseWriter, r *http.Request) {
		//POST = usado para criar dados.

		var body CreateTaskRequest
		//Cria uma variável chamada body.
		//Ela começa vazia.
		//Tipo -> CreateTaskRequest.

		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			//r.Body = Corpo da requisição.
			//json.NewDecoder(r.Body) = Leitor de JSON a partir do body.
			//.Decode(&body) = Tenta transformar o JSON em struct Go.
			//&body = passe o endereço da variável body, para que o Decode consiga preenchê-la.
			//Tratamento de erro -> if err := ...; != nil {} (padrão comum em Go).
			//if err := ...; != nil {} = tente fazer algo; se der erro, trate o erro.

			writeJSON(w, http.StatusBadRequest, map[string]string{
				//http.StatusBadRequest = HTTP 400 (Dados inválidos).

				"error": "invalid JSON body",
			})
			return
			//parar a execução da função.
		}
		title := strings.TrimSpace(body.Title)
		//Pega o título enviado e remove espaços.

		if title == "" {
			//Se o título for vazio, retorna erro.
			writeJSON(w, http.StatusBadRequest, map[string]string{
				"error": "title is required",
			})
			return
		}

		task := store.CreateTask(title)
		//Chama nosso método -> CreateTask.
		//CreateTask -> Cria a tarefa, salva na lista e retorna a tarefa criada.

		writeJSON(w, http.StatusCreated, task)
		//http.StatusCreated = 201 (recurso criado com sucesso).

	})
	log.Println("API running on http://localhost:8080")
	//Mostra no terminal -> saber que o servidor iniciou.

	log.Fatal(http.ListenAndServe(":8080", mux))
	//Inicia o servidor HTTP.
	//":8080" -> escute na porta 8080.
	//mux -> roteador com as rotas: GET /health, GET /tasks, POST /tasks.
	//ListenAndServe = fica rodando enquanto o servidor está ativo.
	//Se der erro -> log.Fatal -> imprime o erro e encerra o programa.

}

func writeJSON(w http.ResponseWriter, status int, data any) {
	//w http.ResponseWriter -> resposta HTTP.
	//status int -> status code.
	//data any -> qualquer dado para virar JSON.

	w.Header().Set("Content-Type", "application/json")
	//Define o header -> Content-Type:application/json.
	//Content-Type:application/json -> avisa que a resposta está em JSON.

	w.WriteHeader(status)
	//Define o status HTTP -> 200(OK), 201(Created), 400(Bad Request).

	if err := json.NewEncoder(w).Encode(data); err != nil {
		//json.NewEncoder(w).Encode(data) -> Transforma o dado Go em JSON e escreve na resposta.

		log.Println("failed to encode response:", err)
		//Se der erro ao gerar JSON, ele mostra no terminal.
	}
}
