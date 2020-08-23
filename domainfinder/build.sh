#!/bin/bash
DFDIR=domainfinder
mkdir -p lib
go build -o domainfinder
cd ../synonyms
go build -o ../${DFDIR}/lib/synonyms
cd ../available
go build -o ../${DFDIR}/lib/available
cd ../sprinkle
go build -o ../${DFDIR}/lib/sprinkle
cd ../coolify
go build -o ../${DFDIR}/lib/coolify
cd ../domainfinder
go build -o ../${DFDIR}/lib/domainify
echo "Complete!"
