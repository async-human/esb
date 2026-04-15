import http from 'k6/http';
import { check, sleep } from 'k6';

// Настройки нагрузки
export const options = {
  stages:[
    { duration: '30s', target: 10000 },  // Разгон до 50 виртуальных юзеров за 30 сек
    { duration: '1m', target: 10000 },   // Держим полку в 50 юзеров 1 минуту
    { duration: '30s', target: 0 },   // Плавное снижение до 0
  ],
};

export default function () {
  const url = 'http://inbound-connector:8080/messages'; 

  const payload = JSON.stringify({
    event_id: `evt_${Math.floor(Math.random() * 1000000)}`,
    data: "test payload",
    timestamp: new Date().toISOString()
  });

  const params = {
    headers: {
      'Content-Type': 'application/json',
    },
  };

  const res = http.post(url, payload, params);

  // Проверяем, что inbound ответил 2xx (принял сообщение)
  check(res, {
    'is status 200 or 201': (r) => r.status === 200 || r.status === 201,
  });

  // Небольшая пауза между запросами юзера
  sleep(0.1); 
}