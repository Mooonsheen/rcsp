# RCSP

Учебный проект

В сервисе:
<ol>
  <li>Подключение и подписка на канал в nats-streaming</li>
  <li>Полученные данные записываются в Postgres</li>
  <li>Полученные данные сохраняются in memory в Redis (кэш)</li>
  <li>В случае падения сервиса кэш восстанавливается из Postgres</li>
  <li>http сервер, выдающий данные по id</li>
</ol>

