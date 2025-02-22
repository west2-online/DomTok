# Copyright 2024 The west2-online Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

PROJECT=$(cd $(dirname $0)/..; pwd)

LICENSEHEADERCHECKER_VERSION=v1.5.0

GOBIN=${PROJECT}/output go install github.com/lluissm/license-header-checker/cmd/license-header-checker@${LICENSEHEADERCHECKER_VERSION}

LICENSEIGNORE=$(cat ${PROJECT}/.licenseignore | tr '\n' ',')

${PROJECT}/output/license-header-checker -r -a -v -i ${LICENSEIGNORE} ${PROJECT}/hack/boilerplate.go.txt . go

${PROJECT}/output/license-header-checker -r -a -v -i ${LICENSEIGNORE} ${PROJECT}/hack/boilerplate.shell.txt . sh