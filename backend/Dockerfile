from alpine:3.17.3

RUN apk add --no-cache \
    ffmpeg \
    python3 \
    py3-pip \
    go \
    && pip3 install --no-cache-dir yt-dlp python-dotenv

WORKDIR /app
COPY . .
CMD ["go", "run", "main.go"]
