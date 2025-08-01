# Проект на Go

## Установка и настройка

### 1. Установка Go на сервер (Ubuntu/Debian)

```bash
# Обновление пакетов и установка Go
sudo apt update
sudo apt install -y golang

# Проверка установленной версии
go version

# Настройка GOPATH (если нужно)
echo 'export GOPATH=$HOME/go' >> ~/.bashrc
echo 'export PATH=$PATH:$GOPATH/bin' >> ~/.bashrc
source ~/.bashrc

# Клонирование репозитория
git clone https://github.com/ваш/репозиторий.git
cd репозиторий

# Установка зависимостей
go mod download

# Непосредственный запуск
go run main.go
```
