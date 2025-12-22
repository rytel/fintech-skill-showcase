# ETAP 1: Budowanie (builder)
FROM golang:1.23-alpine AS builder
WORKDIR /app

# 1. Kopiuj tylko go.mod (to musi istnieć po 'go mod init')
COPY go.mod .

# 2. Opcjonalnie kopiuj go.sum (ignoruje błąd, jeśli plik nie istnieje)
# Używamy --chown, ale kluczowe jest to, że nie ma tu osobnej instrukcji 'COPY go.sum'
# dla tego konkretnego kroku.

# 3. Pobieranie zależności. To polecenie *wygeneruje* go.sum w kontenerze.
# Używamy 'go mod download', ponieważ mamy już skopiowany go.mod
RUN go mod download

# 4. Kopiowanie reszty kodu (teraz dopiero main.go)
COPY . .

# Kompilacja aplikacji.
RUN go build -o server ./cmd/server


# ETAP 2: Uruchamianie (runner)
# Używamy bardzo małego obrazu 'alpine' lub 'scratch' (najmniejszy możliwy)
# ponieważ skompilowana binarka Go jest statycznie linkowana i nie potrzebuje 
# większości bibliotek systemowych. To sprawia, że obraz jest bardzo mały i bezpieczny.
FROM alpine:latest

# Ustawienie zmiennej środowiskowej
ENV PORT 8080

# Ustawienie katalogu roboczego dla etapu uruchamiania.
WORKDIR /root/

# Odkrywanie portu 8080. Jest to tylko informacja dla użytkownika obrazu, 
# nie ma wpływu na to, który port jest otwierany na hoście.
EXPOSE 8080

# Kopiowanie skompilowanej binarki z etapu 'builder' do tego, mniejszego, obrazu.
COPY --from=builder /app/server .

# Kopiowanie folderu migracji, aby aplikacja mogła odczytać plik SQL przy starcie
COPY --from=builder /app/migrations ./migrations

# Definicja polecenia, które zostanie uruchomione po starcie kontenera.
# Wskazujemy, aby uruchomić nasz skompilowany plik 'server'.
# Forma 'exec' ([...]) jest preferowana.
CMD ["./server"]