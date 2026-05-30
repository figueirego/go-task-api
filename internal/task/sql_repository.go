package task

import (
	"context"      //Usado para receber -> ctx context.Context.
	"database/sql" //Pacote padrão do Go para trabalhar com SQL.
	"errors"       //Usado para comparar erro.
)

type SQLRepository struct {
	//Essa struct guarda a conexão/pool com o banco.
	//Busca os dados no MySQL.

	db *sql.DB
}

func NewSQLRepository(db *sql.DB) *SQLRepository {
	//Essa função cria um repository SQL.
	//Ela recebe o banco aberto no main.go
	//Ex -> repository := task.NewSQLRepository.

	return &SQLRepository{
		db: db,
	}
}

func (r *SQLRepository) ListTasks(ctx context.Context) ([]Task, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, title, done, created_at
		FROM tasks
		ORDER BY id
	`)
	//Lista tasks do MySQL.
	//Usamos QueryContext quando esperamos várias linhas.
	//Significa -> busque as colunas id, title, done e created_at, da tabela tasks, ordenando pelo id.

	if err != nil {
		//Se o banco falhar, retornamos erro para o Handler responder 500.
		return nil, err
	}

	defer rows.Close()
	//Sempre que usamos QueryContext, precisamos fechar as linhas no final.
	//defer garante que isso vai acontecer ao sair da função.

	tasks := []Task{}
	//Cria uma slice vazia.

	for rows.Next() {
		//Enquanto houver linha no resulltado, entra no loop.

		var task Task
		//cria uma task vazia para preencher com os dados do banco.

		if err := rows.Scan(
			//Scan pega os valores da linha SQL e coloca dentro dos campos da struct.
			//Ordem importa.

			&task.ID,
			&task.Title,
			&task.Done,
			&task.CreatedAt,
		); err != nil {
			return nil, err
		}

		tasks = append(tasks, task)
		//Adiciona a task na lista.

	}

	if err := rows.Err(); err != nil {
		//Depois do loop, verifica se houve erro final.

		return nil, err
	}

	return tasks, nil
	//Retorna a lista e nenhum erro.

}

func (r *SQLRepository) CreateTask(ctx context.Context, title string) (Task, error) {
	//Esse método cria task no MySQL.

	result, err := r.db.ExecContext(ctx, `
		INSERT INTO tasks (title, done)
		VALUES (?, ?)
	`, title, false)
	//Usamos ExecContext quando a query não retorna linhas diretamente.
	//Ex -> INSERT, UPDATE, DELETE.
	//VALUES (?, ?) -> São placeholders, os valores reais vêm depois.
	//(?, ?) -> title, false. -> ajuda a evitar SQL injection.

	if err != nil {
		return Task{}, err
	}

	id, err := result.LastInsertId()
	//Depois do INSERT, o MySQL gera um ID automático.
	//Esse método pega o ID gerado.

	if err != nil {
		return Task{}, err
	}

	return r.FindTaskByID(ctx, int(id))
	//Depois de inserir, buscamos a task no banco.
	//Buscar de novo garante que retornamos a task completa.

}

func (r *SQLRepository) FindTaskByID(ctx context.Context, id int) (Task, error) {
	//Esse método busca uma task pelo ID.

	var task Task

	err := r.db.QueryRowContext(ctx, `
		SELECT id, title, done, created_at
		FROM tasks
		WHERE id = ?
	`, id).Scan(
		&task.ID,
		&task.Title,
		&task.Done,
		&task.CreatedAt,
	)
	//QueryRowContext -> Usamos quando esperamos uma única linha.

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Task{}, ErrTaskNotFound
		}
		//Quando QueryRowContext não encontra nenhuma linha, retorna esse erro (404).

		return Task{}, err
	}

	return task, nil
}

func (r *SQLRepository) UpdateTaskDone(ctx context.Context, id int, done bool) (Task, error) {
	//Esse método atualiza o campo done.

	result, err := r.db.ExecContext(ctx, `
		UPDATE tasks
		SET done = ?
		WHERE id = ?
	`, done, id)
	if err != nil {
		return Task{}, err
	}

	rowsAffected, err := result.RowsAffected()
	//Isso diz quantas linhas foram alteradas.
	//Se for 0, significa que nenhuma task com aquele ID existe.

	if err != nil {
		return Task{}, err
	}

	if rowsAffected == 0 {
		return Task{}, ErrTaskNotFound
	}
	//Se não encontrou, transformamos em erro conhecido.

	return r.FindTaskByID(ctx, id)
	//Depois de atualizar, buscamos a task atualizada e retornamos.

}

func (r *SQLRepository) DeleteTask(ctx context.Context, id int) error {
	//Esse método deleta uma task.

	result, err := r.db.ExecContext(ctx, `
		DELETE FROM tasks
		WHERE id = ?
	`, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrTaskNotFound
	}
	//Se rowsAffected == 0, nenhuma task foi deletada -> return ErrTaskNotFound.

	return nil
	//Se deletou retorna nil.

}
