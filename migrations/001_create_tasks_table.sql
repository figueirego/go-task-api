CREATE TABLE IF NOT EXISTS tasks (
--     CREATE TABLE -> Cria a tabela tasks.
--     IF NOT EXISTS -> Isso evita erro caso a tabela já exista.

    id INT AUTO_INCREMENT PRIMARY KEY,
--     id -> coluna de identificação.
--     INT -> número inteiro.
--     AUTO_INCREMENT -> MySQL gera automaticamente.
--     PRIMARY KEY -> identificador único da tabela.

    title VARCHAR(255) NOT NULL,
--     title -> título da task.
--     VARCHAR(255) -> texto até 255 caracteres.
--     NOT NULL -> obrigatório.

    done BOOLEAN NOT NULL DEFAULT FALSE,
--     done -> concluída ou não.
--     BOOLEAN -> verdadeiro/false
--     NOT NULL -> obrigatório.
--     DEFAULT FALSE -> começa como false.

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
--     created_at -> data de criação.
--     TIMESTAMP -> data/hora.
--     NOT NULL -> obrigatório.
--     DEFAULT CURRENT_TIMESTAMP -> MySQL coloca a hora atual automaticamente.

);