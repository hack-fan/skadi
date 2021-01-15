# Skadi

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

## Job Status
Job has status in its lifecycle.
* queuing: after user pushed the job to cloud
* sent: after the agent got the job
* expired: the job has been sent `10 minute` but no result, or the agent is offline when job is queuing
* succeeded: after the agent reported a succeeded result
* failed: after the agent reported a failed result

## Agent Status
All agent must check job every minute, if an agent has not checked job in `3 minute`,
it's status would be tagged as `offline`. All queuing job for this agent would be tagged as `expired`.
