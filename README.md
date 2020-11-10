# Server Fan

> This project is under development and is not available yet.

Easily run something on your server by sending message to IM.  slack, teams, wechat, etc...

## Architecture

```text
+-------------+                                  +-----------------------+
|             |                                  |                       |
|   E-mail    +---+                              |     Your  Servers     |
|             |   |                              |                       |
+-------------+   |                              |     +-----------+     |
                  |                              |     |           |     |
                  |                         +----------+   Agent   |     |
+-------------+   |   +--------------+      |    |     |           |     |
|             |   |   |              |      |    |     +-----------+     |
|  Chat Bot   +------->  Server Fan  <------+    |                       |
|             |   |   |              |      |    |                       |
+-------------+   |   +--------------+      |    |     +-----------+     |
                  |                         |    |     |           |     |
                  |                         +----------+   Agent   |     |
+-------------+   |                              |     |           |     |
|             |   |                              |     +-----------+     |
|   CI / CD   +---+                              |                       |
|             |                                  |                       |
+-------------+                                  +-----------------------+
```