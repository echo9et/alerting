# cmd/staticlint

Включенные анализаторы:
    - atomic: проверяет корректность использования пакета sync/atomic.
    - bools: обнаруживает подозрительные булевы выражения.
    - errorsas: проверяет правильность использования errors.As.
    - printf: проверяет форматирование строковых операций.
    - shadow: обнаруживает затенение переменных.
    - structtag: проверяет корректность тегов структур.
    - tests: проверяет тестовые файлы.
    - unmarshal: проверяет корректность работы с JSON/XML.
    - unreachable: обнаруживает недостижимый код.
    - unsafeptr: проверяет корректность использования unsafe.Pointer.
    - staticcheck (SA): анализаторы из пакета staticcheck.io.
    - errcheck: проверяет необработанные ошибки.
    - unused: обнаруживает неиспользуемые параметры функций.
    - noosexitinmain: собственный анализатор, запрещающий использование os.Exit в main.

Запуск multichecker:
    go run cmd/staticlint/main.go ./...
