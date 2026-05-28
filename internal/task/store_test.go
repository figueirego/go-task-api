package task

//Esse teste está dentro do mesmo package da Store.
//Usando o mesmo package task.
//Permite testar diretamente -> NewStore().

import "testing"

//Traz o pacote oficial de testes do Go.
//Permite criar testes usando funções como -> func TestAlgumaCoisa(t *testing.T) {}
//Como Go identifica teste -> Todo teste precisa começar com -> func Test...
//Exemplo -> func TestCreateTask(t *testing.T) {}
//t *testing.T -> Objeto que o Go entrega para o teste.
//t *testing.T -> t = marcar erro, falhar teste, mostrar mensagem, criar subtestes.
//t *testing.T -> t = t.Errorf(...).

func TestNewStoreStartsEmpty(t *testing.T) {
	//Testa se a Store começa vazia.

	store := NewStore()
	//Cria uma Store nova.

	tasks := store.ListTasks()
	//Lista as tarefas. (esperamos lista vazia)

	if len(tasks) != 0 {
		//len(tasks) pega o tamanho da lista.
		//Se o tamanho for diferente de 0, tem algo errado.

		t.Errorf("expected 0 tasks, got %d", len(tasks))
		//Marca o teste como erro e mostra uma mensagem.
		//%d é usado para número inteiro.
		//Se vier 3 tarefas, a mensagem seria -> expected 0 tasks, got 3.
	}
}

func TestCreateTask(t *testing.T) {
	//Verifica se uma task criada vem com os dados corretos.

	store := NewStore()
	//Começa do zero.

	task := store.CreateTask("Estudar testes em Go")
	//Cria uma task.

	if task.ID != 1 {
		t.Errorf("expected task ID 1, got %d", task.ID)
		//A primeira task precisa ter ID 1.

	}

	if task.Title != "Estudar testes em Go" {
		//O título retornado precisa ser igual ao título enviado.

		t.Errorf("expected title %q, got %q", "Estudar testes em Go", task.Title)
		//%q mostra texto com aspas. -> expected title "Estudar testes em Go", got "Outro título".

	}

	if task.Done != false {
		//A task precisa começar como não concluída.

		t.Errorf("expected task Done false, got %v", task.Done)
	}

	if task.CreatedAt.IsZero() {
		//IsZero() verifica se a data está vazia.

		t.Errorf("expected CreatedAt to be set")
	}
}

func TestCreateTaskIncrementsID(t *testing.T) {
	//Verifica se o ID aumenta.
	//Se as duas viessem com ID 1 teria bug.

	store := NewStore()

	first := store.CreateTask("Primeira tarefa")
	second := store.CreateTask("Segunda tarefa")

	if first.ID != 1 {
		t.Errorf("expected first ID, got %d", first.ID)
	}

	if second.ID != 2 {
		t.Errorf("expected second ID, got %d", second.ID)
	}
}

func TestListTasksReturnsCreatedTasks(t *testing.T) {
	//Verifica se a lista retorna as tasks criadas.

	store := NewStore()

	store.CreateTask("Estudar Go")
	store.CreateTask("Estudar Docker")
	//Não guardamos o retorno em variável porque não precisamos dele diretamente.
	//Queremos que elas sejam salvas na Store.

	tasks := store.ListTasks()
	//Pega a lista.

	if len(tasks) != 2 {
		//Verificando tamanho -> Esperamos duas tasks.

		t.Errorf("expected 2 tasks, got %d", len(tasks))
	}

	if tasks[0].Title != "Estudar Go" {
		//Acessando índice da lista -> tasks[0] = Primeira task.
		//tasks[1] = Segunda task.

		t.Errorf("expected first task title %q, got %q", "Estudar Docker", tasks[1].Title)
	}
}
