package task

//Esse arquivo pertence ao package task.
//Pode usar coisas do mesmo package como Store, ErrTaskNotFound...
//Não precisa importar -> "github.com/jomatheusdev/go-task-api/internal/task".
//"github.com/jomatheusdev/go-task-api/internal/task", já está dentro do proprio package task.

import (
	"encoding/json" //Serve para ler e escrever JSON.
	"errors"        //Serve para comparar erros.
	"log"           //Serve para imprimir erro no terminal.
	"net/http"      //Pacote HTTP nativo do Go
	"strconv"       //Converte string para int, ou int para string.
	"strings"       //Remove espaços do começo e do fim.
)

type Handler struct {
	//Criamos um tipo chamado Handler.

	store *Store
	//O Handler guarda uma refêrencia para a Store.
	//Store é onde estão as tarefas em memória
	//Em main.go acessava diretamente, agora quem faz isso é o Handler.
	//store = quem está fora do package task, não pode acessar diretamente.
	//handler.store = apenas no package task (protege a estrutura interna).

}

type CreateTaskRequest struct {
	Title string `json:"title"`
	//Representa o JSON recebido no POST /tasks.

}

type UpdateTaskDoneRequest struct {
	Done bool `json:"done"`
	//Representa o JSON recebido no PATCH /tasks/{id}.

}

func NewHandler(store *Store) *Handler {
	//Essa função cria um novo Handler.
	//Recebe uma Store.
	//Retorna um ponteiro para Handler.
	//store *Store -> Handler precisa manipular tarefas.
	//Store é uma dependência do Handler -> injeção de dependêndia (simples).
	//Injeção de dependência -> É quando você entrega para um objeto/função aquilo que ele precisa para trabalhar.
	//store := task.NewStore(), handler := task.NewHandler(store).
	//main monta as dependências, handler só usa as dependências.

	return &Handler{
		store: store,
	}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	//Esse método registra as rotas de tarefas.
	//(h *Handler) RegisterRoutes() -> Significa que é um método do tipo Handler.
	//mux *http.ServeMux -> mux é o roteador HTTP

	mux.HandleFunc("GET /tasks", h.ListTasks)
	mux.HandleFunc("POST /tasks", h.CreateTask)
	mux.HandleFunc("GET /tasks/{id}", h.FindTaskByID)
	mux.HandleFunc("PATCH /tasks/{id}", h.UpdateTaskDone)
	mux.HandleFunc("DELETE /tasks/{id}", h.DeleteTask)
}

func (h *Handler) ListTasks(w http.ResponseWriter, r *http.Request) {
	//Esse método responde -> GET /tasks
	//Parâmetros -> w http.ResponseWriter, r *http.Request.
	//w = usado para responder.
	//r = representa a request recebida.

	tasks := h.store.ListTasks()
	//Busca tarefas.
	//Antes = store.ListTasks().
	//Agora = h.store.ListTasks().

	writeJSON(w, http.StatusOK, tasks)
	//Responder JSON.
	//Retorna -> status 200 e lista de tasks em JSON.

}

func (h *Handler) CreateTask(w http.ResponseWriter, r *http.Request) {
	//Esse método responde POST /task.
	//Lê JSON.
	//Valida title.
	//Cria task.
	//Responde task criada.

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
	task := h.store.CreateTask(title)
	//Chama nosso método -> CreateTask.
	//CreateTask -> Cria a tarefa, salva na lista e retorna a tarefa criada.

	writeJSON(w, http.StatusCreated, task)
	//http.StatusCreated = 201 (recurso criado com sucesso).

}

func (h *Handler) FindTaskByID(w http.ResponseWriter, r *http.Request) {
	//Esse método responde GET /tasks/{id}.
	//Pega ID da URL.
	//Busca task na Store.
	//Se não achar, responde 404.
	//se achar, responde task.

	id, ok := parseIDParam(w, r)
	if !ok {
		return
	}
	//Parse do Id -> Chamamos uma função auxiliar para extrair e validar o ID.
	//Retorna -> int e bool. -> id convertido e se deu certo ou não.
	//Se ok for false, paramos com return.

	foundTask, err := h.store.FindTaskByID(id)
	//Busca a task na Store.
	//Deu certo -> fondTask = task encontrada, err = nil.
	//Deu errado -> fondTask = task vazia, err = ErrTaskNotFound.

	if err != nil {
		//Se deu erro, entramos nesse bloco.
		handleStoreError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, foundTask)
}

func (h *Handler) UpdateTaskDone(w http.ResponseWriter, r *http.Request) {
	//Esse método responde PATCH /tasks/{id}.
	//Pega ID da URL.
	//Lê JSON.
	//Pega done.
	//Atualiza task.
	//Responde task atualizada.

	id, ok := parseIDParam(w, r)
	if !ok {
		return
	}

	var body UpdateTaskDoneRequest

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		//Mesma lógica do POST.
		//Criamos uma variável vazia e tentamos preencher com o JSON da request.

		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "invalid JSON body",
		})
		return
	}

	updatedTask, err := h.store.UpdateTaskDone(id, body.Done)
	//Chama o método da Store.
	//Se achar a tarefa, atualiza.
	//Se não achar, retorna erro.

	if err != nil {
		handleStoreError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, updatedTask)
}

func (h *Handler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	//Esse método responde DELETE /tasks/{id}.
	//Pega ID da URL.
	//Deleta task.
	//Se não achar, responde 404.
	//Se deletar, responde 204.

	id, ok := parseIDParam(w, r)
	if !ok {
		return
	}

	if err := h.store.DeleteTask(id); err != nil {
		handleStoreError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	//http.StatusNoContent -> status 204.
	//Significa -> Deu certo, mas não há conteúdo para retornar.
	//Por isso não usamos writeJSON aqui.

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

func handleStoreError(w http.ResponseWriter, err error) {
	//Centraliza o tratamento de erro da Store.

	if errors.Is(err, ErrTaskNotFound) {
		writeJSON(w, http.StatusNotFound, map[string]string{
			"error": "task not found",
		})
		return
	}

	writeJSON(w, http.StatusInternalServerError, map[string]string{
		"error": "internal server error",
	})
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

		log.Println("failed to write response:", err)
		//Se der erro ao gerar JSON, ele mostra no terminal.
	}
}
