## FyneAI

A simple tool to run AI models on your computer.

### Prequisites

- Golang(1.20+)
- [Ollama](https://ollama.ai/)
- [Fyne](https://fyne.io/)

### Building

- Clone the repository
- Install the dependencies

```bash
go mod tidy
```

- Build the project

```bash
go build -o fyneai
```

- Run the application

```bash
./fyneai
```

### Install On System

- Clone the repository
- Build with fyne, setting your respective OS. For linux, you can use:

```bash
fyne package -os linux -icon ai.png
```
- Extract the package
- Install the package. For linux it would be:

```bash
cd FyneAI && make user-install
```

- You'll find the application in your applications menu.

