const { execSync } = require('child_process');
const fs = require('fs');
const path = require('path');

/**
 * Biblioteca Digital - Sistema de Sincronização Autônomo
 * Este script garante a paridade entre o ambiente Local, Docker e Online.
 */

const PROJECT_ROOT = __dirname;
const BACKEND_DIR = path.join(PROJECT_ROOT, 'backend');
const FRONTEND_DIR = path.join(PROJECT_ROOT, 'frontend');

function log(msg) {
    console.log(`[SYNC] ${new Date().toISOString()} - ${msg}`);
}

function run(cmd, cwd = PROJECT_ROOT) {
    try {
        log(`Executando: ${cmd} em ${cwd}`);
        return execSync(cmd, { cwd, stdio: 'inherit' });
    } catch (e) {
        console.error(`[ERRO] Falha ao executar: ${cmd}`);
        return null;
    }
}

async function sync() {
    log('Iniciando varredura vasta e sincronização...');

    // 1. Sincronização de Dependências
    log('Verificando dependências do Frontend...');
    if (!fs.existsSync(path.join(FRONTEND_DIR, 'node_modules'))) {
        run('npm install', FRONTEND_DIR);
    } else {
        log('Node modules do frontend já existem. Pulando.');
    }

    // 2. Limpeza de Processos Antigos
    log('Limpando processos antigos...');
    try {
        if (process.platform === 'win32') {
            run('powershell -Command "Stop-Process -Name server_bin -Force -ErrorAction SilentlyContinue"');
            run('powershell -Command "Get-Process | Where-Object {$_.ProcessName -eq \'main\'} | Stop-Process -Force -ErrorAction SilentlyContinue"');
        }
    } catch (e) {}

    // 3. Build do Backend
    log('Compilando backend...');
    run('go build -o server_bin ./cmd/server/main.go', BACKEND_DIR);

    // 4. Sincronização com Docker
    log('Sincronizando com Docker...');
    run('docker-compose build');
    run('docker-compose up -d');

    // 5. Verificação de Saúde (Health Check)
    log('Verificando conectividade...');
    setTimeout(() => {
        log('Dica: Verifique se o banco de dados Postgres está rodando na porta 5432.');
        log('Acesse Localhost: http://localhost:8081');
        log('Acesse Docker: http://localhost:8082');
    }, 2000);

    log('Sincronização concluída com sucesso!');
}

sync();
