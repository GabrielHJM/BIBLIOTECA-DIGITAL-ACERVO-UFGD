-- Consolidated Database Schema for Biblioteca Digital

-- Clean start (Optional: only use if you want to wipe the database)
-- DROP TABLE IF EXISTS mensagens CASCADE;
-- DROP TABLE IF EXISTS amizades CASCADE;
-- DROP TABLE IF EXISTS comentarios CASCADE;
-- DROP TABLE IF EXISTS curtidas CASCADE;
-- DROP TABLE IF EXISTS emprestimos CASCADE;
-- DROP TABLE IF EXISTS historico_leitura CASCADE;
-- DROP TABLE IF EXISTS avaliacoes CASCADE;
-- DROP TABLE IF EXISTS favoritos CASCADE;
-- DROP TABLE IF EXISTS interesses_usuario CASCADE;
-- DROP TABLE IF EXISTS flashcards CASCADE;
-- DROP TABLE IF EXISTS materiais CASCADE;
-- DROP TABLE IF EXISTS usuarios CASCADE;

-- Extensions
CREATE EXTENSION IF NOT EXISTS unaccent;

-- Tabela de Usuários
CREATE TABLE IF NOT EXISTS usuarios (
    id SERIAL PRIMARY KEY,
    nome TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    senha TEXT NOT NULL,
    tipo INTEGER DEFAULT 1, -- 1: Estudante, 2: Professor, 3: Admin
    foto_url TEXT,
    meta_paginas_semana INTEGER DEFAULT 100,
    data_criacao TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Tabela de Materiais (Livros/Artigos)
CREATE TABLE IF NOT EXISTS materiais (
    id SERIAL PRIMARY KEY,
    titulo TEXT NOT NULL,
    autor TEXT NOT NULL,
    isbn TEXT,
    categoria TEXT NOT NULL,
    ano_publicacao INTEGER,
    descricao TEXT,
    capa_url TEXT,
    pdf_url TEXT,
    disponivel BOOLEAN DEFAULT TRUE,
    media_nota NUMERIC(3,2) DEFAULT 0.0,
    total_avaliacoes INTEGER DEFAULT 0,
    paginas INTEGER DEFAULT 0,
    externo_id TEXT UNIQUE,
    fonte TEXT,
    status TEXT DEFAULT 'aprovado',
    curador_id INTEGER REFERENCES usuarios(id),
    dificuldade INTEGER DEFAULT 1,
    xp INTEGER DEFAULT 10,
    relevancia INTEGER DEFAULT 0,
    search_vector tsvector
);

-- Tabela de Flashcards
CREATE TABLE IF NOT EXISTS flashcards (
    id SERIAL PRIMARY KEY,
    usuario_id INTEGER NOT NULL REFERENCES usuarios(id) ON DELETE CASCADE,
    material_id INTEGER REFERENCES materiais(id) ON DELETE SET NULL,
    pergunta TEXT NOT NULL,
    resposta TEXT NOT NULL,
    dificuldade INTEGER DEFAULT 0,
    proxima_revisao TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    data_criacao TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Tabela de Favoritos
CREATE TABLE IF NOT EXISTS favoritos (
    usuario_id INTEGER REFERENCES usuarios(id) ON DELETE CASCADE,
    material_id INTEGER REFERENCES materiais(id) ON DELETE CASCADE,
    PRIMARY KEY (usuario_id, material_id)
);

-- Tabela de Histórico de Leitura
CREATE TABLE IF NOT EXISTS historico_leitura (
    id SERIAL PRIMARY KEY,
    usuario_id INTEGER REFERENCES usuarios(id) ON DELETE CASCADE,
    material_id INTEGER REFERENCES materiais(id) ON DELETE CASCADE,
    data TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Tabela de Interesses do Usuário
CREATE TABLE IF NOT EXISTS interesses_usuario (
    usuario_id INTEGER REFERENCES usuarios(id) ON DELETE CASCADE,
    interesse TEXT NOT NULL,
    PRIMARY KEY (usuario_id, interesse)
);

-- Tabela de Avaliações
CREATE TABLE IF NOT EXISTS avaliacoes (
    id SERIAL PRIMARY KEY,
    usuario_id INTEGER REFERENCES usuarios(id) ON DELETE CASCADE,
    material_id INTEGER REFERENCES materiais(id) ON DELETE CASCADE,
    nota INTEGER CHECK (nota >= 1 AND nota <= 5),
    comentario TEXT,
    data TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Tabela de Empréstimos
CREATE TABLE IF NOT EXISTS emprestimos (
    id SERIAL PRIMARY KEY,
    usuario_id INTEGER REFERENCES usuarios(id) ON DELETE CASCADE,
    material_id INTEGER REFERENCES materiais(id) ON DELETE CASCADE,
    data_emprestimo TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    data_devolucao TIMESTAMP,
    status TEXT DEFAULT 'ativo'
);

-- Tabela de Curtidas
CREATE TABLE IF NOT EXISTS curtidas (
    usuario_id INTEGER REFERENCES usuarios(id) ON DELETE CASCADE,
    material_id INTEGER REFERENCES materiais(id) ON DELETE CASCADE,
    data TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (usuario_id, material_id)
);

-- Tabela de Comentários
CREATE TABLE IF NOT EXISTS comentarios (
    id SERIAL PRIMARY KEY,
    usuario_id INTEGER REFERENCES usuarios(id) ON DELETE CASCADE,
    material_id INTEGER REFERENCES materiais(id) ON DELETE CASCADE,
    texto TEXT NOT NULL,
    data TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Tabela de Amizades
CREATE TABLE IF NOT EXISTS amizades (
    usuario_id1 INTEGER REFERENCES usuarios(id) ON DELETE CASCADE,
    usuario_id2 INTEGER REFERENCES usuarios(id) ON DELETE CASCADE,
    status TEXT DEFAULT 'pendente',
    data_solicitacao TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (usuario_id1, usuario_id2),
    CHECK (usuario_id1 != usuario_id2)
);

-- Tabela de Anotações
CREATE TABLE IF NOT EXISTS anotacoes (
    id SERIAL PRIMARY KEY,
    usuario_id INTEGER NOT NULL REFERENCES usuarios(id) ON DELETE CASCADE,
    material_id INTEGER REFERENCES materiais(id) ON DELETE SET NULL,
    titulo TEXT,
    conteudo TEXT NOT NULL,
    cor TEXT DEFAULT '#FFFFFF',
    data_criacao TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    data_atualizacao TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Tabela de Mensagens
CREATE TABLE IF NOT EXISTS mensagens (
    id SERIAL PRIMARY KEY,
    remetente_id INTEGER REFERENCES usuarios(id) ON DELETE CASCADE,
    destinatario_id INTEGER REFERENCES usuarios(id) ON DELETE CASCADE,
    material_id INTEGER REFERENCES materiais(id) ON DELETE SET NULL,
    conteudo TEXT,
    lida BOOLEAN DEFAULT FALSE,
    data TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for Performance
CREATE INDEX IF NOT EXISTS idx_materiais_search_vector ON materiais USING gin(search_vector);
CREATE INDEX IF NOT EXISTS idx_materiais_status ON materiais(status);
CREATE INDEX IF NOT EXISTS idx_materiais_categoria ON materiais(categoria);
CREATE INDEX IF NOT EXISTS idx_materiais_fonte ON materiais(fonte);
CREATE INDEX IF NOT EXISTS idx_materiais_ano ON materiais(ano_publicacao);
CREATE INDEX IF NOT EXISTS idx_materiais_externo_id ON materiais(externo_id) WHERE externo_id IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_historico_usuario ON historico_leitura(usuario_id);

-- Full-Text Search Triggers
CREATE OR REPLACE FUNCTION materiais_search_trigger() RETURNS trigger AS $$
BEGIN
  new.search_vector :=
    setweight(to_tsvector('portuguese', coalesce(new.titulo,'')), 'A') ||
    setweight(to_tsvector('portuguese', coalesce(new.autor,'')), 'B') ||
    setweight(to_tsvector('portuguese', coalesce(new.descricao,'')), 'C');
  RETURN NEW;
END
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS tsvectorupdate ON materiais;
CREATE TRIGGER tsvectorupdate BEFORE INSERT OR UPDATE
ON materiais FOR EACH ROW EXECUTE FUNCTION materiais_search_trigger();
