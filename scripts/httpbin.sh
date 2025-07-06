# Copyright 2025 The packetd Authors
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

HEADER="X-APP-NAME: packetd"

while true; do
  curl -H "$HEADER" http://httpbin.org/get
  curl -H "$HEADER" -X POST http://httpbin.org/post
  curl -H "$HEADER" -X DELETE http://httpbin.org/delete
  curl -H "$HEADER" -X PUT http://httpbin.org/put
  curl -H "$HEADER" -X OPTION http://httpbin.org/option
done
