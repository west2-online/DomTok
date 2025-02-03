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

#!/bin/bash

################ 密码 ################
PASSWORD=fzuhelper-password
################ 密码 ################


BASE_DIR=/mnt/disk/fzuhelper # SSL各种生成文件的基础路径

if test -d "$BASE_DIR/certs"; then # 判断证书是否已经存在了
    exit 0
fi

mkdir -p $BASE_DIR

CERT_OUTPUT_PATH="$BASE_DIR/certs" # 证书文件生成路径

KEY_STORE="$CERT_OUTPUT_PATH/kafka.keystore.jks" # Kafka keystore文件路径
TRUST_STORE="$CERT_OUTPUT_PATH/kafka.truststore.jks" # Kafka truststore文件路径
KEY_PASSWORD=$PASSWORD # keystore的key密码
STORE_PASSWORD=$PASSWORD # keystore的store密码
TRUST_KEY_PASSWORD=$PASSWORD # truststore的key密码
TRUST_STORE_PASSWORD=$PASSWORD # truststore的store密码
CLUSTER_NAME=kafka-cluster # 指定别名
CERT_AUTH_FILE="$CERT_OUTPUT_PATH/ca-cert" # CA证书文件路径
CLUSTER_CERT_FILE="$CERT_OUTPUT_PATH/${CLUSTER_NAME}-cert" # 集群证书文件路径
DAYS_VALID=3650 # key有效期
D_NAME="CN=FzuHelper, OU=West2Online, O=West2Online, L=Fujian, ST=Fujian, C=CN" # distinguished name

mkdir -p $CERT_OUTPUT_PATH

echo "1. 创建集群证书到keystore......"
keytool -keystore $KEY_STORE -alias $CLUSTER_NAME -validity $DAYS_VALID -genkey -keyalg RSA \
-storepass $STORE_PASSWORD -keypass $KEY_PASSWORD -dname "$D_NAME"

echo "2. 创建CA......"
openssl req -new -x509 -keyout $CERT_OUTPUT_PATH/ca-key -out "$CERT_AUTH_FILE" -days "$DAYS_VALID" \
-passin pass:"$PASSWORD" -passout pass:"$PASSWORD" \
-subj "/C=CN/ST=Beijing/L=Beijing/O=YourCompany/CN=Xi Hu"

echo "3. 导入CA文件到truststore......"
keytool -keystore "$TRUST_STORE" -alias CARoot \
-import -file "$CERT_AUTH_FILE" -storepass "$TRUST_STORE_PASSWORD" -keypass "$TRUST_KEY_PASSWORD" -noprompt

echo "4. 从key store中导出集群证书......"
keytool -keystore "$KEY_STORE" -alias "$CLUSTER_NAME" -certreq -file "$CLUSTER_CERT_FILE" -storepass "$STORE_PASSWORD" -keypass "$KEY_PASSWORD" -noprompt

echo "5. 签发证书......"
openssl x509 -req -CA "$CERT_AUTH_FILE" -CAkey $CERT_OUTPUT_PATH/ca-key -in "$CLUSTER_CERT_FILE" \
-out "${CLUSTER_CERT_FILE}-signed" \
-days "$DAYS_VALID" -CAcreateserial -passin pass:"$PASSWORD"

echo "6. 导入CA文件到keystore......"
keytool -keystore "$KEY_STORE" -alias CARoot -import -file "$CERT_AUTH_FILE" -storepass "$STORE_PASSWORD" \
 -keypass "$KEY_PASSWORD" -noprompt

echo "7. 导入已签发证书到keystore......"
keytool -keystore "$KEY_STORE" -alias "${CLUSTER_NAME}" -import -file "${CLUSTER_CERT_FILE}-signed" \
  -storepass "$STORE_PASSWORD" -keypass "$KEY_PASSWORD" -noprompt

rm $CERT_AUTH_FILE $CLUSTER_CERT_FILE "${CLUSTER_CERT_FILE}-signed"
