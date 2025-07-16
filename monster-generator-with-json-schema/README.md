# Generate Monster - JSON format

This project provides a simple way to generate unique and immersive monster for fantasy role-playing games. It uses a language model to create detailed monster sheets that include various aspects of the monster's life, personality, and background.

The file is genrated in the `contents` folder, with the name: `monster_sheet.json`.

## Model

```bash
docker model pull ai/qwen2.5:latest
```
> If you use Docker Compose, this will be pulled automatically.

## Start the application

**With Docker Compose**:
```bash
docker compose up --build
```
> Specify the `MONSTER_KIND` environment variable to generate a character of a specific kind (e.g., `werewolf`, `dragon`, etc.).

**From a container**:
```bash
MODEL_RUNNER_BASE_URL=http://model-runner.docker.internal/engines/llama.cpp/v1 \
MODEL_RUNNER_CHAT_MODEL=ai/qwen2.5:latest \
MONSTER_KIND=werewolf \
go run main.go
```

**From a local machine**:
```bash
MODEL_RUNNER_BASE_URL=http://localhost:12434/engines/llama.cpp/v1 \
MODEL_RUNNER_CHAT_MODEL=ai/qwen2.5:latest \
MONSTER_KIND=werewolf \
go run main.go
```
