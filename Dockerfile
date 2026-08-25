FROM heroiclabs/nakama-pluginbuilder:3.26.0 AS builder

ENV GO111MODULE=on
ENV CGO_ENABLED=1

WORKDIR /backend

# Copy mod first for layer caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build Go plugin — matches template: trimpath + plugin mode
# Use vendor if present otherwise mod download; here we use mod=mod (no vendor dir yet)
RUN go build --trimpath --buildmode=plugin -o ./backend.so

FROM heroiclabs/nakama:3.26.0

# Copy compiled plugin into Nakama modules (Go runtime scans *.so)
COPY --from=builder /backend/backend.so /nakama/data/modules/rummy_backend.so
# Copy local.yml (config) — alternative to mounting via compose
COPY ./nakama/data/local.yml /nakama/data/local.yml
