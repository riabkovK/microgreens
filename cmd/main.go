package main

import (
	"github.com/riabkovK/microgreens/internal/app"
)

// TODO Добавить:
// 1) рефреш токен +
// 2) сборку бинаря + проброс конфига (соответственно, флаги), поиск переменных из окружения (viper)
// 3) docker compose,
// 4) Тесты
// 5) swagger
// 6) CI?
// 7) UI (на реакте)

func main() {
	app.Run()
}
