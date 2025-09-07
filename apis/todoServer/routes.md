| Method | Pattern                 | Handler     | Action                                  |
|--------|-------------------------|-------------|-----------------------------------------|
| GET    | /todo                   | rootHandler | Retrieve all to-do items                |
| GET    | /todo/{number}          |             | Retrieve a to-do item {number}          |
| POST   | /todo                   |             | Create a to-do item                     |
| PATCH  | /todo/{number}?complete |             | Mark a to-do item {number} as completed |
| DELETE | /todo/{number}          |             | Delete a to-do item {number}            |
