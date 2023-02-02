FROM debian:bullseye as fonts
RUN echo 'deb http://deb.debian.org/debian bullseye main contrib non-free' > /etc/apt/sources.list
RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates ttf-mscorefonts-installer 

# ---

FROM debian:bullseye as wkhtmltopdf
RUN apt-get update && apt-get install -y --no-install-recommends wget ca-certificates
RUN wget https://github.com/wkhtmltopdf/packaging/releases/download/0.12.6.1-2/wkhtmltox_0.12.6.1-2.bullseye_amd64.deb
RUN apt-get install -y --no-install-recommends /wkhtmltox_0.12.6.1-2.bullseye_amd64.deb

# ---

FROM golang:1.19 as builder
WORKDIR /build
COPY . .

RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN swag init
RUN go build -o html2pdf .

# ---

FROM debian:bullseye
WORKDIR /app
ENV WKHTMLTOPDF_PATH=/usr/bin/wkhtmltopdf
ENV GIN_MODE=release
EXPOSE 8080/tcp

RUN apt-get update && \
    apt-get install -y --no-install-recommends libjpeg62 libpng16-16 libxrender1 libfontconfig1 && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/* /usr/share/doc/*

COPY --from=fonts /usr/share/fonts/truetype/msttcorefonts/ /usr/share/fonts/truetype/msttcorefonts/
COPY --from=wkhtmltopdf /usr/local/bin/wkhtmltopdf /usr/bin/
COPY --from=builder /build/html2pdf /app/

HEALTHCHECK CMD curl --fail http://localhost:8080/api/healthcheck || exit 1

CMD ["./html2pdf"]
