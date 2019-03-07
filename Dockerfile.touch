FROM djs55/go AS build

COPY pkg /go/src/github.com/djs55/hyperkit-measure-memory/pkg
COPY cmd /go/src/github.com/djs55/hyperkit-measure-memory/cmd

WORKDIR /go/src/github.com/djs55/hyperkit-measure-memory/cmd/touch
RUN cc -o touch main.c -static -pie

FROM scratch
COPY --from=build /go/src/github.com/djs55/hyperkit-measure-memory/cmd/touch/touch /touch
CMD [ "/touch" ]

