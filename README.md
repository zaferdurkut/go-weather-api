# Weather API

Hexagonal Architecture ile geliştirilmiş Go Weather API projesi.

## Proje Yapısı

## Özellikler

- Hexagonal Architecture
- Circuit Breaker Pattern
- RESTful API
- Weather API entegrasyonu
- Test coverage

## Kurulum

```bash
go mod tidy
go run cmd/server/main.go
```

## API Endpoints

- `GET /weather/:city` - Şehir için hava durumu bilgisi
- `GET /health` - Sağlık kontrolü