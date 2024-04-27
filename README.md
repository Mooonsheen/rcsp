# РКСП

Учебный проект

В сервисе:
<ol>
  <li>Подключение и подписка на канал в nats-streaming</li>
  <li>Полученные данные записываются в Postgres</li>
  <li>Полученные данные сохраняются in memory в сервисе (кэш)</li>
  <li>В случае падения сервиса кэш восстанавливается из Postgres</li>
  <li>http сервер, выдающий данные по id из кэша</li>
</ol>

Сборка:
<ol>
  <li>make up - поднять контейнеры
  <li>make migrate - создать таблицу в Postgres
  <li>make server - поднять сервер
  <li>make publisher - собрать publisher
  <li>./publisher - запустить publisher
</ol>
publisher ожидает путь до json, который мы хотим добавить
