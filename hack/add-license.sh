
PROJECT=$(cd $(dirname $0)/..; pwd)

LICENSEHEADERCHECKER_VERSION=v1.5.0

GOBIN=${PROJECT}/output go install github.com/lluissm/license-header-checker/cmd/license-header-checker@${LICENSEHEADERCHECKER_VERSION}

LICENSEIGNORE=$(cat ${PROJECT}/.licenseignore | tr '\n' ',')

${PROJECT}/output/license-header-checker -r -a -v -i ${LICENSEIGNORE} ${PROJECT}/hack/boilerplate.go.txt . go

${PROJECT}/output/license-header-checker -r -a -v -i ${LICENSEIGNORE} ${PROJECT}/hack/boilerplate.shell.txt . sh
