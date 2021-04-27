# Skadi

Skadi is a cloud message exchanger.

You can easily run something on your server by sending message to IM.Slack, teams, wechat, etc...

[LetServerRun](https://letserver.run) use this project as it infrastructure.[中文文档](https://letserver.run/docs/)

## Architecture

```text
+-------------+                                  +-----------------------+
|             |                                  |                       |
|   E-mail    +---+                              |     Your  Servers     |
|             |   |                              |                       |
+-------------+   |                              |     +-----------+     |
                  |                              |     |           |     |
                  |                         +----------+   Agent   |     |
+-------------+   |   +---------------+     |    |     |           |     |
|             |   |   |               |     |    |     +-----------+     |
|  Chat Bot   +------->  Skadi Cloud  <-----+    |                       |
|             |   |   |               |     |    |                       |
+-------------+   |   +---------------+     |    |     +-----------+     |
                  |                         |    |     |           |     |
                  |                         +----------+   Agent   |     |
+-------------+   |                              |     |           |     |
|             |   |                              |     +-----------+     |
|   CI / CD   +---+                              |                       |
|             |                                  |                       |
+-------------+                                  +-----------------------+
```

## Requirements

* Redis 3.2+
* MySQL

## Concepts

### Job
Job is just a message string, which sent by you from anywhere and would pull by specified Agent.

Job has status in its lifecycle.
* queuing: after user pushed the job to cloud
* sent: after the agent got the job
* expired: the job has been sent `10 minute` but no result, or the agent is offline when job is queuing
* succeeded: after the agent reported a succeeded result
* failed: after the agent reported a failed result

### Agent
Agent is a daemon running in your server, or a thread embedded in your service.

You can use our [open source agent](https://github.com/hack-fan/skadi-agent-shell),
or write your own using our [HTTP API](https://letserver.run/ref/).

All agent must check job every minute, if an agent has not checked job in `3 minute`,
it's status would be tagged as `offline`. All queuing job for this agent would be tagged as `expired`.

### Event

There are several kinds of events, you must handle them.

* EventMessage
* EventJobStatus

They will publish to a queue in redis.
