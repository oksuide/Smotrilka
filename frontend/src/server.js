const WebSocket = require('ws');

const wss = new WebSocket.Server({ port: 5000 });

wss.on('connection', (ws) => {
    ws.on('message', (message) => {
        try {
            const parsedMessage = JSON.parse(message);
            console.log('Получено сообщение:', parsedMessage.text);

            // Отправляем сообщение всем подключенным клиентам
            wss.clients.forEach(client => {
                if (client.readyState === WebSocket.OPEN) {
                    client.send(JSON.stringify({ text: parsedMessage.text }));
                }
            });

        } catch (error) {
            console.error('Ошибка обработки сообщения:', error);
        }
    });
});

console.log('WebSocket сервер запущен на порту 5000');
