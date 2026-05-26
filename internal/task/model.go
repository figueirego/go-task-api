package task

import "time" //Serve para trabalhar com data e hora.

type Task struct {
	//type cria um novo tipo. (Tipo = Task)
	//struct é uma estrutura de dados.

	ID        int       `json:"id"`        //A tarefa tem um ID do tipo inteiro.
	Title     string    `json:"title"`     //A tarefa tem um título do tipo texto.
	Done      bool      `json:"done"`      //A tarefa pode estar concluída ou não.
	CreatedAt time.Time `json:"createdAt"` //A tarefa tem uma data/hora de criação.

	//Ex: Campo=ID, Tipo=int, tag='json:"id"' (como campo vai aparecer em JSON)
}
