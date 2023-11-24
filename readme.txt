route: GET /todos
cURL
curl 'http://localhost:3000/todos'

response
{
    "success": true,
    "message": "You have no todos. Try adding one."
}

route: POST /create/todo
cURL
curl 'http://localhost:3000/create/todo' -H 'Content-Type: text/plain' -d '{
  "title": "Go shopping",
  "description": "Buy items from the supermarket"
}'

response:
{
    "success": true,
    "message": "Todo was created successfully."
}

Route: PUT /update/todo/{ID}
cURL
curl -X PUT 'http://localhost:3000/update/todo/1' -H 'Content-Type: text/plain' -d '{
  "title": "Put shopping",
  "description": "Take shopping out the bags and put away"
}'

response should be a 204 status code

Route: PATCH /complete/{ID}
cURL
curl -X PATCH 'http://localhost:3000/complete/1' -d ''

response should be 204 statud code

route: DELETE "/delete/{ID}
cURL
curl -X DELETE 'http://localhost:3000/delete/1' -d ''

resposne should be 204 stsus code
