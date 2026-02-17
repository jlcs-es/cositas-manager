FROM node:24 as frontend-builder

COPY cositas-manager/ /cositas-manager
WORKDIR /cositas-manager
RUN npm ci
RUN npm run build


FROM golang:1.26.0 as backend-builder

COPY . /cositas-manager-back
COPY --from=frontend-builder /cositas-manager/dist/ /cositas-manager-back/cositas-manager/dist/
WORKDIR /cositas-manager-back

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o cositas-manager-bin main.go


FROM scratch
COPY --from=backend-builder /cositas-manager-back/cositas-manager-bin /bin/cositas-manager
EXPOSE 8080
CMD ["/bin/cositas-manager"]
