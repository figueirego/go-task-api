package task

import (
	"sync" //Serve para recursos de sincronização.
	"time" //Serve para trabalhar com data e hora.
)

type Store struct {
	//Store será nosso "banco de dados em memória".

	mu     sync.Mutex //Mutex serve para travar e destravar acesso a um dado.
	nextID int        //Guarda o próximo ID da tarefa.
	tasks  []Task     //Isso é uma lista de tarefas.
}

func NewStore() *Store {
	//func cria uma função.
	//É uma função que cria uma nova Store.
	//Essa função retorna um ponteiro para Store.
	//*Store = referência para uma Store

	return &Store{
		//Isso cria uma Store e retorna o endereço dela.
		//& = Pegue o endreço da memória.

		nextID: 1,        //A primeira tarefa terá ID 1.
		tasks:  []Task{}, //Começamos com uma lista vazia (slice vazio de task)
	}
}

func (s *Store) ListTasks() []Task {
	//Isso é um método.
	//Um método é uma função ligada a um tipo.
	//ListTasks pertence a Store.
	//(s *Store) = receiver ->  Significa essa funçao pertence a Store.
	// Na função Store = s. -> ex: s.tasks, s.mu...
	//[]Task = Função retorna uma lista de Task.

	s.mu.Lock()
	//Trava o mutex
	//s.mu.Lock() = vou mexer ou ler dados protegidos. Ninguém mais pode mexer agora.

	defer s.mu.Unlock()
	//defer = execute isso no final da função
	//defer s.mu.Unlock() = quando a função terminar, destrave o mutex.

	copied := make([]Task, len(s.tasks))
	//"copied :=" = crie uma variável nova chamada copied
	//O := é uma forma curta de criar variável em Go.
	//make([]Task, len(s.tasks)) = Cria uma lista nova de Task com o mesmo tamanho de s.tasks.
	//len(s.tasks) pega o tamanho da lista.
	//Ex: len(s.tasks) = 3, então, make([]Task, 3) cria uma lista vazia com 3 posições.

	copy(copied, s.tasks)
	//Copia as tarefas de s.tasks para copied
	//É uma proteção simples, para não devolver a lista interna diretamente.

	return copied
	//Retorna a cópia da lista de tarefas

}

func (s *Store) CreateTask(title string) Task {
	//CreateTask é um método de Store.
	//Ele recebe um title do tipo string.
	//Ele retorna uma Task
	//Parâmetro -> tittle string = A função precisa receber um título.

	s.mu.Lock()
	defer s.mu.Unlock()

	//Trava antes de mexer na lista
	//Destrava ao terminar

	task := Task{
		ID:        s.nextID,   //Usa o próximo ID disponível.
		Title:     title,      //Usa o título recebido pela função.
		Done:      false,      //A tarefa começa como não concluída.
		CreatedAt: time.Now(), //Define a data/hora atual.
	}

	s.nextID++
	//Incrementando ID
	//Isso aumenta nextID em 1.
	//s.nextID++ -> s.nextID = s.nextID + 1

	s.tasks = append(s.tasks, task)
	//append adiciona item em uma slice.
	//Antes = [], Depois = [task].
	//Se antes = [task1], Depois = [task1, task2].

	return task
	//Retorna a tarefa criada para a api responder ao usuário.

}
