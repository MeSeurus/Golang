package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Структура задания
type Job struct {
	ID  int
	URL string
}

// Структура результата
type Result struct {
	Job      Job
	Status   string
	Duration time.Duration
}

// Имитация HTTP-запроса с задержкой
func simulateHttpRequest(job Job) (*Result, error) {
	// Задаем случайную задержку - до 1000 мс
	delay := rand.Intn(500) + 500
	time.Sleep(time.Millisecond * time.Duration(delay))
	// Возвращаем ответ по результатам обработки
	result := &Result{
		Job:      job,
		Status:   "Success",
		Duration: time.Duration(delay),
	}
	return result, nil
}

// Worker
func worker(jobs <-chan Job, results chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()

	for job := range jobs {
		result, err := simulateHttpRequest(job)
		if err != nil {
			results <- Result{
				Job:    job,
				Status: "Error",
			}
		} else {
			results <- *result
		}
	}
}

const numWorkers = 5 // Количество Worker

func main() {
	rand.Seed(time.Now().UnixNano())

	var wg sync.WaitGroup
	// Зададим буферизованный канал для заданий
	jobs := make(chan Job, 100)
	// Зададим буферизованный канал для результатов
	results := make(chan Result, 100)

	// Список URL-адресов
	urls := []string{"https://google.com", "https://yandex.com", "https://yahoo.com", "https://rambler.ru", "https://bing.com"}

	// Запускаем Worker
	wg.Add(numWorkers)
	for range numWorkers {
		go worker(jobs, results, &wg)
	}

	// Отправка заданий в канал
	for id, url := range urls {
		job := Job{
			ID:  id,
			URL: url,
		}
		jobs <- job
	}
	// Закрываем канал
	close(jobs)

	// Ждем завершения работы всех Worker
	go func() {
		wg.Wait()
		close(results)
	}()

	// Чтение результатов и сбор статистики
	allResults := make([]Result, 0, len(urls))
	totalTime := time.Duration(0)

	for result := range results {
		allResults = append(allResults, result)
		totalTime += result.Duration
	}

	// Среднее время выполнения запросов
	avgTime := totalTime / time.Duration(len(allResults))

	// Вывод статистики
	fmt.Println("\nСтатистика:")
	for _, res := range allResults {
		fmt.Printf("%s\t%s\tВремя выполнения: %v\n", res.Job.URL, res.Status, res.Duration)
	}
	fmt.Printf("\nСреднее время выполнения: %v\n", avgTime)
}
