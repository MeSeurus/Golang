```package main```

import (
    "errors"
    "fmt"
    "reflect"
    "regexp"
    "strconv"
    "strings"
)

```// Функция Validate принимает интерфейс любого типа и возвращает ошибку,
// если найдены проблемы с проверкой полей.
func Validate(v any) error {
    // получаем отражённую версию переданного аргумента
    val := reflect.ValueOf(v)
    typ := val.Type()

    // Проверяем, является ли аргумент структурой
    if typ.Kind() != reflect.Struct {
        return errors.New("неправильный тип аргумента: ожидается структура")
    }

    // проходим по всем полям структуры
    for i := 0; i < typ.NumField(); i++ {
        field := typ.Field(i)
        // Получаем тег validate
        tagValue := field.Tag.Get("validate")
        if tagValue == "" {
            continue
        }

        // Значение текущего поля
        fieldValue := val.Field(i).Interface()

        // Разбираем правила, отделённые друг от друга точкой с запятой
        rules := strings.Split(tagValue, ";")

        for _, rule := range rules {
            // Каждая проверка состоит из двух частей: ключ и значение
            kvPair := strings.SplitN(rule, "=", 2)
            if len(kvPair) != 2 {
                continue
            }

            checkType := kvPair[0]
            param := kvPair[1]

            switch checkType {
            case "min": // Минимальное значение
                minVal, err := strconv.Atoi(param)
                if err != nil {
                    continue
                }

                switch fv := fieldValue.(type) {
                case string:
                    runes := []rune(fv)
                    if len(runes) < minVal {
                        return fmt.Errorf("%s: минимальная длина '%d', фактическая длина '%d'", field.Name, minVal, len(runes))
                    }
                case int:
                    if fv < minVal {
                        return fmt.Errorf("%s: минимальное значение '%d', фактическое значение '%d'", field.Name, minVal, fv)
                    }
                }

            case "max": // Максимальное значение
                maxVal, err := strconv.Atoi(param)
                if err != nil {
                    continue
                }

                switch fv := fieldValue.(type) {
                case string:
                    runes := []rune(fv)
                    if len(runes) > maxVal {
                        return fmt.Errorf("%s: максимальная длина '%d', фактическая длина '%d'", field.Name, maxVal, len(runes))
                    }
                case int:
                    if fv > maxVal {
                        return fmt.Errorf("%s: максимальное значение '%d', фактическое значение '%d'", field.Name, maxVal, fv)
                    }
                }

            case "regexp": // Регулярное выражение
                re := regexp.MustCompile(param)
                fvStr, ok := fieldValue.(string)
                if !ok || !re.MatchString(fvStr) {
                    return fmt.Errorf("%s: некорректный формат строки", field.Name)
                }
            }
        }
    }

    return nil
}```

``// Пример структуры для тестирования
type User struct {
    Name  string `validate:"min=3"`      // минимум 3 символа
    Age   int    `validate:"min=18;max=65"` // возраст между 18 и 65 годами включительно
    Email string `validate:"regexp=^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$"` // правильный e-mail адрес
}``

```// Примеры тестовых случаев
func main() {
    user1 := User{
        Name:  "Ив",
        Age:   18,
        Email: "test@example.com",
    }
    
    err := Validate(user1)
    if err != nil {
        fmt.Printf("Ошибка валидации: %v\n", err)
    }

    user2 := User{
        Name:  "Иван",
        Age:   70,
        Email: "test@example.com",
    }
    err = Validate(user2)
    if err != nil {
        fmt.Printf("Ошибка валидации: %v\n", err)
    }

    user3 := User{
        Name:  "Иван",
        Age:   35,
        Email: "invalid_email",
    }
    err = Validate(user3)
    if err != nil {
        fmt.Printf("Ошибка валидации: %v\n", err)
    }

    user4 := User{
        Name:  "Иван",
        Age:   35,
        Email: "test@example.com",
    }
    err = Validate(user4)
    if err != nil {
        fmt.Printf("Ошибка валидации: %v\n", err)
    } else {
        fmt.Println("Валидация прошла успешно.")
    }```
}
