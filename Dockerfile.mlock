FROM djs55/go AS build

COPY pkg /go/src/github.com/djs55/hyperkit-measure-memory/pkg
COPY cmd /go/src/github.com/djs55/hyperkit-measure-memory/cmd

WORKDIR /go/src/github.com/djs55/hyperkit-measure-memory/cmd/mlock
RUN CGO_ENABLED=0 go build --ldflags '-s -w -extldflags "-static"' --buildmode pie

FROM scratch
COPY --from=build /go/src/github.com/djs55/hyperkit-measure-memory/cmd/mlock/mlock /mlock
CMD [ "/mlock" ]