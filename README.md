# goBid

![Go](https://img.shields.io/badge/Go-00ADD8?style=flat&logo=go&logoColor=white)
![License](https://img.shields.io/badge/license-MIT-blue.svg)
![Build Status](https://img.shields.io/badge/build-passing-brightgreen.svg)

**goBid** é uma aplicação de leilão online desenvolvida em Go, projetada para ser rápida, escalável e fácil de usar. Com ela, você pode criar leilões, gerenciar lances em tempo real e oferecer uma experiência fluida para usuários participarem de leilões digitais. Ideal para quem quer construir plataformas de bidding com performance e confiabilidade.

## ✨ Funcionalidades

- **Criação de Leilões**: Configure leilões com descrição, preço inicial e tempo de duração.
- **Lances em Tempo Real**: Suporte a lances via WebSocket para atualizações instantâneas.
- **Gerenciamento de Usuários**: Cadastro e autenticação de usuários para participação segura nos leilões.
- **Histórico de Lances**: Registro completo de todos os lances feitos em cada leilão.
- **Notificações**: Alertas para lances superados e fim de leilões.
- **API RESTful**: Endpoints para integração com frontends ou outros sistemas.

## 🚀 Começando

Siga os passos abaixo para rodar o `goBid` localmente.

### Pré-requisitos

- [Go](https://golang.org/dl/) (versão 1.16 ou superior)
- [Git](https://git-scm.com/downloads)

### Instalação

1. Clone o repositório:
   ```bash
   git clone https://github.com/DevKayoS/goBid.git
   cd goBid
   ```

2. Instale as dependências:
   ```bash
   go mod tidy
   ```

3. Configure as variáveis de ambiente:
   Crie um arquivo `.env` com base no `.env.example`:
   ```bash
   cp .env.example .env
   ```
   Edite o `.env` com suas configurações (ex.: porta, banco de dados, etc.).

4. Rode a aplicação:
   ```bash
   go run main.go
   ```

   O servidor estará disponível em `http://localhost:8080`.

## 🤝 Contribuindo

Quer ajudar a melhorar o `goBid`? Bora lá!

1. Faça um fork do repositório.
2. Crie uma branch para sua feature:
   ```bash
   git checkout -b minha-feature
   ```
3. Commit suas alterações:
   ```bash
   git commit -m "Adiciona minha feature"
   ```
4. Envie para o repositório remoto:
   ```bash
   git push origin minha-feature
   ```
5. Abra um Pull Request.

## 🌟 Agradecimentos

- À comunidade Go por ferramentas incríveis.
- A todos os contribuidores que ajudarem a tornar o `goBid` ainda melhor!

## 📬 Contato

Dúvidas ou ideias? Abre uma [issue](https://github.com/DevKayoS/goBid/issues) ou me contate:

- **GitHub**: [DevKayoS](https://github.com/DevKayoS)

---

⭐ **Curtiu? Dê uma estrela no repositório pra apoiar o projeto!**
