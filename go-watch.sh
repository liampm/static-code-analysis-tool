#!/bin/sh

if ! which inotifyd > /dev/null; then
    go run main.go;
    exit 0;
fi

go run main.go &

while true; do
    inotifywait -r -e modify -e attrib -e moved_to -e moved_from -e create -e delete -e move /go/src;
    # TODO: check what changed
    echo -n 'Killing current processes... ';
    ps -C go -o pid= | xargs kill -9;
    ps -C main -o pid= | xargs kill -9;
    echo -e "\033[32mDone!\033[0m";
    go run main.go &
done

