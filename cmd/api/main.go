package main

//todo arquivo Go começa com um package.
//Um package é um agrupamento de código.
//package main = este projeto gera um executável.
//Se fosse uma biblioteca, poderia ser package task ou package users.

import (
	//import serve para trazer código pronto de outros pacotes.

	"encoding/json" //Serve para trabalhar com JSON.
	"errors"        //Usamos para comparar erros.
	"log"           //Serve para imprimir logs no terminal.
	"net/http"      //Esse é o pacote HTTP nativo do Go. Permite criar servidor web sem framework externo.
	"strconv"       //Serve para converter string para número e número para string.
	"strings"       //Serve para manipular textos.
	"time"          //Serve para trabalhar com data e hora.

	"github.com/jomatheusdev/go-task-api/internal/task"
	//Acessa o pacote interno e usa task.NewStore().
)

type CreateTaskRequest struct {
	//Essa stuct representa o corpo da requisição POST /tasks.

	Title string `json:"title"`
}

type UpdateTaskRequest struct {
	//Essa struct representa o body do PATCH.

	Done bool `json:"done"`
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

	mux.HandleFunc("GET /tasks/{id}", func(w http.ResponseWriter, r *http.Request) {
		//{id} -> Parâmetro de rota.

		id, ok := parseIDParam(w, r)
		if !ok {
			return
		}
		//Parse do Id -> Chamamos uma função auxiliar para extrair e validar o ID.
		//Retorna -> int e bool. -> id convertido e se deu certo ou não.
		//Se ok for false, paramos com return.

		foundTask, err := store.FindTaskById(id)
		//Busca a task na Store.
		//Deu certo -> fondTask = task encontrada, err = nil.
		//Deu errado -> fondTask = task vazia, err = ErrTaskNotFound.

		if err != nil {
			//Se deu erro, entramos nesse bloco.

			if errors.Is(err, task.ErrTaskNotFound) {
				//Se o erro for "task not found" -> erro 404.
				writeJSON(w, http.StatusNotFound, map[string]string{
					//http.StatusNotFound -> status 404.
					//significa recurso não encontrado.

					"error": "task not found",
				})
				return
			}

			writeJSON(w, http.StatusInternalServerError, map[string]string{
				//Erro inesperado
				//http.StatusInternalServerError -> status 500.
				//Significa erro interno no servidor.

				"error": "internal server error",
			})
			return
		}

		writeJSON(w, http.StatusOK, foundTask)
	})

	mux.HandleFunc("PATCH /tasks/{id}", func(w http.ResponseWriter, r *http.Request) {
		//PATCH é usado para atualizar parcialmente um recurso.
		//Não vamos atualizar a task inteira, apenas "done", por isso PATCH.

		id, ok := parseIDParam(w, r)
		if !ok {
			return
		}

		var body UpdateTaskRequest

		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			//Mesma lógica do POST.
			//Criamos uma variável vazia e tentamos preencher com o JSON da request.

			writeJSON(w, http.StatusBadRequest, map[string]string{
				"error": "invalid JSON body",
			})
			return
		}

		updatedTask, err := store.UpdateTaskDone(id, body.Done)
		//Chama o método da Store.
		//Se achar a tarefa, atualiza.
		//Se não achar, retorna erro.

		if err != nil {
			if errors.Is(err, task.ErrTaskNotFound) {
				writeJSON(w, http.StatusNotFound, map[string]string{
					"error": "task not found",
				})
				return
			}
			writeJSON(w, http.StatusInternalServerError, map[string]string{
				"error": "internal server error",
			})
			return
		}

		writeJSON(w, http.StatusOK, updatedTask)
	})

	mux.HandleFunc("DELETE /tasks/{id}", func(w http.ResponseWriter, r *http.Request) {
		//DELETE é usado para deletar um recurso.
		//Significa -> delete a task de ID 1.

		id, ok := parseIDParam(w, r)
		if !ok {
			return
		}

		if err := store.DeleteTask(id); err != nil {
			if errors.Is(err, task.ErrTaskNotFound) {
				writeJSON(w, http.StatusNotFound, map[string]string{
					"error": "task not found",
				})
				return
			}
			writeJSON(w, http.StatusInternalServerError, map[string]string{
				"error": "internal server error",
			})
			return
		}

		w.WriteHeader(http.StatusNoContent)
		//http.StatusNoContent -> status 204.
		//Significa -> Deu certo, mas não há conteúdo para retornar.
		//Por isso não usamos writeJSON aqui.

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

func parseIDParam(w http.ResponseWriter, r *http.Request) (int, bool) {
	//Nova função auxiliar.
	//Recebe -> w = resposta HTTP, r = request HTTP.
	//Retorna -> int = ID convertido, bool = true se válido, false se inválido.

	rawID := r.PathValue("id")
	//PathValue("id") pega o valor da rota.

	id, err := strconv.Atoi(rawID)
	//Converte string para int.

	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "invalid task id",
		})
		return 0, false
	}
	//Se o usuário chamar string (/tasks/abc, responde com invalid task id. (400)

	if id <= 0 {
		//Não queremos aceitar -> /tasks/0 ou /tasks/-1

		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "invalid id must be positive",
		})
		return 0, false
	}
	return id, true
	//Retorno de sucesso.

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
