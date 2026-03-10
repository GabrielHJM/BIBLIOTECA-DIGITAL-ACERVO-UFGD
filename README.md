<div align="center">
  <img src="https://capsule-render.vercel.app/api?type=waving&color=00ADD8&height=200&section=header&text=Biblioteca%20Digital%20Colaborativa&fontSize=40&fontAlignY=38&desc=Revolucionando%20o%20acesso%20acadêmico&descAlignY=55&descAlign=50" alt="Banner Biblioteca Digital" width="100%" />

  <br>

  [![Vue.js](https://img.shields.io/badge/Vue.js-35495E?style=for-the-badge&logo=vue.js&logoColor=4FC08D)](https://vuejs.org/)
  [![Go](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://go.dev/)
  [![PostgreSQL](https://img.shields.io/badge/PostgreSQL-316192?style=for-the-badge&logo=postgresql&logoColor=white)](https://www.postgresql.org/)
  [![Status](https://img.shields.io/badge/Status-TCC_Concluído-success?style=for-the-badge)](#)

  <p align="center">
    <strong>A Biblioteca Digital Colaborativa unifica acervos acadêmicos da UFGD em uma única plataforma. Com arquitetura robusta em Go, Vue.js e PostgreSQL, ela moderniza a pesquisa. O sistema democratiza o acesso à informação, centralizando buscas e a leitura de materiais para transformar o estudo universitário de forma ágil, eficiente e inovadora.</strong>
  </p>
</div>

---

# 📖 PARTE I: APRESENTAÇÃO DO PROJETO

## O Resgate do Conhecimento Centralizado
Historicamente, as bibliotecas físicas eram o ponto central inquestionável do saber acadêmico. Com a digitalização, paradoxalmente, o conhecimento fragmentou-se. Livros, teses e artigos acabaram isolados em dezenas de repositórios institucionais que não se comunicam. 

Desenvolvido como Trabalho de Conclusão de Curso (TCC) em Sistemas de Informação pela Universidade Federal da Grande Dourados (UFGD), este projeto tem uma visão clara: **resgatar o valor do acervo unificado**, utilizando tecnologia de ponta para centralizar o ecossistema educacional. A solução otimiza o tempo do estudante, conectando-o diretamente à fonte de estudo.

## 💻 Telas e Interface (Native-First)

A experiência do usuário foi priorizada para ser fluida, livre de distrações e orientada à conversão em leitura.

| Autenticação Segura | Vitrine Virtual | Leitura e Download |
| :---: | :---: | :---: |
| <img src="docs/login.png" width="250" alt="Login"> | <img src="docs/home.png" width="250" alt="Home"> | <img src="docs/detalhes.png" width="250" alt="Detalhes"> |
| *Acesso restrito e seguro.* | *Exploração em tempo real.* | *Acesso imediato ao material.* |

> **Nota para deploy:** Substituir `docs/*.png` pelas imagens reais armazenadas no repositório.

## ✨ Visão Geral das Funcionalidades
- 📚 **Vitrine Inteligente:** Organização visual intuitiva por áreas do conhecimento.
- 🧠 **Espaço de Estudo:** Suporte a Flashcards, marcações de favoritos e anotações.
- 📱 **PWA Ready:** Experiência de aplicativo instalável via Service Workers.

---

# ⚙️ PARTE II: ESPECIFICAÇÕES TÉCNICAS

Esta seção detalha a engenharia por trás da plataforma, construída para suportar alta concorrência com o mínimo de sobrecarga de hardware.

## 🛠️ Stack Tecnológica

<div align="center">
  <table>
    <tr>
      <td align="center" width="33%"><b>Frontend (Port: 8081)</b><br><br>Vue.js 3 (Composition API)<br>Vuetify 3<br>GSAP<br>Axios</td>
      <td align="center" width="33%"><b>Backend (Port: 8080)</b><br><br>Golang 1.25<br>Clean Architecture<br>Middlewares (Zap, Cors)</td>
      <td align="center" width="33%"><b>Database & Infra (Port: 5432)</b><br><br>PostgreSQL<br>Redis (Cache Opcional)<br>Node.js (Concurrently)</td>
    </tr>
  </table>
</div>

## 🏗️ Arquitetura do Sistema

O ecossistema é dividido em aplicações independentes que se comunicam estritamente via API RESTful. 

```mermaid
graph TD
    subgraph Frontend [SPA / PWA]
        V[Vue.js 3 App]
        VR[Vue Router]
    end

    subgraph Backend [Golang API]
        G[Go 1.25 HTTP Handlers]
        M[Automated Migrations]
        H[Harvester Background Sync]
    end

    subgraph Database [Relational & Search]
        P[PostgreSQL Local]
        FTS[Full-Text Search Engine]
    end

    V -- Axios (JSON HTTP) --> G
    G -- Queries & GORM --> P
    G -- Boot Setup --> M
    G -- Goroutine (30min) --> H
    H -- Sync APIs Externas --> G
    P -- tsvector / unaccent --> FTS
