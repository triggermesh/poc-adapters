# Copyright 2022 TriggerMesh Inc.

# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at

#     http://www.apache.org/licenses/LICENSE-2.0

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

FROM ubuntu:latest as builder
RUN mkdir /app
RUN apt-get update && apt-get install wget -y && apt-get install unzip
RUN mkdir /app/.dw
RUN cd /app/.dw && wget -c https://github.com/mulesoft-labs/data-weave-cli/releases/download/v1.0.16/dw-1.0.16-Linux && unzip dw-1.0.16-Linux
RUN chmod +x /app/.dw/bin/dw
RUN wget -c https://go.dev/dl/go1.17.6.linux-amd64.tar.gz && tar -C /usr/local -xzf go1.17.6.linux-amd64.tar.gz


FROM builder
WORKDIR /app
ENV DW_HOME=/app/.dw
ENV DW_LIB_PATH=/app/.dw/libs
ENV PATH=/app/.dw/bin:$PATH
ENV PATH=$PATH:/usr/local/go/bin
COPY . ./
RUN cd cmd && go build -o /dwadapter
EXPOSE 8080
CMD ["/dwadapter"]
