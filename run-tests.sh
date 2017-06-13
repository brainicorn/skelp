#!/usr/bin/env bash

set -e

workdir=.coverage
covhtml=$workdir/all.html
outfile=coverage.txt

test -d $workdir || mkdir -p $workdir
[ -f "$covprofile" ] && rm ${covprofile}
[ -f "$outfile" ] && rm ${outfile}

echo "" > ${outfile}

i=0
for d in $(go list ./... | grep -v vendor); do
    go test -coverprofile=${workdir}/profile_$i.out -covermode=atomic $d
	i=$((i+1))
done

gocovmerge ${workdir}/*.out > ${outfile}
go tool cover -html=${outfile} -o ${covhtml}
rm $workdir/*.out