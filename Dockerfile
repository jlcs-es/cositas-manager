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
COPY --from=crazymax/7zip /usr/local/bin/7za /bin/7za
COPY --from=crazymax/7zip /usr/local/bin/7zr /bin/7zr
COPY --from=crazymax/7zip /usr/local/bin/7z /bin/7z
EXPOSE 8080
CMD ["/bin/cositas-manager"]
